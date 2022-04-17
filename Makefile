bin:
	go build ./cmd/collect
	go build ./cmd/process
	go build ./cmd/insert

gen:
	protoc -I="protobuf" --go_out="pkg" protobuf/leadpipe.proto
