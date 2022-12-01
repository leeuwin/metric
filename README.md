
## 指标类型
- counter
  计数值，整数，永远递增(可通过取两个时刻的值相减得到一段时间内的计数值)
- gauge
  状态值, 支持浮点数
- summary
  统计值，支持分位统计，不支持合并; 当前配置了50%,90%,95%,99%份位
- histogram
  统计值, 支持自定义分组区间统计，支持合并；
  

## 特性
- 指标会内置强制携带app和kind的label
    - 其中app支持后置补充设置（即可以先定义你指标，后提供app的值，也可以成功上报, 通过InitMonitor, UpdateAppName接口设置)
    - 但kind需要在指标定义的时候及时提供，有封闭的base和business以及默认的default可选，也可以开放完全自定义
- 指标支持自定义添加若干组label， NewCounter("metric", AddLabel("k", "v"), AddLabel("k2","v2"))
- 支持动态启停指标更新SwitchEnable()

## 使用三部曲
### 静态label
1. 初始化app与暴露地址，如：metric.Init(":8088", "appname")
2. 定义指标，如：cnt := metric.NewCounter("metric")
3. 更新指标，如：cnt.Inc()
ps:上述1，2顺序可随意切换，只要3在2之后，就ok

### 动态态label
- 例1
    1. 初始化app与暴露地址，如：metric.Init(":8088", "appname")
    2. 定义指标，如：cnt := metric.NewCounterWith("metric", "mylabel")
    3. 更新指标，如：cnt.IncWith("value")

- 例2
    1. 初始化app与暴露地址，如：metric.Init(":8088", "appname")
    2. 定义指标，如：cnt := metric.NewCounterWith("metric", "mylabel01", "mylabel02")
    3. 更新指标，如：cnt.IncWith("value01", "value02")

## 注意事项
- **metric必须保证唯一，否则创建指标metric.NewXXX()会panic抛出错误**
