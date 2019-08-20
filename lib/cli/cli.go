package cli

import (
	"flag"
	"fmt"
	"github.com/ryotarai/prometheus-rough-proxy/lib/client"
	"github.com/ryotarai/prometheus-rough-proxy/lib/handler"
	"log"
	"net/http"
	"net/url"
)

type options struct {
	listen string
	prometheusURL string
	apiConcurrency int
}

func Start(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	opts := options{}
	fs.StringVar(&opts.listen, "listen", ":9090", "Listen address")
	fs.StringVar(&opts.prometheusURL, "prometheus-url", "", "URL to Prometheus")
	fs.IntVar(&opts.apiConcurrency, "api-concurrency", 1, "Concurrency limit for Prometheus API")

	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	if opts.prometheusURL == "" {
		return fmt.Errorf("-prometheus-url is required")
	}

	prometheusURL, err := url.Parse(opts.prometheusURL)
	if err != nil {
		return err
	}

	c, err := client.New(prometheusURL, opts.apiConcurrency)
	if err != nil {
		return err
	}

	h, err := handler.New(prometheusURL, c)
	if err != nil {
		return err
	}

	log.Printf("Starting prometheus-rough-proxy (listening on %s)", opts.listen)
	if err := http.ListenAndServe(opts.listen, h); err != nil {
		return err
	}

	return nil
}
