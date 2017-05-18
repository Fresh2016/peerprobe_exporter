#!/bin/bash
# This is our first script.

#while true; do
# date
# sleep 30
#done

# 1,Start node peer probing
/opt/bin/netnode_probe &

#2,Start probe exporter to connect with prometheus 
/opt/bin/peerprobe_exporter -web.listen-address 0.0.0.0:${NODE_EXPORTER_LISTEN_PORT} 

# 1,Start node peer probing
#/opt/bin/netnode_probe $@
