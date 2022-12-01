package metric

import "github.com/prometheus/client_golang/prometheus"

type (
	Param  prometheus.Labels
	option func(Param)
)

const (
	kindDefault  = "default"
	kindBusiness = "business"
	kindBase     = "base"
)

func (p Param) Init(opts ...option) Param {

	for _, opt := range opts {
		opt(p)
	}

	p.checkKind()
	return p
}

//用户自定义method的值，标识不通的restful请求方法
func Method(method string) option {
	return func(p Param) {
		p["method"] = method
	}
}

//用户自定义endpoint的值,标识指标发生现场，如api的path， rpc的服务名字等
func Endpoint(endpoint string) option {
	return func(p Param) {
		p["endpoint"] = endpoint
	}
}

//若用户没有指定kind,设置默认kind
func (p Param) checkKind() {

	if _, ok := p["kind"]; !ok {
		p["kind"] = kindDefault
	}
}

//用户自定义kind的值
func Kind(kind string) option {
	return func(p Param) {
		p["kind"] = kind
	}
}

//设置kind值为business, 表示业务监控
func BusinessKind() option {
	return Kind(kindBusiness)
}

//设置kind值为base, 表示系统基础监控
func BaseKind() option {
	return Kind(kindBase)
}

//初始化添加静态label
func AddLabel(k, v string) option {
	return func(p Param) {
		p[k] = v
	}
}
