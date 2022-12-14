package metric

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Observer interface {
	Observe(float64)
}
type ObserverWith interface {
	ObserveWith(float64, []string)
}

type TimeGranularity int

const (
	GranularityNanoSecond = iota
	GranularityMicroSecond
	GranularityMilliSecond
	GranularitySecond
)

// note:此数组顺序有意义，不能随意调整顺序
var gLabelNames []string = []string{
	"app",
	"center",
	"suffix",
}

type metricManager struct {
	appCenter string
	appSuffix string
	appName   string /* appName = center+"-"+suffix */
	enable    bool

	engine *gin.Engine
	addr   string
}

var (
	gMetricManager *metricManager
)

func init() {

	gMetricManager = &metricManager{
		appCenter: "default",
		appName:   "default",
		enable:    true,
	}
}

func Init(addr string, appCenter string, appSuffix ...string) error {

	if gMetricManager.isRunning() {
		return fmt.Errorf("monitor is already running at addr:%s", gMetricManager.addr)
	}

	err := gMetricManager.updateAppName(appCenter, appSuffix...)
	if err != nil {
		return err
	}

	gMetricManager.asyncRun(addr)

	return nil
}

/*
	记录跟踪一个函数的执行时间（单位:毫秒ms)
	必须通过defer关键字调用
	如要跟踪test函数的处理耗时,如下:
	func test {
		defer ObserveCostTime(obs, time.Now())
		//process
		time.Sleep(time.Second)
	}
*/
func ObserveCostTime(obs Observer, now time.Time) {

	costTime := time.Since(now).Milliseconds()
	obs.Observe(float64(costTime))
}

func ObserveCostTimeWith(obs ObserverWith, labelValues []string, now time.Time) {

	costTime := time.Since(now).Milliseconds()
	obs.ObserveWith(float64(costTime), labelValues)
}

// 可自定义要统计的时间 粒度，支持nano,micro,milli,and 1
func ObserveCostTimeGranularity(obs Observer, now time.Time, granularity TimeGranularity) {

	costTime := time.Since(now)
	dur := transGranularity(costTime, granularity)
	obs.Observe(float64(dur))
}

func ObserveCostTimeGranularityWith(obs ObserverWith, labelValues []string, now time.Time, granularity TimeGranularity) {

	costTime := time.Since(now)
	dur := transGranularity(costTime, granularity)
	obs.ObserveWith(float64(dur), labelValues)
}

func IsEnable() bool {

	return gMetricManager.enable
}

func SwitchEnable(isEnable bool) {

	gMetricManager.enable = isEnable
}

func (m *metricManager) isRunning() bool {

	return m.engine != nil
}

func (m *metricManager) asyncRun(addr string) {

	m.addr = addr
	m.engine = gin.Default()
	m.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	go m.engine.Run(addr)
}

func (m *metricManager) updateAppName(appCenter string, appSuffix ...string) error {

	if appCenter == "" {
		return errors.New("app name is unvalid")
	} else if 50 < len(appCenter) {
		return errors.New("app name lenght must <=50")
	}
	m.appCenter = appCenter

	if 0 < len(appSuffix) && 0 < len(appSuffix[0]) {
		m.appSuffix = appSuffix[0]
		m.appName = m.appCenter + "-" + m.appSuffix
	} else {
		m.appSuffix = ""
		m.appName = m.appCenter
	}

	return nil
}

func AppCenter() string {

	return gMetricManager.appCenter
}

func AppSuffix() string {

	return gMetricManager.appSuffix
}

func AppName() string {

	return gMetricManager.appName
}

func makeLabelNames(labelNames []string) []string {

	lns := make([]string, 0, len(gLabelNames)+len(labelNames))
	lns = append(lns, gLabelNames...)
	lns = append(lns, labelNames...)
	return lns
}

func makeLabelValues(labelValues []string) []string {

	lvs := make([]string, 0, len(gLabelNames)+len(labelValues))
	lvs = append(lvs, AppName(), AppCenter(), AppSuffix())
	lvs = append(lvs, labelValues...)
	return lvs
}

func transGranularity(costTime time.Duration, granularity TimeGranularity) int64 {
	var dur int64
	switch granularity {
	case GranularityNanoSecond:
		dur = costTime.Nanoseconds()
	case GranularityMicroSecond:
		dur = costTime.Microseconds()
	case GranularityMilliSecond:
		dur = costTime.Milliseconds()
	case GranularitySecond:
		dur = int64(costTime.Seconds())
	}
	return dur
}
