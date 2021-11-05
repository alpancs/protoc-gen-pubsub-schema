package main

import (
	"errors"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type protoFilesType map[string]*descriptorpb.FileDescriptorProto
type messageTypesType map[string]*descriptorpb.DescriptorProto
type enumsType map[string]*descriptorpb.EnumDescriptorProto

type responseBuilder struct {
	request         *pluginpb.CodeGeneratorRequest
	schemaSyntax    string
	messageEncoding string
	protoFiles      protoFilesType
	messageTypes    messageTypesType
	enums           enumsType
}

func newResponseBuilder(request *pluginpb.CodeGeneratorRequest) (*responseBuilder, error) {
	if request == nil {
		return nil, errors.New("newResponseBuilder(request *pluginpb.CodeGeneratorRequest): request is nil")
	}
	return &responseBuilder{
		request,
		getSyntax(request),
		getEncoding(request),
		getProtoFiles(request),
		getMessageTypes(request),
		getEnums(request),
	}, nil
}

func getSyntax(request *pluginpb.CodeGeneratorRequest) string {
	if strings.Contains(request.GetParameter(), "schema-syntax=proto3") {
		return "proto3"
	}
	return "proto2"
}

func getEncoding(request *pluginpb.CodeGeneratorRequest) string {
	if strings.Contains(request.GetParameter(), "message-encoding=json") {
		return "json"
	}
	return "binary"
}

func getProtoFiles(request *pluginpb.CodeGeneratorRequest) protoFilesType {
	protoFiles := make(protoFilesType)
	for _, protoFile := range request.GetProtoFile() {
		protoFiles[protoFile.GetName()] = protoFile
	}
	return protoFiles
}

func getMessageTypes(request *pluginpb.CodeGeneratorRequest) messageTypesType {
	messageTypes := make(messageTypesType)
	for _, protoFile := range request.GetProtoFile() {
		messageTypes.setUsingFile(protoFile)
	}
	return messageTypes
}

func (ms messageTypesType) setUsingFile(file *descriptorpb.FileDescriptorProto) {
	packageName := strings.TrimSuffix("."+file.GetPackage(), ".")
	for _, m := range file.GetMessageType() {
		ms.setUsingMessage(packageName+".", m)
	}
}

func (ms messageTypesType) setUsingMessage(namePrefix string, message *descriptorpb.DescriptorProto) {
	fullMessageName := namePrefix + message.GetName()
	ms[fullMessageName] = message
	for _, m := range message.GetNestedType() {
		ms.setUsingMessage(fullMessageName+".", m)
	}
}

func getEnums(request *pluginpb.CodeGeneratorRequest) enumsType {
	enums := make(enumsType)
	for _, protoFile := range request.GetProtoFile() {
		enums.setUsingFile(protoFile)
	}
	return enums
}

func (es enumsType) setUsingFile(file *descriptorpb.FileDescriptorProto) {
	packageName := strings.TrimSuffix("."+file.GetPackage(), ".")
	for _, enum := range file.GetEnumType() {
		es[packageName+"."+enum.GetName()] = enum
	}
}

func (b *responseBuilder) build() *pluginpb.CodeGeneratorResponse {
	resp := new(pluginpb.CodeGeneratorResponse)
	for _, fileName := range b.request.GetFileToGenerate() {
		respFile, err := b.buildFile(fileName)
		if err != nil {
			errorMessage := err.Error()
			resp.Error = &errorMessage
			return resp
		}
		resp.File = append(resp.File, respFile)
	}
	return resp
}

func (b *responseBuilder) buildFile(reqFileName string) (*pluginpb.CodeGeneratorResponse_File, error) {
	respFileName := strings.TrimSuffix(reqFileName, ".proto") + ".pps"
	content, err := newContentBuilder(b).build(b.protoFiles[reqFileName])
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &respFileName,
		Content: &content,
	}, err
}
