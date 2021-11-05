package main

import (
	"testing"
)

func Test_getParentName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty name",
			args: args{""},
			want: "",
		},
		{
			name: "normal name",
			args: args{".example.UserAddComment.User"},
			want: ".example.UserAddComment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getParentName(tt.args.name); got != tt.want {
				t.Errorf("getParentName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getChildName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty name",
			args: args{""},
			want: "",
		},
		{
			name: "normal name",
			args: args{".example.UserAddComment.User"},
			want: "User",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getChildName(tt.args.name); got != tt.want {
				t.Errorf("getChildName() = %v, want %v", got, tt.want)
			}
		})
	}
}
