# gomonitor

## Description

monitor a process and store message, \
no matter on bare metal or kubernetes

## Agent
1. 启动入口`agent/entrypoint.sh`
2. 容器中运行需要设置环境变量：
   - `MONITOR_PID` 要监控的进程号
   - `MONITOR_SERVICE` 要监控的service名
   - `REPORT_DBURL` 数据库url (influxdb)
   - `REPORT_DBBUCKET` influxdb buket
   - `REPORT_DBORG` influxdb organization
   - `REPORT_DBTOKEN` influxdb token
   - `MONITOR_IP` 要监控进程的ip
   - `MONITOR_INTERVAL` 监控时间间隔 (second)
3. 构建docker容器命令：
   ```sh
    cd agent
    docker build -t image:tag -f Dockerfile.yaml ../
   ```

## Manager
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
3. 启动入口`backend/entrypoint.sh`
4. 容器中运行需要设置环境变量：
    - `MONITOR_SERVICES` 要监控的服务名（k8s环境不需要设置）
    - `NACOS_IP` nacos服务器ip地址（k8s环境不需要设置）
    - `NACOS_PORT` nacos服务器port（默认8848，k8s环境不需要设置）
    - `NACOS_NS` nacos namespaceid（默认public，k8s环境不需要设置）
    - `NACOS_GROUP` nacos group（默认DEFAULT_GROUP，k8s环境不需要设置）
    - `REPORT_DBURL` influxdb url（http://......）
    - `REPORT_DBBUCKET` influxdb bucket
    - `REPORT_DBORG` influxdb organization
    - `REPORT_DBTOKEN` influxdb token
    - `MONITOR_INTERVAL` 监控采样间隔（默认3s）
5. 构建docker容器命令：
   ```sh
    cd backend
    docker build -t image:tag -f Dockerfile.yaml ../
   ```
6. kubernetes 部署文件
    ```
    backend/deploy.yaml
    ```
## TODO

- [x] agent cmdline flags
- [x] monitor message storage (influxdb now)
- [x] backend monitor service
- [ ] k8s environment

## Reference
- _[docker sdk reference](https://docs.docker.com/engine/api/v1.41/)_
- _[kubebuilder reference](https://book.kubebuilder.io/)_
- _[nacos doc](https://nacos.io/zh-cn/docs/quick-start.html)_
- _[k8s controller example](https://github.com/kubernetes-sigs/controller-runtime/tree/master/examples/builtins)_
- _[telegraf-operator](https://github.com/influxdata/telegraf-operator)_