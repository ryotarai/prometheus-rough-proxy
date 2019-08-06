package client

import (
	"context"
	"github.com/prometheus/common/model"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	promapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type Client struct {
	prometheusURL *url.URL
	api           promv1.API
	concurrencyCh chan struct{}
}

func New(prometheusURL *url.URL, concurrency int) (*Client, error) {
	client, err := promapi.NewClient(promapi.Config{
		RoundTripper: http.DefaultTransport,
		Address:      prometheusURL.String(),
	})
	if err != nil {
		return nil, err
	}

	api := promv1.NewAPI(client)

	return &Client{
		prometheusURL: prometheusURL,
		api:           api,
		concurrencyCh: make(chan struct{}, concurrency),
	}, nil
}

func (c *Client) QueryRangeByQuery(query string, start time.Time, end time.Time, step time.Duration) (model.Matrix, error) {
	type respT struct {
		value    model.Value
		warnings promapi.Warnings
		err      error
	}
	respCh := make(chan respT)

	count := 0
	t := end
	for !t.Before(start) {
		u := t
		go func() {
			c.concurrencyCh <- struct{}{}
			defer func() {
				<-c.concurrencyCh
			}()

			log.Printf("DEBUG: Querying to upstream (t: %s, query: '%v')", u, query)
			value, warnings, err := c.api.Query(context.Background(), query, u)
			respCh <- respT{value, warnings, err}
		}()

		count++
		t = t.Add(step * -1)
	}

	resps := make([]respT, count)
	for i := 0; i < count; i++ {
		resps[i] = <-respCh
	}

	metrics := map[string]model.Metric{}
	values := map[string][]model.SamplePair{}

	for _, resp := range resps {
		if resp.err != nil {
			return nil, resp.err
		}

		vec := resp.value.(model.Vector)

		for _, sample := range vec {
			key := metricToKey(sample.Metric)
			if _, ok := metrics[key]; !ok {
				metrics[key] = sample.Metric
				values[key] = []model.SamplePair{}
			}
			values[key] = append(values[key], model.SamplePair{
				Timestamp: sample.Timestamp,
				Value:     sample.Value,
			})
		}
	}

	result := model.Matrix{}
	for key, metric := range metrics {
		v := values[key]
		sort.Slice(v, func(i, j int) bool {
			return v[i].Timestamp.Before(v[j].Timestamp)
		})
		result = append(result, &model.SampleStream{
			Metric: metric,
			Values: v,
		})
	}

	return result, nil
}

func metricToKey(metric model.Metric) string {
	keys := []string{}
	for k := range metric {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	escape := func(s string) string {
		return strings.ReplaceAll(s, "\n", "\n\n")
	}

	var builder strings.Builder
	for _, k := range keys {
		builder.WriteString(escape(k))
		builder.WriteString("\n")
		builder.WriteString(escape(string(metric[model.LabelName(k)])))
		builder.WriteString("\n")
	}

	return builder.String()
}
