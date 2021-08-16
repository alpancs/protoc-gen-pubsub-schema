package main

import (
	"os"
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

func Test_responseBuilder_build(t *testing.T) {
	tests := []struct {
		name    string
		request *pluginpb.CodeGeneratorRequest
		want    *pluginpb.CodeGeneratorResponse
	}{
		{
			name:    "sanity check",
			request: loadTestInput("test/user_add_comment.in"),
			want:    loadTestOutput("test/user_add_comment.out"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := newResponseBuilder(tt.request)
			got := b.build().GetFile()[0]
			want := tt.want.GetFile()[0]
			if got.GetName() != want.GetName() {
				t.Errorf("responseBuilder.build() = name=%v, want name=%v", got.GetName(), want.GetName())
			}
			if got.GetContent() != want.GetContent() {
				t.Errorf("responseBuilder.build() = content=%v, want content=%v", got.GetContent(), want.GetContent())
			}
		})
	}
}

func loadTestInput(path string) *pluginpb.CodeGeneratorRequest {
	return loadTest(path, new(pluginpb.CodeGeneratorRequest)).(*pluginpb.CodeGeneratorRequest)
}

func loadTestOutput(path string) *pluginpb.CodeGeneratorResponse {
	return loadTest(path, new(pluginpb.CodeGeneratorResponse)).(*pluginpb.CodeGeneratorResponse)
}

func loadTest(path string, message protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	raw, _ := os.ReadFile(path)
	prototext.Unmarshal(raw, message)
	return message
}
