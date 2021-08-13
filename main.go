package main

import (
	"io"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	err := process(os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
}

func process(input io.Reader, output io.Writer) error {
	req, err := decodeRequest(input)
	if err != nil {
		return encodeResponse(buildResponseError(err.Error()), output)
	}

	builder, err := newResponseBuilder(req)
	if err != nil {
		return encodeResponse(buildResponseError(err.Error()), output)
	}

	return encodeResponse(builder.build(), output)
}

func decodeRequest(input io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	rawInput, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	req := new(pluginpb.CodeGeneratorRequest)
	err = proto.Unmarshal(rawInput, req)
	return req, err
}

func encodeResponse(resp *pluginpb.CodeGeneratorResponse, output io.Writer) error {
	rawOutput, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = output.Write(rawOutput)
	return err
}
