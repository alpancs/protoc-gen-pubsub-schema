package main

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type contentBuilder struct {
	rb     responseBuilder
	output *strings.Builder
}

func newContentBuilder(req responseBuilder) contentBuilder {
	builder := contentBuilder{req, new(strings.Builder)}
	return builder
}

func (b contentBuilder) build(protoFile *descriptorpb.FileDescriptorProto) (string, error) {
	if protoFile == nil {
		return "", errors.New("protoFile is nil")
	}

	if len(protoFile.GetMessageType()) != 1 {
		return "", errors.New(protoFile.GetName() + ": only one top-level type may be defined in a file. use nested types or use imports. see https://developers.google.com/protocol-buffers/docs/proto3 for details.")
	}

	b.output.WriteString(`syntax = "` + b.getOutputSyntax() + `";` + "\n\n")
	b.buildMessage(protoFile.GetMessageType()[0], 0)
	return b.output.String(), nil
}

func (b contentBuilder) getOutputSyntax() string {
	if strings.Contains(b.rb.request.GetParameter(), "syntax=proto3") {
		return "proto3"
	}
	return "proto2"
}

func (b contentBuilder) buildMessage(message *descriptorpb.DescriptorProto, level int) {
	b.output.WriteString(buildIndent(level) + "message " + message.GetName() + " {\n")
	for _, field := range message.GetField() {
		b.buildField(field, level+1)
	}
	b.output.WriteString(buildIndent(level) + "}\n")
}

func (b contentBuilder) buildField(field *descriptorpb.FieldDescriptorProto, level int) {
	fieldType := strings.ToLower(strings.TrimPrefix(field.GetType().String(), "TYPE_"))
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		fieldType = b.buildFieldType(field.GetTypeName(), level)
	}
	b.output.WriteString(buildIndent(level))
	b.buildFieldLabel(field.GetLabel())
	b.output.WriteString(fmt.Sprintf("%s %s = %d;\n", fieldType, field.GetName(), field.GetNumber()))
}

func (b contentBuilder) buildFieldType(typeName string, level int) string {
	if typeName, ok := wktMapping[typeName]; ok && b.hasJSONEncoding() {
		return typeName
	}

	b.output.WriteString("\n")
	b.buildMessage(b.rb.messageTypes[typeName], level)
	b.output.WriteString("\n")
	return typeName[strings.LastIndexByte(typeName, '.')+1:]
}

func (b contentBuilder) hasJSONEncoding() bool {
	return strings.Contains(b.rb.request.GetParameter(), "encoding=json")
}

func (b contentBuilder) buildFieldLabel(label descriptorpb.FieldDescriptorProto_Label) {
	if label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		b.output.WriteString("repeated ")
	} else if b.getOutputSyntax() == "proto2" {
		b.output.WriteString(strings.ToLower(strings.TrimPrefix(label.String(), "LABEL_")) + " ")
	}
}

func buildIndent(level int) string {
	return strings.Repeat("  ", level)
}
