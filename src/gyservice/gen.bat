protoc -I ./ ./service.proto --go_out=plugins=grpc:./go --java_out=./java