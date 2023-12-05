package main

import (
	"errors"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type responseBuilder struct {
	request         *pluginpb.CodeGeneratorRequest
	schemaSyntax    string
	messageEncoding string
	protoFiles      map[string]*descriptorpb.FileDescriptorProto
	messageTypes    map[string]*descriptorpb.DescriptorProto
	enums           map[string]*descriptorpb.EnumDescriptorProto
	fileTypeNames   map[*descriptorpb.FileDescriptorProto][]string
}

func newResponseBuilder(request *pluginpb.CodeGeneratorRequest) (*responseBuilder, error) {
	if request == nil {
		return nil, errors.New("newResponseBuilder(request *pluginpb.CodeGeneratorRequest): request is nil")
	}
	builder := &responseBuilder{
		request,
		getSyntax(request),
		getEncoding(request),
		make(map[string]*descriptorpb.FileDescriptorProto),
		make(map[string]*descriptorpb.DescriptorProto),
		make(map[string]*descriptorpb.EnumDescriptorProto),
		make(map[*descriptorpb.FileDescriptorProto][]string),
	}
	builder.initProtoFiles()
	builder.initTypes()
	return builder, nil
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

func (b *responseBuilder) initProtoFiles() {
	for _, protoFile := range b.request.GetProtoFile() {
		b.protoFiles[protoFile.GetName()] = protoFile
	}
}

func (b *responseBuilder) initTypes() {
	for _, protoFile := range b.request.GetProtoFile() {
		b.initTypesInFile(protoFile)
	}
}

func (b *responseBuilder) initTypesInFile(file *descriptorpb.FileDescriptorProto) {
	packageName := strings.TrimSuffix("."+file.GetPackage(), ".")
	for _, m := range file.GetMessageType() {
		fullName := packageName + "." + m.GetName()
		b.messageTypes[fullName] = m
		b.fileTypeNames[file] = append(b.fileTypeNames[file], fullName)
		b.initTypesInMessage(file, fullName, m)
	}
	for _, e := range file.GetEnumType() {
		fullName := packageName + "." + e.GetName()
		b.enums[fullName] = e
		b.fileTypeNames[file] = append(b.fileTypeNames[file], fullName)
	}
}

func (b *responseBuilder) initTypesInMessage(file *descriptorpb.FileDescriptorProto, parentName string, message *descriptorpb.DescriptorProto) {
	for _, m := range message.GetNestedType() {
		fullName := parentName + "." + m.GetName()
		b.messageTypes[fullName] = m
		b.fileTypeNames[file] = append(b.fileTypeNames[file], fullName)
		b.initTypesInMessage(file, fullName, m)
	}
	for _, e := range message.GetEnumType() {
		fullName := parentName + "." + e.GetName()
		b.enums[fullName] = e
		b.fileTypeNames[file] = append(b.fileTypeNames[file], fullName)
	}
}

func (b *responseBuilder) build() *pluginpb.CodeGeneratorResponse {
	// Enable experimental support of `optional` fields in `proto3`
	supportedFeatures := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	resp := &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: &supportedFeatures,
	}

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
	content, err := newContentBuilder(b, b.protoFiles[reqFileName]).build()
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &respFileName,
		Content: &content,
	}, err
}
