FROM golang:latest

RUN pwd; ls -l;
RUN git clone https://github.com/k2rth1k/grpc_learning.git
WORKDIR grpc_learning
RUN pwd; ls -l;
RUN go mod tidy
RUN make run_server

EXPOSE 50051