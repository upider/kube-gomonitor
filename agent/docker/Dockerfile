FROM 1445277435/centos:gopacket AS builder
WORKDIR /workspace
COPY ./ src/
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
RUN cd src && go build -o kube-gomonitor-agent agent/agent.go

FROM centos:7
WORKDIR /workspace
RUN yum install -y deltarpm epel-release && yum install -y libpcap-devel
COPY --from=builder /workspace/src/kube-gomonitor-agent .
ENTRYPOINT ./kube-gomonitor-agent --monitorPid $MONITOR_PID --monitorIP $MONITOR_IP --monitorService $MONITOR_SERVICE \
            --dburl $REPORT_DBURL --bucket $REPORT_DBBUCKET --organization $REPORT_DBORG --token $REPORT_DBTOKEN \
            --monitorInterval $MONITOR_INTERVAL --nacosIPs $NACOS_IP --nacosPort $NACOS_PORT