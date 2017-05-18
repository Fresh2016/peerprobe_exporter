# 功能简介
**peerconn.go**基于https://github.com/prometheus/node_exporter 的node_exporter进行功能增强，读取netnode_probe.go写入/etc/proberesultnew文件中的IP/PORT检测信息，通过增加的node_peerconnect_success_rate与node_peerconnect_timedelay_seconds两个metric，把通过netnode_probe监测到的每个网络节点连通状态信息进行汇总、计算后，将结果与Promethues对接。

# 开发过程
## 1，   开发工作目录
/opt/gopath/src/github.com/prometheus/
## 2，   node_exporter源代码
go get -u github.com/prometheus/node_exporter （下载代码），基于node_exporter.go代码修改为peerprobe_exporter.go，新增peerconn.go相关的defaultCollectors
## 3，   新增peerconn.go
在/opt/gopath/src/github.com/prometheus/node_exporter/collector 下新增peerconn.go，实现新增node_peerconnect_success_rate与node_peerconnect_timedelay_average的功能
## 4，   修改Makefile
使GOPATH/GOROOT与环境一致
## 5，   编译、试运行
make编译go代码
## 6，   准备探测与exporter二合一容器镜像
在/opt/gopath/src/github.com/prometheus 目录下新增Dockerfile，在/opt/gopath/src/github.com/prometheus/node_exporter 新增start_peerprobe_exporter.sh，
修改makefile的镜像tag为node-exporter-ext，make docker制作容器镜像
## 7，   制作镜像、运行应用容器
docker build -t shangaoshuichang2017/go-peerprobe-exporter:1.0 .

docker run -d --volume=/etc:/etc  --net="host" --name go-netprobe shangaoshuichang2017/go-peerprobe-exporter:1.0
## 8，访问exporter
以服务IP：9000访问peerprobe-exporter，可以查询新增node_peerconnect_success_rate与node_peerconnect_timedelay_average的功能
