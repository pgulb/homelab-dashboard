FROM --platform=$BUILDPLATFORM golang:1.26.2-alpine3.23 AS builder
ARG TARGETOS TARGETARCH
WORKDIR /src
ADD go.mod go.sum /src/
RUN go mod download
ADD cmd/ /src/
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /dist/dash .

FROM --platform=$TARGETPLATFORM scratch
COPY --from=builder /dist/dash /dash
EXPOSE 9091
CMD ["/dash"]
