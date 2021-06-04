# Build the manager binary
FROM golang:1.16 as builder
WORKDIR /workspace/
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
COPY ./ src/
RUN cd src && go build -o gomonitor-manager backend/backend.go

FROM centos:7
WORKDIR /
COPY --from=builder /workspace/src/gomonitor-manager .
COPY backend/entrypoint.sh .
RUN chmod +x entrypoint.sh
USER 65532:65532
ENTRYPOINT ["./entrypoint.sh"]
