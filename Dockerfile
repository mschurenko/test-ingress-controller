FROM golang:1.8

RUN go get -v k8s.io/client-go/...
COPY main.go /opt
WORKDIR /opt
RUN go build -o ingress_ctrl main.go && chmod +x ingress_ctrl
