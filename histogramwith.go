package metric

import "github.com/prometheus/client_golang/prometheus"

type histogramWith struct {
	h histogram
}

//带自定义label的指标，必须配套调用with的操作且提供分别匹配的value值
func NewHistogramWith(metric string, bucketSchema int, labelNames []string, opts ...option) histogramWith {

	return histogramWith{
		h: histogram{
			vec: newHistogramVec(metric, bucketSchema, makeLabelNames(labelNames), opts...),
		},
	}
}

func (hw histogramWith) ObserveWith(v float64, labelValues []string) {

	if obsv := hw.getCollectorWithLabelValues(labelValues); obsv != nil {
		obsv.Observe(v)
	}
}

func (hw histogramWith) getCollectorWithLabelValues(labelValues []string) prometheus.Observer {

	if IsEnable() {
		return hw.h.vec.WithLabelValues(makeLabelValues(labelValues)...)
	}
	return nil
}
