package main

import (
	"fmt"
	"io"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type responseBuilder struct {
	request      *pluginpb.CodeGeneratorRequest
	protoFiles   map[string]*descriptorpb.FileDescriptorProto
	messageTypes map[string]*descriptorpb.DescriptorProto
}

func buildResponseError(errorMessage string) *pluginpb.CodeGeneratorResponse {
	return &pluginpb.CodeGeneratorResponse{Error: &errorMessage}
}

func newResponseBuilder(req *pluginpb.CodeGeneratorRequest) responseBuilder {
	builder := responseBuilder{
		req,
		make(map[string]*descriptorpb.FileDescriptorProto),
		make(map[string]*descriptorpb.DescriptorProto),
	}
	builder.initIndex()
	return builder
}

func (b *responseBuilder) initIndex() {
	for _, protoFile := range b.request.GetProtoFile() {
		b.protoFiles[protoFile.GetName()] = protoFile
		packageName := strings.TrimSuffix("."+protoFile.GetPackage(), ".")
		b.initIndexByMessages(packageName, protoFile.GetMessageType())
	}
}

func (b *responseBuilder) initIndexByMessages(messageNamePrefix string, messages []*descriptorpb.DescriptorProto) {
	for _, message := range messages {
		messageName := messageNamePrefix + "." + message.GetName()
		b.messageTypes[messageName] = message
		b.initIndexByMessages(messageName, message.GetNestedType())
	}
}

func (b responseBuilder) build() *pluginpb.CodeGeneratorResponse {
	resp := new(pluginpb.CodeGeneratorResponse)
	for _, fileName := range b.request.GetFileToGenerate() {
		respFile, err := b.buildFile(fileName)
		if err != nil {
			return buildResponseError(err.Error())
		}
		resp.File = append(resp.File, respFile)
	}
	return resp
}

func (b responseBuilder) buildFile(reqFileName string) (*pluginpb.CodeGeneratorResponse_File, error) {
	respFileName := strings.TrimSuffix(reqFileName, ".proto") + ".pubsub.proto"
	content, err := b.buildContent(b.protoFiles[reqFileName])
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &respFileName,
		Content: &content,
	}, err
}

func (b responseBuilder) buildContent(protoFile *descriptorpb.FileDescriptorProto) (string, error) {
	if len(protoFile.GetMessageType()) != 1 {
		return "", fmt.Errorf(
			"only one top-level type may be defined in the file \"%s\". use nested types instead (https://developers.google.com/protocol-buffers/docs/proto3#nested)",
			protoFile.GetName(),
		)
	}

	contentBuilder := new(strings.Builder)
	fmt.Fprint(contentBuilder, "syntax = \"proto3\";\n\n")
	b.buildMessage(contentBuilder, protoFile.GetMessageType()[0], 0)
	return contentBuilder.String(), nil
}

func (b responseBuilder) buildMessage(output io.Writer, message *descriptorpb.DescriptorProto, level int) {
	fmt.Fprintf(output, "%smessage %s {\n", buildIndent(level), message.GetName())
	for _, field := range message.GetField() {
		b.buildField(output, field, level+1)
	}
	fmt.Fprintf(output, "%s}\n", buildIndent(level))
}

func (b responseBuilder) buildField(output io.Writer, field *descriptorpb.FieldDescriptorProto, level int) {
	fieldType := strings.ToLower(strings.TrimPrefix(field.GetType().String(), "TYPE_"))

	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		fmt.Fprintln(output)
		defer fmt.Fprintln(output)
		b.buildMessage(output, b.messageTypes[field.GetTypeName()], level)
		fieldType = getShortTypeName(field.GetTypeName())
	}

	fmt.Fprint(output, buildIndent(level))
	if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		fmt.Fprint(output, "repeated ")
	}
	fmt.Fprintf(output, "%s %s = %d;\n", fieldType, field.GetName(), field.GetNumber())
}

func buildIndent(level int) string {
	return strings.Repeat("  ", level)
}

func getShortTypeName(typeName string) string {
	return typeName[strings.LastIndexByte(typeName, '.')+1:]
}
