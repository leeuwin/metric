package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type counter struct {
	vec *prometheus.CounterVec
}

func NewCounter(metric string, opts ...option) counter {

	return counter{
		vec: newCounterVec(metric, gLabelNames, opts...),
	}
}

func (c counter) Inc() {

	if cnt := c.getCollector(); cnt != nil {
		cnt.Inc()
	}
}

func (c counter) Add(v float64) {

	if cnt := c.getCollector(); cnt != nil {
		cnt.Add(v)
	}
}

func (c counter) getCollector() prometheus.Counter {

	if IsEnable() {
		return c.vec.WithLabelValues(makeLabelValues(nil)...)
	}
	return nil
}

func newCounterVec(metric string, labelNames []string, opts ...option) *prometheus.CounterVec {

	labels := make(Param).Init(opts...)

	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        metric,
		ConstLabels: prometheus.Labels(labels),
		Help:        "",
	}, labelNames)
	prometheus.DefaultRegisterer.MustRegister(counterVec)
	return counterVec
}
