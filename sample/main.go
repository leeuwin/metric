package main

import (
	"fmt"
	"time"

	"github.com/leeuwin/metric"
)

var (
	//counter
	cnt0 = metric.NewCounter("cnt_default")
	cnt1 = metric.NewCounter("cnt_base", metric.BaseKind())
	cnt2 = metric.NewCounter("cnt_business", metric.BusinessKind())
	cnt3 = metric.NewCounter("cnt_custome", metric.Kind("custome"))
	cnt4 = metric.NewCounter("cnt_default_label", metric.AddLabel("l1", "v1"), metric.Endpoint("/user/login"))
)

var (
	//counter with label dynamic value
	cntWith0 = metric.NewCounterWith("cnt_with", []string{"method", "endpoint"})
)

var (
	//gauge
	gg1 = metric.NewGauge("gg_default")
)

var (
	//histogram
	hg1 = metric.NewHistogram("hg_statistic_default", metric.BucketSchemaDefault)
	hg2 = metric.NewHistogram("hg_statistic_100", metric.BucketSchemaHundred)
	hg3 = metric.NewHistogram("hg_statistic_1000", metric.BucketSchemaThousand)
	hg4 = metric.NewHistogram("hg_statistic_10000", metric.BucketSchema10Thousand)
)

var (
	//summary
	sm1 = metric.NewSummary("sm_statistic_millicecond")
	sm2 = metric.NewSummary("sm_statistic_second")
)

func main() {

	// 访问暴露指标命令:curl 0:8088/metrics
	err1 := metric.Init("0.0.0.0:8088", "user-center", "api")
	if err1 != nil {
		fmt.Println(err1)
	}

	//演示重复初始化报错
	err2 := metric.Init("0.0.0.0:8088", "user-center", "srv")
	if err2 != nil {
		fmt.Println(err2)
	}

	//count with
	cntWith0.IncWith([]string{"GET", "/user/login"})
	cntWith0.IncWith([]string{"GET", "/user/login"})
	cntWith0.IncWith([]string{"POST", "/user/login"})
	cntWith0.IncWith([]string{"POST", "/user/logout"})

	//统计匿名函数执行耗时（ms)
	func() {

		defer metric.ObserveCostTime(sm1, time.Now())

		//模拟处理耗时2s
		time.Sleep(2 * time.Second)
	}()

	//统计匿名函数执行耗时（s)
	func() {

		defer metric.ObserveCostTimeGranularity(sm2, time.Now(), metric.GranularitySecond)

		//模拟处理耗时2s
		time.Sleep(2 * time.Second)
	}()

	//可根据系统配置 or 通过接口调用，切换指标是否激活有效更新,默认是 激活可用状态
	metric.SwitchEnable(true)

	//counter
	cnt5 := metric.NewCounter("cnt_default_label2", metric.AddLabel("l1", "v1"), metric.AddLabel("l2", "v2"))
	go func() {
		for {
			cnt0.Inc()
			cnt1.Inc()
			cnt2.Inc()
			cnt3.Inc()
			cnt4.Inc()
			cnt5.Inc()

			time.Sleep(1 * time.Second)
		}
	}()

	//gauge
	go func() {
		i := 0.0
		for {
			i += 1.0
			gg1.Set(i)
			time.Sleep(1 * time.Second)
		}
	}()

	//histogram
	go func() {

		nums := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		i := 0

		for {
			hg1.Observe(nums[i%len(nums)])
			hg2.Observe(nums[i%len(nums)])
			hg3.Observe(nums[i%len(nums)])
			hg4.Observe(nums[i%len(nums)])

			time.Sleep(2 * time.Second)
			i++
		}
	}()

	//summary
	go func() {

		nums := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		i := 0

		for {
			sm1.Observe(nums[i%len(nums)])

			time.Sleep(2 * time.Second)
			i++
		}

	}()

	fmt.Printf("monitor switch:%v\n", metric.IsEnable())
	time.Sleep(10 * time.Second)
	metric.SwitchEnable(false)
	fmt.Printf("monitor switch:%v\n", metric.IsEnable())
	time.Sleep(10 * time.Second)
	metric.SwitchEnable(true)
	fmt.Printf("monitor switch:%v\n", metric.IsEnable())

	select {}
}
