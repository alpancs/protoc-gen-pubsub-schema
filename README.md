# protoc-gen-pubsub-schema

[![Go](https://github.com/alpancs/protoc-gen-pubsub-schema/actions/workflows/go.yml/badge.svg)](https://github.com/alpancs/protoc-gen-pubsub-schema/actions/workflows/go.yml)

This is a [protoc](https://github.com/protocolbuffers/protobuf) plugin
that assembles protocol buffer messages into a single-valid message for defining [Pub/Sub schemas](https://cloud.google.com/pubsub/docs/schemas).

## Installation

Run the following command to install `protoc-gen-pubsub-schema`.

```sh
go install github.com/alpancs/protoc-gen-pubsub-schema@latest
```

## Usage

You need to have `protoc` installed.
Follow <https://grpc.io/docs/protoc-installation> for instruction.

To use this plugin, just run `protoc` with an option `--pubsub-schema_out`.
`protoc` will automatically use `protoc-gen-pubsub-schema` executable file.
`protoc` and `protoc-gen-pubsub-schema` must be found in shell's `$PATH`.

```sh
# generate pubsub-proto-schema files using proto2 syntax that accept binary message encoding
protoc PROTO_FILES --pubsub-schema_out=OUT_DIR

# generate pubsub-proto-schema files using proto2 syntax that accept JSON message encoding
protoc PROTO_FILES --pubsub-schema_out=OUT_DIR --pubsub-schema_opt=message-encoding=json

# generate pubsub-proto-schema files using proto3 syntax that accept JSON message encoding
protoc PROTO_FILES --pubsub-schema_out=OUT_DIR --pubsub-schema_opt=message-encoding=json --pubsub-schema_opt=schema-syntax=proto3
```

## Example

The following example shows how to generate [example/article_commented.pps](example/article_commented.pps) from [example/article_commented.proto](example/article_commented.proto).

```sh
# include go compiled binaries in the $PATH if it hasn't been there yet
export PATH=$PATH:$(go env GOPATH)/bin

# generate example/article_commented.pps
protoc example/article_commented.proto --pubsub-schema_out=.
```
