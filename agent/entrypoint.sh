#!/bin/sh

if [ -z $MONITOR_PID ];then
    echo "NO MONITOR_PID"
    return 1
else
    echo "MONITOR_PID = "$MONITOR_PID
fi

if [ -z $MONITOR_SERVICE ];then
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

if [ -z $MONITOR_IP ];then
    export MONITOR_IP=$(hostname -I| awk '{print $1}')
    echo "No MONITOR_IP found, set it to $MONITOR_IP"
fi

if [ -z $MONITOR_INTERVAL ];then
    echo "No MONITOR_INTERVAL, set it to 3s"
    export MONITOR_INTERVAL=3
fi

./agent -p $MONITOR_PID -l $MONITOR_INTERVAL -i $MONITOR_IP \
-s $MONITOR_SERVICE -d $REPORT_DBURL - b $REPORT_DBBUCKET -o $REPORT_DBORG -t $REPORT_DBTOKEN