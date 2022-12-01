package metric

import "github.com/prometheus/client_golang/prometheus"

var objectives map[float64]float64 = map[float64]float64{
	0.5:  0.05,  //0.45 ~ 0.55
	0.9:  0.01,  //0.89 ~ 0.91
	0.95: 0.005, //0.945 ~ 0.955
	0.99: 0.001, //0.989 ~ 0.991
}

type summary struct {
	vec *prometheus.SummaryVec
}

func NewSummary(metric string, opts ...option) summary {

	return summary{
		vec: newSummaryVec(metric, gLabelNames, opts...),
	}
}

func (s summary) Observe(v float64) {

	if obsv := s.getCollector(); obsv != nil {
		obsv.Observe(v)
	}
}

func (s summary) getCollector() prometheus.Observer {

	if IsEnable() {
		return s.vec.WithLabelValues(makeLabelValues(nil)...)
	}
	return nil
}

func newSummaryVec(metric string, labelNames []string, opts ...option) *prometheus.SummaryVec {

	labels := make(Param).Init(opts...)

	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   "",
		Subsystem:   "",
		Name:        metric,
		ConstLabels: prometheus.Labels(labels),
		Help:        "",
		Objectives:  objectives,
	}, labelNames)

	prometheus.DefaultRegisterer.MustRegister(summaryVec)
	return summaryVec
}
