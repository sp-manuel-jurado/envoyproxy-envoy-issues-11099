package main

import (
	"github.com/golang/protobuf/proto"
	sp_rpc_hello "github.com/socialpoint/envoyproxy-envoy-issues-11099/service/hello/pkg/SP/Rpc/Hello"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	addr = ":8888"
)

func main() {
	http.HandleFunc("/sp.rpc.hello.HelloService/Hello", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		req := &sp_rpc_hello.HelloRequest{}
		err = proto.Unmarshal(body, req)
		if err != nil {
			panic(err)
		}

		if req.Name == "" {
			w.Header().Set("Content-Type", "application/grpc+proto")
			w.Header().Set("Content-Length", "0")
			w.Header().Set("grpc-status", "3")
			w.Header().Set("grpc-message", "invalid name")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/grpc+proto")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusOK)

		content, err := proto.Marshal(&sp_rpc_hello.HelloResponse{
			HelloMessage: "Hello " + req.Name,
		})
		if err != nil {
			panic(err)
		}

		_, err = w.Write(content)
		if err != nil {
			panic(err)
		}
	})

	log.Printf("http server in localhost:" + addr + " started")
	log.Printf("\t- route (grpc-http1-reverse-bridge): /sp.rpc.HelloService/Hello" + addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
