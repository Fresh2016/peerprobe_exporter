peerprobe_exporter功能描述
go-peerprobe-exporter镜像实现了两个功能：一方面实现了间隔10s对目录/etc/hostsmap文件中定义的主机ip+端口进行连通性检测，并记录是否连通以及访问时延信息，信息保存超过5分钟会自动老化删除；另一方面实现了将检测结果进行汇总统计，可以通过node_peerconnect_success_rate（连通率，取值0~1）和node_peerconnect_timedelay_average（时延平均值，单位：ms，如果不通取值为-1）两个metric访问，以便prometheus监控系统进行结果查询以及后续处理。

镜像获取
docker pull shangaoshuichang2017/go-peerprobe-exporter

镜像启动
docker run -d --volume=/etc:/etc --env NODE_EXPORTER_LISTEN_PORT=8000 --net="host" --name go-netprobe shangaoshuichang2017/go-peerprobe-exporter:1.0

注意：

1，必须在主机/etc目录下提供hostsmap文件，格式如下（信息之间可以用空格或者tab键分隔）：

App Configuration:

############################################################################
CheckInterval(default:10s):   15     RecordObsoleteTime(default:5minutes):   6     DialWaitTime(default:3s):   5  (The discription and the value must be seperated by spacebar or tab key)
############################################################################



IP TYPE PORT

127.0.0.1 tcp 22

127.0.0.1 tcp 8058

127.0.0.1 udp 67

192.168.122.1 udp 53

192.168.200.14 tcp 22

2，如果不配置NODE_EXPORTER_LISTEN_PORT，缺省访问端口为9000

exporter访问
服务IP:9000（或者NODE_EXPORTER_LISTEN_PORT）
