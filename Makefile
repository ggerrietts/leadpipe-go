bin:
	go build ./cmd/collect
	go build ./cmd/process
	go build ./cmd/insert

gen:
	protoc -I="api/protobuf" --go_out="internal/pb" api/protobuf/leadpipe.proto
