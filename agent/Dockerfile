FROM centos:gopacket AS builder
WORKDIR /workspace
COPY ./ src/
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
RUN cd src && go build -o gomonitor-agent agent/agent.go

FROM centos:7
WORKDIR /
RUN yum install -y deltarpm epel-release && yum install -y libpcap-devel
COPY --from=builder /workspace/src/gomonitor-agent .
COPY agent/entrypoint.sh .
RUN chmod +x entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]