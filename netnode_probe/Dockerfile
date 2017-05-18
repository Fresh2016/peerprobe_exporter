FROM       docker.io/alpine:latest
#FROM       docker.io/golang:alpine
#FROM        alpine:3.5
MAINTAINER shangaoshuichang2017 <just as shangaoshuichang2017's app>

ENV         LISTEN_PORT=${LISTEN_PORT:-'9009'} \
            NODE_EXPORTER_LISTEN_PORT=${NODE_EXPORTER_LISTEN_PORT:-'9000'}

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main" > /etc/apk/repositories

RUN apk add --update curl bash && \
    rm -rf /var/cache/apk/*

#RUN apk add --update curl git && \
#    rm -rf /var/cache/apk/*

#RUN mkdir -p /opt/gopath/src/github.com/prometheus/node_exporter && \
#    mkdir -p /opt/gopath/src/github.com/prometheus/netnode_probe

RUN  mkdir -p /opt/bin

#COPY /node_exporter/peerprobe_exporter /opt/gopath/src/github.com/prometheus/node_exporter/
#COPY /node_exporter/start_peerprobe_exporter.sh /opt/gopath/src/github.com/prometheus/node_exporter/
#COPY /netnode_probe/netnode_probe /opt/gopath/src/github.com/prometheus/netnode_probe/

COPY /node_exporter/peerprobe_exporter /opt/bin/
COPY /node_exporter/start_peerprobe_exporter.sh /opt/bin/
COPY /netnode_probe/netnode_probe /opt/bin/


#RUN chmod a+x /opt/gopath/src/github.com/prometheus/node_exporter/peerprobe_exporter
#RUN chmod a+x /opt/gopath/src/github.com/prometheus/node_exporter/peerprobe_exporter && \
#    chmod a+x /opt/gopath/src/github.com/prometheus/node_exporter/start_peerprobe_exporter.sh && \
#    chmod a+x /opt/gopath/src/github.com/prometheus/netnode_probe/netnode_probe 

RUN chmod a+x /opt/bin/peerprobe_exporter && \
    chmod a+x /opt/bin/start_peerprobe_exporter.sh && \
    chmod a+x /opt/bin/netnode_probe 

EXPOSE      ${LISTEN_PORT}

#ENTRYPOINT  [ "/opt/gopath/src/github.com/prometheus/node_exporter/start_peerprobe_exporter.sh" ]
#ENTRYPOINT  [ "/opt/gopath/src/github.com/prometheus/netnode_probe/netnode_probe" ]
#ENTRYPOINT  [ "/opt/gopath/src/github.com/prometheus/node_exporter/peerprobe_exporter" ]
#ENTRYPOINT  [ "/opt/bin/netnode_probe" ]
CMD  [ "sh","/opt/bin/start_peerprobe_exporter.sh" ]
