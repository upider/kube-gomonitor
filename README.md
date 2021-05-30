# gomonitor

## description

monitor a process and store message, \
no matter on bare metal or kubernetes

## agent
监控agent, 使用docker, 启动入口`entrypoint.sh` \
需要设置环境变量：
1. `MONITOR_PID` 要监控的进程号
2. `MONITOR_SERVICE` 要监控的service名
3. `REPORT_DBURL` 数据库url (influxdb)
4. `REPORT_DBBUCKET` influxdb buket
5. `REPORT_DBORG` influxdb organization
6. `REPORT_DBTOKEN` influxdb token
7. `MONITOR_IP` 要监控进程的ip
8. `MONITOR_INTERVAL` 监控时间间隔 (second)

## TODO

- [x] agent cmdline flags
- [x] monitor message storage (influxdb now)
- [ ] backend monitor service
- [ ] k8s environment

## docker sdk reference
[docker sdk reference](https://docs.docker.com/engine/api/v1.41/)