pb:
	@rm -rf ./protobuf/*.go && \
	protoc --proto_path=./protobuf \
	--go_out=./protobuf/ --go_opt=paths=source_relative \
	--go-grpc_out=./protobuf --go-grpc_opt=paths=source_relative \
	./protobuf/*.proto && \
	echo "proto code generation successful"