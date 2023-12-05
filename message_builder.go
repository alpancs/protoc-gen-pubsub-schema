package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type messageBuilder struct {
	*contentBuilder
	level            int
	message          *descriptorpb.DescriptorProto
	externalMessages []*descriptorpb.DescriptorProto
	externalEnums    []*descriptorpb.EnumDescriptorProto
}

func newMessageBuilder(b *contentBuilder, level int, message *descriptorpb.DescriptorProto) *messageBuilder {
	return &messageBuilder{b, level, message, nil, nil}
}

func (b *messageBuilder) build() {
	fmt.Fprintf(b.output, "%smessage %s {\n", buildIndent(b.level), b.message.GetName())
	b.buildFields()
	b.buildMessages(b.level+1, append(b.message.GetNestedType(), b.externalMessages...))
	b.buildEnums(b.level+1, append(b.message.GetEnumType(), b.externalEnums...))
	fmt.Fprintf(b.output, "%s}\n", buildIndent(b.level))
}

func (b *messageBuilder) buildFields() {
	for _, field := range b.message.GetField() {
		fmt.Fprint(b.output, buildIndent(b.level+1))
		label := field.GetLabel()
		if label == descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL || label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			fmt.Fprintf(b.output, "%s ", strings.ToLower(strings.TrimPrefix(label.String(), "LABEL_")))
		}
		fmt.Fprintf(b.output, "%s %s = %d;\n", b.buildFieldType(field), field.GetName(), field.GetNumber())
	}
}

func (b *messageBuilder) buildFieldType(field *descriptorpb.FieldDescriptorProto) string {
	switch {
	case b.isInternalDefinition(field):
		return getChildName(field.GetTypeName())
	case field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if b.messageEncoding == "json" && wktMapping[field.GetTypeName()] != "" {
			return wktMapping[field.GetTypeName()]
		}
		internalName := pascalCase(field.GetTypeName())
		message := b.messageTypes[field.GetTypeName()]
		message.Name = &internalName
		b.externalMessages = append(b.externalMessages, message)
		return internalName
	case field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		internalName := pascalCase(field.GetTypeName())
		enum := b.enums[field.GetTypeName()]
		enum.Name = &internalName
		b.externalEnums = append(b.externalEnums, enum)
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
		if c == '.' || c == '_' || c == '-' {
			continue
		}
		if i > 0 && (name[i-1] == '.' || name[i-1] == '_' || name[i-1] == '-') {
			sb.WriteString(strings.ToUpper(string(c)))
		} else {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}
