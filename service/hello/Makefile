test:
	sleep 10
	go clean -testcache && go test -v ./...
.PHONY: test

grpc-curl-proxy-level-2:
	grpcurl -v -plaintext -d '{"name": "${NAME}"}' -proto pkg/SP/Rpc/Hello/hello_service.proto localhost:10001 sp.rpc.hello.HelloService/Hello
.PHONY:

grpc-curl-proxy-level-1:
	grpcurl -v -plaintext -d '{"name": "${NAME}"}' -proto pkg/SP/Rpc/Hello/hello_service.proto localhost:11001 sp.rpc.hello.HelloService/Hello
.PHONY:
