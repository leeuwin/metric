package metric

import "github.com/prometheus/client_golang/prometheus"

type summaryWith struct {
	s summary
}

//带自定义label的指标，必须配套调用with的操作且提供分别匹配的value值
func NewSummaryWith(metric string, labelNames []string, opts ...option) summaryWith {

	return summaryWith{
		s: summary{
			vec: newSummaryVec(metric, makeLabelNames(labelNames), opts...),
		},
	}
}

func (sw summaryWith) ObserveWith(v float64, labelValues []string) {

	if obsv := sw.getCollectorWithLabelValues(labelValues); obsv != nil {
		obsv.Observe(v)
	}
}

func (sw summaryWith) getCollectorWithLabelValues(labelValues []string) prometheus.Observer {

	if IsEnable() {
		return sw.s.vec.WithLabelValues(makeLabelValues(labelValues)...)
	}
	return nil
}
