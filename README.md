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

构建容器命令：
```sh
cd agent
docker build -t image:tag -f dockerfile.yaml ../
```

## server
1. 裸机环境和k8s环境实现原理不同，裸机环境依赖nacos，应用需要在nacos注册自己在对应的服务名下，并且把pid写在元数据里
2. k8s环境中，agent需要检测挂载的配置文件(/tmp/monitor-config.json)，所以应用需要自己创建此文件
   文件格式：
   ```json
    {
        "pid": 1111,
        "serviceName": "sssss",
        "ip": "12.13.41.11",
    }
   ```

## TODO

- [x] agent cmdline flags
- [x] monitor message storage (influxdb now)
- [x] backend monitor service
- [ ] k8s environment

## docker sdk reference
[docker sdk reference](https://docs.docker.com/engine/api/v1.41/)