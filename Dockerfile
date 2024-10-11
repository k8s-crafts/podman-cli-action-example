FROM golang:1.22-alpine as builder

WORKDIR /workspace

COPY go.mod go.sum ./

# cache deps before building
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -o app main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest as runner

WORKDIR /workspace

COPY --from=builder /workspace/app .

EXPOSE 8080

ENTRYPOINT [ "/workspace/app" ]
