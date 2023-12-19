package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
)

var index = "go-api-app-metrics"

type Metrics struct {
	Metrics             map[string]any
	elasticsearchClient *elasticsearch.Client
}

func NewMetrics(elasticsearchClient *elasticsearch.Client) *Metrics {
	m := &Metrics{
		Metrics:             make(map[string]any),
		elasticsearchClient: elasticsearchClient,
	}

	m.Add("_timestamp", time.Now().Format(time.RFC3339))
	return m
}

func (m *Metrics) Add(key string, value any) {
	m.Metrics[key] = value
}

func (m *Metrics) AddAll(from map[string]any) {
	for key, value := range from {
		m.Metrics[key] = value
	}
}

func (m *Metrics) MetrifyExecutionTime(key string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		m.Add(fmt.Sprintf("%s_exec_time", key), elapsed)
	}
}

func (m *Metrics) Send() {
	m.Add("_level", "info")
	m.send()
}

func (m *Metrics) SendError(err error) {
	m.Add("_level", "error")
	m.Add("error", err.Error())
	m.send()
}

func (m *Metrics) send() {
	go func() {
		doc, _ := json.Marshal(m.Metrics)

		req := esapi.IndexRequest{
			Index:      index,
			DocumentID: uuid.New().String(),
			Body:       bytes.NewReader(doc),
			Refresh:    "true",
		}

		_, _ = req.Do(context.Background(), m.elasticsearchClient)
	}()
}
