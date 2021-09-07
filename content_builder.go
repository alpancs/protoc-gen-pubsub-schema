package main

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type contentBuilder struct {
	*responseBuilder
	output *strings.Builder
}

func newContentBuilder(b *responseBuilder) *contentBuilder {
	return &contentBuilder{b, new(strings.Builder)}
}

func (b *contentBuilder) build(protoFile *descriptorpb.FileDescriptorProto) (string, error) {
	if protoFile == nil {
		return "", errors.New("build(protoFile *descriptorpb.FileDescriptorProto): protoFile is nil")
	}

	if len(protoFile.GetMessageType()) != 1 {
		return "", errors.New(protoFile.GetName() + ": only one top-level type may be defined in a file (see https://cloud.google.com/pubsub/docs/schemas#schema_types). use nested types or imports (see https://developers.google.com/protocol-buffers/docs/proto)")
	}

	compVersion := b.request.GetCompilerVersion()
	fmt.Fprintln(b.output, "// Code generated by protoc-gen-pubsub-schema. DO NOT EDIT.")
	fmt.Fprintln(b.output, "// versions:")
	fmt.Fprintln(b.output, "// 	protoc-gen-pubsub-schema v1.5.0")
	fmt.Fprintf(b.output, "// 	protoc                   v%d.%d.%d%s\n", compVersion.GetMajor(), compVersion.GetMinor(), compVersion.GetPatch(), compVersion.GetSuffix())
	fmt.Fprintf(b.output, "// source: %s\n\n", protoFile.GetName())
	fmt.Fprintf(b.output, `syntax = "%s";`, b.schemaSyntax)
	fmt.Fprint(b.output, "\n\n")
	b.buildMessage("", protoFile.GetMessageType()[0], 0)
	b.buildEnums(protoFile.GetEnumType(), 0)
	return b.output.String(), nil
}

func (b *contentBuilder) buildMessage(prefix string, message *descriptorpb.DescriptorProto, level int) {
	fmt.Fprintf(b.output, "%smessage %s%s {\n", buildIndent(level), prefix, message.GetName())
	for _, field := range message.GetField() {
		fmt.Fprintf(b.output, "%s%s%s %s = %d;\n",
			buildIndent(level+1),
			b.getLabelPrefix(field.GetLabel()),
			b.getFieldType(field),
			field.GetName(),
			field.GetNumber(),
		)
	}
	b.buildNestedTypes(message.GetNestedType(), level+1)
	b.buildEnums(message.GetEnumType(), level+1)
	b.buildOtherTypes(message.GetField(), level+1)
	fmt.Fprintf(b.output, "%s}\n", buildIndent(level))
}

func (b *contentBuilder) getLabelPrefix(label descriptorpb.FieldDescriptorProto_Label) string {
	if label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		return "repeated "
	}
	if b.schemaSyntax == "proto2" {
		return strings.ToLower(strings.TrimPrefix(label.String(), "LABEL_")) + " "
	}
	return ""
}

func (b *contentBuilder) getFieldType(field *descriptorpb.FieldDescriptorProto) string {
	typeName := field.GetTypeName()
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if b.messageEncoding == "json" && wktMapping[typeName] != "" {
			return wktMapping[typeName]
		}
		return b.getLocalName(typeName)
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return typeName[strings.LastIndexByte(typeName, '.')+1:]
	default:
		return strings.ToLower(strings.TrimPrefix(field.GetType().String(), "TYPE_"))
	}
}

func (b *contentBuilder) getLocalName(name string) string {
	if b.isNestedType(name) {
		return name[strings.LastIndexByte(name, '.')+1:]
	}
	sb := new(strings.Builder)
	for i, c := range name {
		if i > 0 && name[i-1] == '.' {
			sb.WriteString(strings.ToUpper(string(c)))
		} else if c != '.' {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}

func (b *contentBuilder) buildNestedTypes(messages []*descriptorpb.DescriptorProto, level int) {
	for _, message := range messages {
		fmt.Fprintln(b.output)
		b.buildMessage("", message, level)
	}
}

func (b *contentBuilder) buildEnums(enums []*descriptorpb.EnumDescriptorProto, level int) {
	for _, enum := range enums {
		fmt.Fprintln(b.output)
		b.buildEnum(enum, level)
	}
}

func (b *contentBuilder) buildEnum(enum *descriptorpb.EnumDescriptorProto, level int) {
	fmt.Fprintf(b.output, "%senum %s {\n", buildIndent(level), enum.GetName())
	for _, value := range enum.GetValue() {
		fmt.Fprintf(b.output, "%s%s = %d;\n", buildIndent(level+1), value.GetName(), value.GetNumber())
	}
	fmt.Fprintf(b.output, "%s}\n", buildIndent(level))
}

func (b *contentBuilder) buildOtherTypes(fields []*descriptorpb.FieldDescriptorProto, level int) {
	built := make(map[string]bool)
	for _, field := range fields {
		typeName := field.GetTypeName()
		if field.GetType() != descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			continue
		}
		if b.messageEncoding == "json" && wktMapping[typeName] != "" {
			continue
		}
		if b.isNestedType(typeName) {
			continue
		}
		if built[typeName] {
			continue
		}
		fmt.Fprintln(b.output)
		b.buildMessage("Generated", b.messageTypes[typeName], level)
		built[typeName] = true
	}
}

func (b *contentBuilder) isNestedType(name string) bool {
	return b.messageTypes[name[:strings.LastIndexByte(name, '.')]] != nil
}

func buildIndent(level int) string {
	return strings.Repeat("  ", level)
}
