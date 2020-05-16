install:
	go get \
		github.com/gogo/protobuf/protoc-gen-gogo \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/golang/mock/mockgen

vendor:
	go mod vendor

generate:
	mkdir -p pkg/proto
	mkdir -p docs/api
	mkdir -p docs/proto

	# Copy proto files and add go_package to them.
	cp -r ../proto docs
	for proto in $$(find docs/proto -type f -name "*.proto"); do \
  		export package=$$(echo "$$proto" | sed -E 's#\docs/(.*)/[^/]+\.proto#github.com/liam923/Kript/server/pkg/\1#'); \
  		sed -i.bak -E "s#(package .+;)#\1 option go_package = \"$${package}\";#" $$proto; \
  		rm "$$proto".bak; \
  	done

  	# Compile proto files.
	protoc \
		-I docs/proto \
		-I vendor/github.com/grpc-ecosystem/grpc-gateway/ \
		-I vendor/github.com/gogo/googleapis/ \
		-I vendor/ \
		--gogo_out=plugins=grpc,paths=source_relative,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:\
./pkg/proto/ \
		--grpc-gateway_out=allow_patch_feature=false,paths=source_relative,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:\
./pkg/proto \
		--swagger_out=allow_merge,merge_file_name=kript:./docs/api \
		docs/proto/kript/api/*.proto

	# Workaround for https://github.com/grpc-ecosystem/grpc-gateway/issues/229.
	sed -i.bak "s/empty.Empty/types.Empty/g" pkg/proto/kript/api/*.pb.gw.go && \
rm pkg/proto/kript/api/*.pb.gw.go.bak

	go generate ./...

test:
	go test ./...

build:
	mkdir -p bin
	for main in $$(find cmd -type f -name "main.go"); do \
  		echo $$main; \
  		export out=$$(echo "$$main" | sed -E "s#cmd/(.*)/main.go#bin/\1#"); \
  		echo $$out; \
		go build -o "$$out" "$$main"; \
	done