package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type debtsType = map[string]bool

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

	fmt.Fprintf(b.output, `syntax = "%s";`, b.schemaSyntax)
	b.output.WriteString("\n\n")
	b.buildMessage(protoFile.GetMessageType()[0], 0)
	return b.output.String(), nil
}

func (b *contentBuilder) buildMessage(message *descriptorpb.DescriptorProto, level int) {
	fmt.Fprintf(b.output, "%smessage %s {\n", buildIndent(level), message.GetName())
	defer fmt.Fprintf(b.output, "%s}\n", buildIndent(level))
	b.payDebts(b.buildFields(message, level), level+1)
}

func (b *contentBuilder) buildFields(message *descriptorpb.DescriptorProto, level int) debtsType {
	debts := make(debtsType)
	for _, field := range message.GetField() {
		if debt := b.buildField(field, level+1); debt != "" {
			debts[debt] = true
		}
	}
	return debts
}

func (b *contentBuilder) buildField(field *descriptorpb.FieldDescriptorProto, level int) string {
	fieldType, debt := b.getFieldType(field)
	fmt.Fprintf(b.output, "%s%s%s %s = %d;\n",
		buildIndent(level),
		b.getLabelPrefix(field.GetLabel()),
		fieldType,
		field.GetName(),
		field.GetNumber(),
	)
	return debt
}

func (b *contentBuilder) getFieldType(field *descriptorpb.FieldDescriptorProto) (string, string) {
	if field.GetType() != descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		return strings.ToLower(strings.TrimPrefix(field.GetType().String(), "TYPE_")), ""
	}
	fullMessageName := field.GetTypeName()
	wkt := wktMapping[fullMessageName]
	if b.messageEncoding == "json" && wkt != "" {
		return wkt, ""
	}
	return getLocalName(fullMessageName), fullMessageName
}

func (b *contentBuilder) payDebts(debts debtsType, level int) {
	for debt := range debts {
		b.payDebt(debt, level)
	}
}

func (b *contentBuilder) payDebt(fullMessageName string, level int) {
	message := b.messageTypes[fullMessageName]
	localName := getLocalName(fullMessageName)
	message.Name = &localName
	b.output.WriteString("\n")
	b.buildMessage(message, level)
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

func buildIndent(level int) string {
	return strings.Repeat("  ", level)
}

var localNamePattern = regexp.MustCompile(`\.[A-Za-z0-9]`)

func getLocalName(fullMessageName string) string {
	return localNamePattern.ReplaceAllStringFunc(
		fullMessageName,
		func(s string) string { return strings.ToUpper(s[1:]) },
	)
}
