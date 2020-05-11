install:
	go get \
		github.com/gogo/protobuf/protoc-gen-gogo \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/golang/mock/mockgen

vendor:
	echo "package example; import \"github.com/gogo/googleapis/google/api\"" > proto/temp.go
	go mod vendor
	rm proto/temp.go

generate:
	# Generate gogo, gRPC-Gateway, swagger, go-validators output.
	#
	# -I declares import folders, in order of importance
	# This is how proto resolves the protofile imports.
	# It will check for the protofile relative to each of these
	# folders and use the first one it finds.
	#
	# --gogo_out generates GoGo Protobuf output with gRPC plugin enabled.
	# --grpc-gateway_out generates gRPC-Gateway output.
	# --swagger_out generates an OpenAPI 2.0 specification for our gRPC-Gateway endpoints.
	# --govalidators_out generates Go validation files for our messages types, if specified.
	#
	# The lines starting with Mgoogle/... are proto import replacements,
	# which cause the generated file to import the specified packages
	# instead of the go_package's declared by the imported protof files.
	#
	# ./proto is the output directory.
	#
	# proto/example.proto is the location of the protofile we uerse.
	protoc \
		-I proto \
		-I vendor/github.com/grpc-ecosystem/grpc-gateway/ \
		-I vendor/github.com/gogo/googleapis/ \
		-I vendor/ \
		--gogo_out=plugins=grpc,paths=source_relative,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:\
./proto/ \
		--grpc-gateway_out=allow_patch_feature=false,paths=source_relative,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:\
./proto/ \
		--swagger_out=./ \
		proto/kript/api/*.proto

	# Workaround for https://github.com/grpc-ecosystem/grpc-gateway/issues/229.
	sed -i.bak "s/empty.Empty/types.Empty/g" proto/kript/api/*.pb.gw.go && rm proto/kript/api/*.pb.gw.go.bak

	go generate ./...
