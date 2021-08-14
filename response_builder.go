package main

import (
	"errors"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type protoFilesType map[string]*descriptorpb.FileDescriptorProto
type messageTypesType map[string]*descriptorpb.DescriptorProto

type responseBuilder struct {
	request         *pluginpb.CodeGeneratorRequest
	schemaSyntax    string
	messageEncoding string
	protoFiles      protoFilesType
	messageTypes    messageTypesType
}

func buildResponseError(errorMessage string) *pluginpb.CodeGeneratorResponse {
	return &pluginpb.CodeGeneratorResponse{Error: &errorMessage}
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
	}, nil
}

func getSyntax(request *pluginpb.CodeGeneratorRequest) string {
	if strings.Contains(request.GetParameter(), "syntax=proto3") {
		return "proto3"
	}
	return "proto2"
}

func getEncoding(request *pluginpb.CodeGeneratorRequest) string {
	if strings.Contains(request.GetParameter(), "message-messageEncoding=json") {
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

func (b *responseBuilder) build() *pluginpb.CodeGeneratorResponse {
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

func (b *responseBuilder) buildFile(reqFileName string) (*pluginpb.CodeGeneratorResponse_File, error) {
	respFileName := strings.TrimSuffix(reqFileName, ".proto") + ".pps"
	content, err := newContentBuilder(b).build(b.protoFiles[reqFileName])
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &respFileName,
		Content: &content,
	}, err
}
