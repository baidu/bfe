# 性能数据采集

BFE内置了CPU profile接口，配合使用火焰图工具，用于定位分析性能问题

## 配置

使用和监控项采集相同的端口
```
[server]
monitorPort = 8421
```

## 工具准备

* FlameGragh

```
git clone https://github.com/brendangregg/FlameGraph
```

其中包含stackcollpase-go.pl和flamegraph.pl

## 操作步骤

* 获取性能采样数据
```
go tool pprof -seconds=60 -raw -output=bfe.pprof  http://<addr>:<port>/debug/pprof/profile
```
注：seconds=60 表示抓取60s的采样数据

* 转换并绘制火焰图

```
./stackcollpase-go.pl bfe.pporf > bfe.flame
./flamegraph.pl bfe.flame > bfe.svg
```

在浏览器中打开bfe.svg查看
