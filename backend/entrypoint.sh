#!/bin/sh

if [ -z $MONITOR_SERVICES ];then
    echo "NO MONITOR_SERVICE"
    return 1
else 
    echo "MONITOR_SERVICE = "$MONITOR_SERVICE
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

echo "./gominotor-manager -i $MONITOR_IP -s $MONITOR_SERVICES -d $REPORT_DBURL -b $REPORT_DBBUCKET -o $REPORT_DBORG -t $REPORT_DBTOKEN -l $MONITOR_INTERVAL"

./gominotor-manager -l $MONITOR_INTERVAL -s $MONITOR_SERVICES -d $REPORT_DBURL -b $REPORT_DBBUCKET -o $REPORT_DBORG -t $REPORT_DBTOKEN