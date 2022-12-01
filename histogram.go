package metric

import "github.com/prometheus/client_golang/prometheus"

type histogram struct {
	vec *prometheus.HistogramVec
}

var (
	/*
		1.要求严格递增
		2.不必设置最大值，会自动追加
		3.若bucket为nil或者len为0，采用默认bucket
	*/
	buckets = [][]float64{
		{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}, //this is default bucket
		{.05, .1, .25, .5, 1.0, 2.5, 5, 10, 25, 50, 100},
		{0.5, 1, 2.5, 5, 10, 25, 50, 100, 250, 500, 1000},
		{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
	}
)

const (
	BucketSchemaDefault = iota
	BucketSchemaHundred
	BucketSchemaThousand
	BucketSchema10Thousand
)

func NewHistogram(metric string, bucketSchema int, opts ...option) histogram {

	return histogram{
		vec: newHistogramVec(metric, bucketSchema, gLabelNames, opts...),
	}
}

func (h histogram) Observe(v float64) {

	if obsv := h.getCollector(); obsv != nil {
		obsv.Observe(v)
	}
}

func (h histogram) getCollector() prometheus.Observer {

	if IsEnable() {
		return h.vec.WithLabelValues(makeLabelValues(nil)...)
	}
	return nil
}

func newHistogramVec(metric string, bucketSchema int, labelNames []string, opts ...option) *prometheus.HistogramVec {

	labels := make(Param).Init(opts...)

	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        metric,
		ConstLabels: prometheus.Labels(labels),
		Help:        "",
		Buckets:     buckets[bucketSchema],
	}, labelNames)

	prometheus.DefaultRegisterer.MustRegister(histogramVec)
	return histogramVec
}
