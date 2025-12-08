module github.com/practical6/api-gateway

go 1.23

require (
	github.com/gorilla/mux v1.8.1
	github.com/practical6/proto v0.0.0
	google.golang.org/grpc v1.65.0
)

require (
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240624140628-dc46fd24d27d // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/practical6/proto => ../proto
