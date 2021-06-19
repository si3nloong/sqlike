proto:
	@rm -rf ./x/proto/*.go && \
	protoc --proto_path=./x/proto \
	--go_out=./x/proto/ --go_opt=paths=source_relative \
	--go-grpc_out=./x/proto --go-grpc_opt=paths=source_relative \
	./x/proto/*.proto && \
	echo "proto code generation successful"