package main

import (
	"io"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	req, err := decodeRequest(os.Stdin)
	if err != nil {
		exitWithResponseError(err.Error(), os.Stdout)
	}

	builder, err := newResponseBuilder(req)
	if err != nil {
		exitWithResponseError(err.Error(), os.Stdout)
	}

	encodeResponse(builder.build(), os.Stdout)
}

func exitWithResponseError(errorMessage string, output io.Writer) {
	encodeResponse(&pluginpb.CodeGeneratorResponse{Error: &errorMessage}, output)
	os.Exit(0)
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
