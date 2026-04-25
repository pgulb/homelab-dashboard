FROM golang:1.26.2-alpine3.23 AS builder
WORKDIR /src
ADD go.mod go.sum /src/
RUN go mod download
ADD cmd/ /src/
RUN go build -o /dist/dash .

FROM scratch
COPY --from=builder /dist/dash /dash
EXPOSE 9091
CMD ["/dash"]
