#!/bin/sh

if [ -z $NACOS_IP ];then
    echo "NO NACOS_IP"
    return 1
else 
    echo "NACOS_IP = "$NACOS_IP
fi

if [ -z $NACOS_PORT ];then
    echo "NO NACOS_PORT, set it to 8848"
    export NACOS_PORT=8848
else 
    echo "NACOS_PORT = "$NACOS_PORT
fi

if [ -z $NACOS_NS ];then
    echo "NO NACOS_NS, set it to public"
    export NACOS_NS=public
else 
    echo "NACOS_NS = "$NACOS_NS
fi

if [ -z $NACOS_GROUP ];then
    echo "NO NACOS_GROUP, set it to DEFAULT_GROUP"
    export NACOS_GROUP=DEFAULT_GROUP
else 
    echo "NACOS_GROUP = "$NACOS_GROUP
fi

if [ -z $MONITOR_SERVICES ];then
    echo "NO MONITOR_SERVICES"
else 
    echo "MONITOR_SERVICES = "$MONITOR_SERVICES
fi

if [ -z $REPORT_DBURL ];then
    echo "NO REPORT_DBURL"
    return 1
else 
    echo "REPORT_DBURL = "$REPORT_DBURL
fi

if [ -z $REPORT_DBBUCKET ];then
    echo "NO REPORT_DBBUCKET"
    return 1
else 
    echo "REPORT_DBBUCKET = "$REPORT_DBBUCKET
fi

if [ -z $REPORT_DBORG ];then
    echo "NO REPORT_DBORG"
    return 1
else 
    echo "REPORT_DBORG = "$REPORT_DBORG
fi

if [ -z $REPORT_DBTOKEN ];then
    echo "NO REPORT_DBTOKEN"
    return 1
else 
    echo "REPORT_DBTOKEN = "$REPORT_DBTOKEN
fi

if [ -z $MONITOR_INTERVAL ];then
    echo "No MONITOR_INTERVAL, set it to 3s"
    export MONITOR_INTERVAL=3
fi

echo "./gomonitor-manager -i $NACOS_IP -p $NACOS_PORT -n $NACOS_NS -s $MONITOR_SERVICES -g $NACOS_GROUP -d $REPORT_DBURL -b $REPORT_DBBUCKET -o $REPORT_DBORG -t $REPORT_DBTOKEN -l $MONITOR_INTERVAL"

./gomonitor-manager -i $NACOS_IP -p $NACOS_PORT -n $NACOS_NS -s $MONITOR_SERVICES -g $NACOS_GROUP -d $REPORT_DBURL -b $REPORT_DBBUCKET -o $REPORT_DBORG -t $REPORT_DBTOKEN -l $MONITOR_INTERVAL