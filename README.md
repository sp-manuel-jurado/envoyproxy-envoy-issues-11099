# envoyproxy-envoy-issues-11099

Test for envoy issue:
- https://github.com/envoyproxy/envoy/issues/11099
- envoy version v1.14.1 (latest stable)

Related info:
- https://www.envoyproxy.io/docs/envoy/latest/configuration/best_practices/level_two

Flow chart:

- client -> Envoy L1 (upstream HTTP2/gRPC) -> Envoy L2 (grpc-http1-reverse-bridge; upstream HTTP1.1) -> backend

Service/hello proto (empty string in name property will be an invalid name):
```
syntax = "proto3";

package sp.rpc.hello;

service HelloService {
    rpc Hello(HelloRequest) returns (HelloResponse);
}

message HelloRequest {
    string name = 1;
}

message HelloResponse {
    string hello_message = 1;
}
```

Execute tests (dependencies docker/docker-compose):
- `make up`

```
hello_service_test_1  | === RUN   TestHello_EnvoyProxyLevel2
hello_service_test_1  | ********** client -> L2 (reverse-bridge) -> backend
hello_service_test_1  | hello_message:"Hello Lucas"
hello_service_test_1  | **********
hello_service_test_1  |
hello_service_test_1  |
hello_service_test_1  | --- PASS: TestHello_EnvoyProxyLevel2 (0.01s)
hello_service_test_1  | === RUN   TestHello_EnvoyProxyLevel2InvalidArgumentName
hello_service_test_1  | ********** client -> L2 (reverse-bridge) -> backend (invalid name)
hello_service_test_1  | rpc error: code = InvalidArgument desc = invalid name
hello_service_test_1  | **********
hello_service_test_1  |
hello_service_test_1  |
hello_service_test_1  | --- PASS: TestHello_EnvoyProxyLevel2InvalidArgumentName (0.02s)
hello_service_test_1  | === RUN   TestHello_EnvoyProxyLevel1
hello_service_test_1  | ********** client -> L1 -> L2 (reverse-bridge) -> backend
hello_service_test_1  | hello_message:"Hello Lucas"
hello_service_test_1  | **********
hello_service_test_1  |
hello_service_test_1  |
hello_service_test_1  | --- PASS: TestHello_EnvoyProxyLevel1 (0.02s)
hello_service_test_1  | === RUN   TestHello_EnvoyProxyLevel1InvalidArgumentName
hello_service_test_1  | ********** client -> L1 -> L2 (reverse-bridge) -> backend (invalid name)
hello_service_test_1  | rpc error: code = Unavailable desc = upstream connect error or disconnect/reset before headers. reset reason: connection termination
hello_service_test_1  | **********
hello_service_test_1  |
hello_service_test_1  |
hello_service_test_1  | --- FAIL: TestHello_EnvoyProxyLevel1InvalidArgumentName (0.01s)
hello_service_test_1  |     main_test.go:86:
hello_service_test_1  |         	Error Trace:	main_test.go:86
hello_service_test_1  |         	Error:      	Not equal:
hello_service_test_1  |         	            	expected: "rpc error: code = InvalidArgument desc = invalid name"
hello_service_test_1  |         	            	actual  : "rpc error: code = Unavailable desc = upstream connect error or disconnect/reset before headers. reset reason: connection termination"
hello_service_test_1  |
hello_service_test_1  |         	            	Diff:
hello_service_test_1  |         	            	--- Expected
hello_service_test_1  |         	            	+++ Actual
hello_service_test_1  |         	            	@@ -1 +1 @@
hello_service_test_1  |         	            	-rpc error: code = InvalidArgument desc = invalid name
hello_service_test_1  |         	            	+rpc error: code = Unavailable desc = upstream connect error or disconnect/reset before headers. reset reason: connection termination
hello_service_test_1  |         	Test:       	TestHello_EnvoyProxyLevel1InvalidArgumentName
hello_service_test_1  | FAIL
hello_service_test_1  | FAIL	github.com/socialpoint/envoyproxy-envoy-issues-11099/service/hello	0.062s
hello_service_test_1  | ?   	github.com/socialpoint/envoyproxy-envoy-issues-11099/service/hello/pkg/SP/Rpc/Hello	[no test files]
hello_service_test_1  | FAIL
hello_service_test_1  | make: *** [Makefile:3: test] Error 1
envoyproxy-envoy-issues-11099_hello_service_test_1 exited with code 2
```

- If you want to execute tests using grpcurl:
```
cd service/hello
```

```
~ NAME=Lucas make grpc-curl-proxy-level-1
grpcurl -v -plaintext -d '{"name": "Lucas"}' -proto pkg/SP/Rpc/Hello/hello_service.proto localhost:11001 sp.rpc.hello.HelloService/Hello

Resolved method descriptor:
rpc Hello ( .sp.rpc.hello.HelloRequest ) returns ( .sp.rpc.hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Mon, 11 May 2020 10:04:48 GMT
server: envoy
x-envoy-upstream-service-time: 1

Response contents:
{
  "helloMessage": "Hello Lucas"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response



~ NAME=Lucas make grpc-curl-proxy-level-2
grpcurl -v -plaintext -d '{"name": "Lucas"}' -proto pkg/SP/Rpc/Hello/hello_service.proto localhost:10001 sp.rpc.hello.HelloService/Hello

Resolved method descriptor:
rpc Hello ( .sp.rpc.hello.HelloRequest ) returns ( .sp.rpc.hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Mon, 11 May 2020 10:04:51 GMT
server: envoy
x-envoy-upstream-service-time: 0

Response contents:
{
  "helloMessage": "Hello Lucas"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response




~ NAME= make grpc-curl-proxy-level-1
grpcurl -v -plaintext -d '{"name": ""}' -proto pkg/SP/Rpc/Hello/hello_service.proto localhost:11001 sp.rpc.hello.HelloService/Hello

Resolved method descriptor:
rpc Hello ( .sp.rpc.hello.HelloRequest ) returns ( .sp.rpc.hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
(empty)

Response trailers received:
content-type: application/grpc
date: Mon, 11 May 2020 10:04:57 GMT
server: envoy
Sent 1 request and received 0 responses
ERROR:
  Code: Unavailable
  Message: upstream connect error or disconnect/reset before headers. reset reason: connection termination
make: *** [grpc-curl-proxy-level-1] Error 1




~ NAME= make grpc-curl-proxy-level-2
grpcurl -v -plaintext -d '{"name": ""}' -proto pkg/SP/Rpc/Hello/hello_service.proto localhost:10001 sp.rpc.hello.HelloService/Hello

Resolved method descriptor:
rpc Hello ( .sp.rpc.hello.HelloRequest ) returns ( .sp.rpc.hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
(empty)

Response trailers received:
content-length: 5
content-type: application/grpc
date: Mon, 11 May 2020 10:05:02 GMT
server: envoy
x-envoy-upstream-service-time: 0
Sent 1 request and received 0 responses
ERROR:
  Code: InvalidArgument
  Message: invalid name
make: *** [grpc-curl-proxy-level-2] Error 1
```

Headers returned by http1.1 in case of grpc error:
```
w.Header().Set("Content-Type", "application/grpc+proto")
w.Header().Set("Content-Length", "0")
w.Header().Set("grpc-status", "3")
w.Header().Set("grpc-message", "invalid name")
```
