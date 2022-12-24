FROM golang:1.19.3-bullseye

# Install tools
RUN go install github.com/cweill/gotests/gotests@v1.6.0
RUN go install github.com/fatih/gomodifytags@v1.16.0
RUN go install github.com/josharian/impl@v1.1.0
RUN go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/ramya-rao-a/go-outline@v0.0.0-20210608161538-9736a4bde949

# Install gRPC
RUN apt update
RUN apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN export PATH="$PATH:$(go env GOPATH)/bin"

# Install .proto file formatter
RUN apt install -y clang-format