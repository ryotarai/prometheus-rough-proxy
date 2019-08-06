package handler

import (
	"encoding/json"
	"github.com/prometheus/common/model"
	"github.com/ryotarai/prometheus-rough-proxy/lib/client"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type handler struct {
	prometheusURL *url.URL
	reverseProxy  http.Handler
	client        *client.Client
}

func New(prometheusURL *url.URL, client *client.Client) (*handler, error) {
	return &handler{
		prometheusURL: prometheusURL,
		reverseProxy:  httputil.NewSingleHostReverseProxy(prometheusURL),
		client:        client,
	}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	directMode := r.URL.Query().Get("direct") != ""

	if r.URL.Path == "/api/v1/query_range" && !directMode {
		h.handleQueryRange(w, r)
	} else {
		log.Printf("Proxying to upstream (url: %s)", r.URL)
		h.reverseProxy.ServeHTTP(w, r)
	}
}

func (h *handler) handleQueryRange(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("query")

	start, err := parsePromTime(q.Get("start"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err, errorNone)
		return
	}

	end, err := parsePromTime(q.Get("end"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err, errorNone)
		return
	}

	step, err := parsePromDuration(q.Get("step"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err, errorNone)
		return
	}

	result, err := h.client.QueryRangeByQuery(query, start, end, step)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err, errorNone)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")

	e := json.NewEncoder(w)
	e.Encode(apiResponse{
		Status: "success",
		Data: queryData{
			ResultType: model.ValMatrix,
			Result:     result,
		},
	})
}

func (h *handler) writeError(w http.ResponseWriter, status int, err error, t errorType) {
	w.WriteHeader(status)
	w.Header().Set("content-type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(apiResponse{
		Status:    statusError,
		ErrorType: t,
		Error:     err.Error(),
	})
}
