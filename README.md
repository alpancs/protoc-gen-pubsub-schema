# protoc-gen-pubsub-schema

This is a [protoc](https://github.com/protocolbuffers/protobuf) plugin
that assembles protocol buffer messages into a single-valid message for defining [Pub/Sub schemas](https://cloud.google.com/pubsub/docs/schemas).

## Installation

To install `protoc-gen-pubsub-schema`, run the following command.

```sh
go install github.com/alpancs/protoc-gen-pubsub-schema
```

## Example

To run the example below, you need to have `protoc` installed.
Follow <https://grpc.io/docs/protoc-installation> for the installation.

Run the following command to generate `example/user_add_comment.pubsub.proto` from `example/user_add_comment.proto`.

```sh
# include go compiled binaries in the $PATH if it hasn't been there yet
export PATH=$PATH:$(go env GOPATH)/bin

# generate example/user_add_comment.pubsub.proto
protoc example/user_add_comment.proto --pubsub-schema_out .
```
