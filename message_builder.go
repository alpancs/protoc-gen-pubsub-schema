package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type messageBuilder struct {
	*contentBuilder
	message          *descriptorpb.DescriptorProto
	level            int
	externalMessages []*descriptorpb.DescriptorProto
	externalEnums    []*descriptorpb.EnumDescriptorProto
}

func newMessageBuilder(b *contentBuilder, message *descriptorpb.DescriptorProto, level int) *messageBuilder {
	return &messageBuilder{b, message, level, nil, nil}
}

func (b *messageBuilder) build() {
	b.message.NestedType = nil
	b.message.EnumType = nil
	fmt.Fprintf(b.output, "%smessage %s {\n", buildIndent(b.level), b.message.GetName())
	b.buildFields()
	b.buildMessages(b.message.GetNestedType(), b.level+1)
	b.buildEnums(b.message.GetEnumType(), b.level+1)
	fmt.Fprintf(b.output, "%s}\n", buildIndent(b.level))
}

func (b *messageBuilder) buildFields() {
	for _, field := range b.message.GetField() {
		fmt.Fprint(b.output, buildIndent(b.level+1))
		label := field.GetLabel()
		if b.schemaSyntax == "proto2" || label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			fmt.Fprintf(b.output, "%s ", strings.ToLower(strings.TrimPrefix(label.String(), "LABEL_")))
		}
		fmt.Fprintf(b.output, "%s %s = %d;\n", b.buildFieldType(field), field.GetName(), field.GetNumber())
	}
}

func (b *messageBuilder) buildFieldType(field *descriptorpb.FieldDescriptorProto) string {
	typeName := field.GetTypeName()
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if b.messageEncoding == "json" && wktMapping[typeName] != "" {
			return wktMapping[typeName]
		}
		internalName := pascalCase(typeName)
		internalMessage := b.messageTypes[field.GetTypeName()]
		internalMessage.Name = &internalName
		b.message.NestedType = append(b.message.NestedType, internalMessage)
		return internalName
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		internalName := pascalCase(typeName)
		internalEnum := b.enums[field.GetTypeName()]
		internalEnum.Name = &internalName
		b.message.EnumType = append(b.message.EnumType, internalEnum)
		return internalName
	default:
		return strings.ToLower(strings.TrimPrefix(field.GetType().String(), "TYPE_"))
	}
}

func getParentName(name string) string {
	lastDotIndex := strings.LastIndexByte(name, '.')
	if lastDotIndex == -1 {
		return name
	}
	return name[:lastDotIndex]
}

func getChildName(name string) string {
	return name[strings.LastIndexByte(name, '.')+1:]
}

func pascalCase(name string) string {
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
