package main

import (
	"io"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	processIO(os.Stdin, os.Stdout)
}

func processIO(io.Reader, io.Writer) {
	req, err := decodeRequest(os.Stdin)
	if err != nil {
		err = encodeResponse(buildResponseError(err), os.Stdout)
		if err != nil {
			panic(err)
		}
		return
	}

	err = encodeResponse(responseBuilder{req}.build(), os.Stdout)
	if err != nil {
		panic(err)
	}
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
