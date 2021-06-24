# Build the manager binary
FROM golang:1.16 as builder
WORKDIR /workspace/
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
COPY ./ src/
RUN cd src && go build -o kube-gomonitor-backned backend/backend.go

FROM centos:7
WORKDIR /
COPY --from=builder /workspace/src/kube-gomonitor-backned .
USER 65532:65532
ENTRYPOINT ./kube-gomonitor-backned --nacosIPs $NACOS_IP --nacosPort $NACOS_PORT --agentImage $AGENT_IMAGE \
            --services $SERVICES --groups $GROUPS --namespaces $NAMESPACES \
            --dburl $REPORT_DBURL --bucket $REPORT_DBBUCKET --organization $REPORT_DBORG --token $REPORT_DBTOKEN \
            --interval $MONITOR_INTERVAL