package main

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/pluginpb"
)

func Test_getTopLevelMessage(t *testing.T) {
	tests := []struct {
		name      string
		parameter string
		want      string
	}{
		{
			name:      "empty",
			parameter: "",
			want:      "",
		},
		{
			// Return value on valid value
			name:      "single valid value",
			parameter: fmt.Sprintf("blah,blah,%s=Foo,blah,blah", keyTopLevelMessage),
			want:      "Foo",
		},
		{
			// Return empty string on invalid value
			name:      "single invalid value",
			parameter: "blah,blah,root-message=9,blah,blah",
			want:      "",
		},
		{
			// Return empty string for multiple values
			name:      "multiple values",
			parameter: fmt.Sprintf("blah,blah,%s=Foo,blah,%s=Bar,blah", keyTopLevelMessage, keyTopLevelMessage),
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &pluginpb.CodeGeneratorRequest{
				Parameter: &tt.parameter,
			}
			if got := getTopLevelMessage(request); got != tt.want {
				t.Errorf("getTopLevelMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
