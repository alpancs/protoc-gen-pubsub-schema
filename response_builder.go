package main

import (
	"fmt"
	"io"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type responseBuilder struct {
	request *pluginpb.CodeGeneratorRequest
}

func buildResponseError(err error) *pluginpb.CodeGeneratorResponse {
	errorMessage := err.Error()
	return &pluginpb.CodeGeneratorResponse{Error: &errorMessage}
}

func (b responseBuilder) build() *pluginpb.CodeGeneratorResponse {
	resp := new(pluginpb.CodeGeneratorResponse)
	for _, fileName := range b.request.GetFileToGenerate() {
		respFile, err := b.buildFile(fileName)
		if err != nil {
			errorMessage := err.Error()
			resp.Error = &errorMessage
			break
		}
		resp.File = append(resp.File, respFile)
	}
	return resp
}

func (b responseBuilder) buildFile(reqFileName string) (*pluginpb.CodeGeneratorResponse_File, error) {
	respFileName := strings.TrimSuffix(reqFileName, ".proto") + ".pubsub.proto"
	content, err := b.buildContent(b.findProtoFileByName(reqFileName))
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &respFileName,
		Content: &content,
	}, err
}

func (b responseBuilder) findProtoFileByName(desiredName string) *descriptorpb.FileDescriptorProto {
	for _, protoFile := range b.request.GetProtoFile() {
		if protoFile.GetName() == desiredName {
			return protoFile
		}
	}
	return nil
}

func (b responseBuilder) findMessageByName(desiredName string) *descriptorpb.DescriptorProto {
	for _, protoFile := range b.request.GetProtoFile() {
		packageName := strings.TrimSuffix("."+protoFile.GetPackage(), ".")
		message := b.findNestedMessageByName(desiredName, protoFile.GetMessageType(), packageName)
		if message != nil {
			return message
		}
	}
	return nil
}

func (b responseBuilder) findNestedMessageByName(desiredName string, messages []*descriptorpb.DescriptorProto, messageNamePrefix string) *descriptorpb.DescriptorProto {
	for _, message := range messages {
		messageName := messageNamePrefix + "." + message.GetName()
		if messageName == desiredName {
			return message
		}
		nestedMessage := b.findNestedMessageByName(desiredName, message.GetNestedType(), messageName)
		if nestedMessage != nil {
			return nestedMessage
		}
	}
	return nil
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
		b.buildMessage(output, b.findMessageByName(field.GetTypeName()), level)
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
