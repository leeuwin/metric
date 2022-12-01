package metric

import "github.com/prometheus/client_golang/prometheus"

type gauge struct {
	vec *prometheus.GaugeVec
}

func NewGauge(metric string, opts ...option) gauge {

	return gauge{
		vec: newGaugeVec(metric, gLabelNames, opts...),
	}
}

func (g gauge) Set(v float64) {

	if gg := g.getCollector(); gg != nil {
		gg.Set(v)
	}
}

func (g gauge) Inc() {

	if gg := g.getCollector(); gg != nil {
		gg.Inc()
	}
}

func (g gauge) Dec() {

	if gg := g.getCollector(); gg != nil {
		gg.Dec()
	}
}

func (g gauge) Add(v float64) {

	if gg := g.getCollector(); gg != nil {
		gg.Add(v)
	}
}

func (g gauge) Sub(v float64) {

	if gg := g.getCollector(); gg != nil {
		gg.Sub(v)
	}
}

func (g gauge) SetToCurrentTime() {

	if gg := g.getCollector(); gg != nil {
		gg.SetToCurrentTime()
	}
}

func (g gauge) getCollector() prometheus.Gauge {

	if IsEnable() {
		return g.vec.WithLabelValues(makeLabelValues(nil)...)
	}
	return nil
}

func newGaugeVec(metric string, labelNames []string, opts ...option) *prometheus.GaugeVec {

	labels := make(Param).Init(opts...)

	gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        metric,
		ConstLabels: prometheus.Labels(labels),
		Help:        "",
	}, labelNames)

	prometheus.DefaultRegisterer.MustRegister(gaugeVec)
	return gaugeVec
}
