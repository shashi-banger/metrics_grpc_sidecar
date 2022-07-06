package metrics

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	v1 "github.com/shashi-banger/metrics_grpc_sidecar/metricspb/v1"
)

type MetricType string

var (
	MetricTypeCounter   = "Counter"
	MetricTypeGauge     = "Gauge"
	MetricTypeHistogram = "Histogram"
	MetricTypeSummary   = "Summary"
)

type MetricsContext struct {
	counterCollectors   map[string]prometheus.Collector
	gaugeCollectors     map[string]prometheus.Collector
	histogramCollectors map[string]prometheus.Collector
	summaryCollectors   map[string]prometheus.Collector
}

var (
	metricsContext *MetricsContext
	once           sync.Once
)

/*
type CreateCounterParams struct {
	Name   string
	Labels []string
	Help   string
}

type CreateGaugeParams struct {
	Name   string
	Labels []string
	Help   string
}

type CreateHistogramParams struct {
	Name   string
	Labels []string
	Help   string
	// Usage and initialization of Buckets is described here
	// https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#HistogramOpts
	Buckets []float64
}
*/

func init() {
	once.Do(func() {
		metricsContext = &MetricsContext{}
	})
}

func create_collector_key_from_list(name string, all_labels []string) string {
	sort.Strings(all_labels)
	all_keys_str := strings.Join(all_labels, " ")
	h := fnv.New32a()
	h.Write([]byte(all_keys_str))
	final_str := fmt.Sprintf("%s_%s", "num_bytes", strconv.FormatUint(uint64(h.Sum32()), 10))
	return final_str
}

func create_collector_key_from_map(name string, labelValues map[string]string) string {
	all_labels := []string{}
	for k, _ := range labelValues {
		all_labels = append(all_labels, k)
	}

	return create_collector_key_from_list(name, all_labels)

}

func CreateCounter(params v1.CreateCounterParams) error {
	var err error
	counterCollector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: params.Name,
			Help: params.Help,
		},
		params.Labels,
	)
	collector_key := create_collector_key_from_list(params.Name, params.Labels)
	if _, ok := metricsContext.counterCollectors[collector_key]; ok {
		err = fmt.Errorf("AlreadyExists: ccounter ollector_key %s alreadyExists", collector_key)
	} else {
		metricsContext.counterCollectors[collector_key] = counterCollector
	}
	return err
}

func CounterInc(name string, labelValues map[string]string) error {
	var err error
	collector_key := create_collector_key_from_map(name, labelValues)

	if val, ok := metricsContext.counterCollectors[collector_key]; ok {
		val.Inc()
	} else {
		err = fmt.Errorf("UnknownCollector: counter collector_key %s not registered", collector_key)
	}
	metricsContext.counterCollectors[collector_key].Inc()
	return err
}

func CreateGauge(params CreateGaugeParams) error {
	var err error
	gaugeCollector := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: params.Name,
			Help: params.Help,
		},
		params.Labels,
	)
	collector_key := create_collector_key_from_list(params.Name, params.Labels)
	if _, ok := metricsContext.gaugeCollectors[collector_key]; ok {
		err = fmt.Errorf("AlreadyExists: gauge collector_key %s alreadyExists", collector_key)
	} else {
		metricsContext.gaugeCollectors[collector_key] = gaugeCollector
	}
	return err
}

func GaugeSet(name string, labelValues map[string]string, value float64) error {
	var err error
	collector_key := create_collector_key_from_map(name, labelValues)

	if val, ok := metricsContext.gaugeCollectors[collector_key]; ok {
		val.Inc()
	} else {
		err = fmt.Errorf("UnknownCollector: gauge collector_key %s not registered", collector_key)
	}
	metricsContext.gaugeCollectors[collector_key].Set(value)
	return err
}

func CreateHistogram(params CreateHistogramParams) error {
	var err error
	histogramCollector := prometheus.NewHistogramVec(
		prometheus.GaugeOpts{
			Name:    params.Name,
			Help:    params.Help,
			Buckets: params.Buckets,
		},
		params.Labels,
	)
	collector_key := create_collector_key_from_list(params.Name, params.Labels)
	if _, ok := metricsContext.histogramCollectors[collector_key]; ok {
		err = fmt.Errorf("AlreadyExists: histogram collector_key %s alreadyExists", collector_key)
	} else {
		metricsContext.histogramCollectors[collector_key] = histogramCollector
	}
	return err
}

func HistogramObserve(name string, labelValues map[string]string, value float64) error {
	var err error
	collector_key := create_collector_key_from_map(name, labelValues)

	if val, ok := metricsContext.histogramCollectors[collector_key]; ok {
		val.Inc()
	} else {
		err = fmt.Errorf("UnknownCollector: histogram collector_key %s not registered", collector_key)
	}
	metricsContext.histogramCollectors[collector_key].Observe(value)
	return err
}
