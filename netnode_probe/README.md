# 功能简介
**netnode_probe.go**设计为网络节点+端口的探测功能。当前设计每10秒轮询监测一次/etc/hosts文件中定义的所有主机节点，并记录所有轮询监测结果；另外超过5分钟后，监测记录会老化，将被自动删除。

netnode_probe.go将从/etc/hostsmap文件中（需要预先部署此文件）读取需要进行网络连通性检测的IP与端口信息，并将检测结果记录到/etc/proberesultnew文件中（如果文件不存在，程序会自动创建）。netnode_exporter程序会从/etc/proberesultnew文件中读取检测结果，并按照IP+PORT分别统计计算连通率与连通平均时延，然后通过node_peerconnect_success_rate与node_peerconnect_timedelay_average上报给promethues>系统

# 开发过程
## 1，开发目录
在目录/opt/gopath/src/github.com/prometheus/netnode_probe下，准备探针源代码netnode_probe.go
## 2，hostmap文件要求
从/etc/hostsmap文件中读取需要探测查询的IP/PORT等信息，hostsmap文件内的格式如下，IP/TYPE/PORT信息之间用“空格”或者“tab”键隔离。
## 3，ubuntu系统build验证
用Makefile进行go build，进行代码功能验证
## 4，alpine系统编译
功能验证ok后，制作容器Dockerfile进行go源代码编译，生成二进制代码（执行docker build -t lily/go-build:1.0 . 制作go-build镜像）

docker run -it -v /opt/gopath/src/github.com/prometheus/netnode_probe:/opt/gopath/src/github.com/prometheus/netnode_probe --rm lily/go-build:1.0 /bin/bash
（#cd /opt/gopath/src/github.com/prometheus/netnode_probe）
（#go build netnode_probe.go）
## 5，alpine镜像制作
制作新Dockerfile生成应用镜像，将alpine 和gonetnode_probe打包
docker build -t shangaoshuichang2017/go-netnode-probe:1.0 .
## 6，应用容器启动
启动容器：docker run -d --volume=/etc:/etc --net="host" netnode_probe_exe
docker run -d --volume=/etc:/etc  --net="host" --name go-netprobe shangaoshuichang2017/go-netnode-probe:1.0 

