package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func Test_process(t *testing.T) {
	type args struct {
		input io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantOutput []byte
		wantErr    bool
	}{
		{
			name:       "sanity check",
			args:       args{mustOpen("test/user_add_comment.in")},
			wantOutput: mustRead("test/user_add_comment.out"),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			if err := process(tt.args.input, output); (err != nil) != tt.wantErr {
				t.Errorf("process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput := output.Bytes(); len(gotOutput) != len(tt.wantOutput) {
				t.Errorf("process() = %v\nwant %v", mustDecodeResponse(gotOutput), mustDecodeResponse(tt.wantOutput))
			}
		})

		if closer, ok := tt.args.input.(io.Closer); ok {
			closer.Close()
		}
	}
}

func mustOpen(path string) io.Reader {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}

func mustRead(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}

func mustDecodeResponse(raw []byte) *pluginpb.CodeGeneratorResponse {
	resp := new(pluginpb.CodeGeneratorResponse)
	err := proto.Unmarshal(raw, resp)
	if err != nil {
		panic(err)
	}
	return resp
}
