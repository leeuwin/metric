package metric

import "github.com/prometheus/client_golang/prometheus"

type gaugeWith struct {
	g gauge
}

//带自定义label的指标，必须配套调用with的操作且提供分别匹配的value值
func NewGaugeWith(metric string, labelNames []string, opts ...option) gaugeWith {

	return gaugeWith{
		g: gauge{
			vec: newGaugeVec(metric, makeLabelNames(labelNames), opts...),
		},
	}
}

func (gw gaugeWith) SetWith(v float64, labelValues []string) {

	if gg := gw.getCollectorWithLabelValues(labelValues); gg != nil {
		gg.Set(v)
	}
}

func (gw gaugeWith) IncWith(labelValues []string) {

	if gg := gw.getCollectorWithLabelValues(labelValues); gg != nil {
		gg.Inc()
	}
}

func (gw gaugeWith) DecWith(labelValues []string) {

	if gg := gw.getCollectorWithLabelValues(labelValues); gg != nil {
		gg.Dec()
	}
}

func (gw gaugeWith) AddWith(v float64, labelValues []string) {

	if gg := gw.getCollectorWithLabelValues(labelValues); gg != nil {
		gg.Add(v)
	}
}

func (gw gaugeWith) SubWith(v float64, labelValues []string) {

	if gg := gw.getCollectorWithLabelValues(labelValues); gg != nil {
		gg.Sub(v)
	}
}

func (gw gaugeWith) SetToCurrentTimeWith(labelValues []string) {

	if gg := gw.getCollectorWithLabelValues(labelValues); gg != nil {
		gg.SetToCurrentTime()
	}
}

func (gw gaugeWith) getCollectorWithLabelValues(labelValues []string) prometheus.Gauge {

	if IsEnable() {
		return gw.g.vec.WithLabelValues(makeLabelValues(labelValues)...)
	}
	return nil
}
