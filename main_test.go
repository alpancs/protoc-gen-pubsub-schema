package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func Test_processIO(t *testing.T) {
	userAddCommentFile, err := os.Open("test/user_add_comment.protobuf")
	if err != nil {
		panic(err)
	}
	userAddCommentOutput, err := os.ReadFile("test/user_add_comment.pubsub.protobuf")
	if err != nil {
		panic(err)
	}

	type args struct {
		input io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantOutput []byte
	}{
		{
			name:       "sanity check",
			args:       args{userAddCommentFile},
			wantOutput: userAddCommentOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			processIO(tt.args.input, output)
			if gotOutput := output.Bytes(); !bytes.Equal(gotOutput, tt.wantOutput) {
				t.Errorf("processIO() = len=%v, want len=%v", len(gotOutput), len(tt.wantOutput))
			}
		})
	}
}
