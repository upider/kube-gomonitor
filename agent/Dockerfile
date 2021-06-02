FROM centos:gopacket AS dev
WORKDIR /workspace
COPY agent/ agent/
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
RUN export GOPROXY=https://goproxy.cn && export GO111MODULE=on
RUN go mod tidy && go build -o gomonitor-agent agent/agent.go

FROM centos:7 AS prod
WORKDIR /
RUN yum install -y deltarpm epel-release && yum install -y libpcap-devel
COPY --from=dev /workspace/gomonitor-agent .
COPY agent/entrypoint.sh .
RUN chmod +x entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]