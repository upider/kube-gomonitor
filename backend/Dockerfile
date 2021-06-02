# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on

# Copy the go source
COPY ./ src/

# Build
RUN cd src && go build -o gominotor-manager backend/backend.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# FROM gcr.io/distroless/static:nonroot
FROM alpine:3
WORKDIR /
COPY --from=builder /workspace/src/gominotor-manager .
COPY backend/entrypoint.sh .
RUN chmod +x entrypoint.sh
USER 65532:65532
ENTRYPOINT ["./entrypoint.sh"]
