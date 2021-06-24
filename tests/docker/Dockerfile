FROM golang:1.16 as builder
WORKDIR /workspace
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
COPY ./ src/
RUN cd src/ && go build -o gomonitor-testprog tests/registry.go tests/busy_loop.go tests/test_prog.go

FROM centos:7
WORKDIR /
COPY --from=builder /workspace/src/gomonitor-testprog .
USER 65532:65532
ENTRYPOINT ["./gomonitor-testprog"]
