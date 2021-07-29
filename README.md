# protoc-gen-pubsub-schema

This is a [protoc](https://github.com/protocolbuffers/protobuf) plugin
that assembles protocol buffer messages into a single-valid message for defining [Pub/Sub schemas](https://cloud.google.com/pubsub/docs/schemas).

## Installation

Run the following command to install `protoc-gen-pubsub-schema`.

```sh
go install github.com/alpancs/protoc-gen-pubsub-schema
```

## Usage

You need to have `protoc` installed.
Follow <https://grpc.io/docs/protoc-installation> for instruction.

To use this plugin, just run `protoc` with an option `--pubsub-schema_out`.
`protoc` will automatically use `protoc-gen-pubsub-schema` executable file.
`protoc` and `protoc-gen-pubsub-schema` must be found in shell's `$PATH`.

```sh
protoc PROTO_FILES --pubsub-schema_out=OUT_DIR
```

## Example

The following example shows how to generate [example/user_add_comment.pubsub.proto](example/user_add_comment.pubsub.proto) from [example/user_add_comment.proto](example/user_add_comment.proto).

```sh
# include go compiled binaries in the $PATH if it hasn't been there yet
export PATH=$PATH:$(go env GOPATH)/bin

# generate example/user_add_comment.pubsub.proto
protoc example/user_add_comment.proto --pubsub-schema_out=.
```
