package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type counterWith struct {
	c counter
}

//带自定义label的指标，必须配套调用with的操作且提供分别匹配的value值
func NewCounterWith(metric string, labelNames []string, opts ...option) counterWith {

	return counterWith{
		c: counter{
			vec: newCounterVec(metric, makeLabelNames(labelNames), opts...),
		},
	}
}

func (cw counterWith) IncWith(labelValues []string) {

	cw.c.getCollectorWithLavelValues(labelValues).Inc()
}

func (cw counterWith) AddWith(v float64, labelValues []string) {

	cw.c.getCollectorWithLavelValues(labelValues).Add(v)
}

func (c counter) getCollectorWithLavelValues(labelValues []string) prometheus.Counter {

	if IsEnable() {
		return c.vec.WithLabelValues(makeLabelValues(labelValues)...)
	}
	return nil
}
