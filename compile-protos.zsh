#!/bin/zsh
# compile protocol buffers

protoc -I=/home/sam/projects/protos/syncapod-protos/ \
	--go_out=internal/protos/ \
	--go-grpc_out=internal/protos/ \
	/home/sam/projects/protos/syncapod-protos/*

# add _id to all protobuf objects
add-bson-id internal/protos
