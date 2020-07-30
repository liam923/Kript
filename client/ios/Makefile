setup:
	go get \
			github.com/gogo/googleapis \
			github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
			github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go mod vendor

	mkdir -p docs/proto
	cp -r ../../proto docs

	mkdir -p proto
	protoc \
			-I docs/proto \
			-I vendor/github.com/grpc-ecosystem/grpc-gateway/ \
			-I vendor/github.com/gogo/googleapis/ \
			--swift_out=proto \
			--grpc-swift_out=Client=true,Server=false:proto \
			docs/proto/kript/api/*.proto
			
	pod install
