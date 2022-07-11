package metrics

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "github.com/shashi-banger/metrics_grpc_sidecar/metricspb/metricspb_v1"
)

type MetricType string

var (
	MetricTypeCounter   = "Counter"
	MetricTypeGauge     = "Gauge"
	MetricTypeHistogram = "Histogram"
	MetricTypeSummary   = "Summary"
)

var (
	ErrorAlreadyExists    = "ErrorAlreadyExists"
	ErrorUnknownCollector = "ErrorUnknownCollector"
)

type MetricsContext struct {
	counterCollectors   map[string]*prometheus.CounterVec
	gaugeCollectors     map[string]*prometheus.GaugeVec
	histogramCollectors map[string]*prometheus.HistogramVec
	summaryCollectors   map[string]*prometheus.SummaryVec
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
		metricsContext.counterCollectors = make(map[string]*prometheus.CounterVec)
		metricsContext.gaugeCollectors = make(map[string]*prometheus.GaugeVec)
		metricsContext.histogramCollectors = make(map[string]*prometheus.HistogramVec)
		metricsContext.summaryCollectors = make(map[string]*prometheus.SummaryVec)
	})
}

// Following call bloacks and hence should be run as go routine
func StartMetricsEndpointServer() {
	fmt.Println("Starting http server")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
	fmt.Println("Stopping http server")
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

func CreateCounter(params *v1.CreateCounterParams) error {
	var err error
	counterCollector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: params.Name,
			Help: params.Help,
		},
		params.Labels,
	)

	prometheus.MustRegister(counterCollector)

	collector_key := create_collector_key_from_list(params.Name, params.Labels)
	if _, ok := metricsContext.counterCollectors[collector_key]; ok {
		err = fmt.Errorf("%s: ccounter ollector_key %s alreadyExists", ErrorAlreadyExists, collector_key)
	} else {
		metricsContext.counterCollectors[collector_key] = counterCollector
	}
	return err
}

func CounterInc(name string, labelValues map[string]string) error {
	var err error
	collector_key := create_collector_key_from_map(name, labelValues)

	if val, ok := metricsContext.counterCollectors[collector_key]; ok {
		val.With(prometheus.Labels(labelValues)).Inc()
	} else {
		err = fmt.Errorf("%s: counter collector_key %s not registered", ErrorUnknownCollector, collector_key)
	}
	return err
}

func CreateGauge(params *v1.CreateGaugeParams) error {
	var err error
	gaugeCollector := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: params.Name,
			Help: params.Help,
		},
		params.Labels,
	)

	prometheus.MustRegister(gaugeCollector)

	collector_key := create_collector_key_from_list(params.Name, params.Labels)
	if _, ok := metricsContext.gaugeCollectors[collector_key]; ok {
		err = fmt.Errorf("%s: gauge collector_key %s alreadyExists", ErrorAlreadyExists, collector_key)
	} else {
		metricsContext.gaugeCollectors[collector_key] = gaugeCollector
	}
	return err
}

func GaugeSet(name string, labelValues map[string]string, value float64) error {
	var err error
	collector_key := create_collector_key_from_map(name, labelValues)

	if val, ok := metricsContext.gaugeCollectors[collector_key]; ok {
		val.With(prometheus.Labels(labelValues)).Set(value)
	} else {
		err = fmt.Errorf("%s: gauge collector_key %s not registered", ErrorUnknownCollector, collector_key)
	}
	return err
}

func CreateHistogram(params *v1.CreateHistogramParams) error {
	var err error
	histogramCollector := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    params.Name,
			Help:    params.Help,
			Buckets: params.Buckets,
		},
		params.Labels,
	)
	prometheus.MustRegister(histogramCollector)
	collector_key := create_collector_key_from_list(params.Name, params.Labels)
	if _, ok := metricsContext.histogramCollectors[collector_key]; ok {
		err = fmt.Errorf("%s: histogram collector_key %s alreadyExists", ErrorAlreadyExists, collector_key)
	} else {
		metricsContext.histogramCollectors[collector_key] = histogramCollector
	}
	return err
}

func HistogramObserve(name string, labelValues map[string]string, value float64) error {
	var err error
	collector_key := create_collector_key_from_map(name, labelValues)

	if val, ok := metricsContext.histogramCollectors[collector_key]; ok {
		val.With(prometheus.Labels(labelValues)).Observe(value)
	} else {
		err = fmt.Errorf("%s: histogram collector_key %s not registered", ErrorUnknownCollector, collector_key)
	}
	return err
}
