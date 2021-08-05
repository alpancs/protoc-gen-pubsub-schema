package main

var wktMapping = map[string]string{
	".google.protobuf.Int32Value":  "int32",
	".google.protobuf.Int64Value":  "int64",
	".google.protobuf.UInt32Value": "uint32",
	".google.protobuf.UInt64Value": "uint64",
	".google.protobuf.DoubleValue": "double",
	".google.protobuf.FloatValue":  "float",
	".google.protobuf.BoolValue":   "bool",
	".google.protobuf.StringValue": "string",
	".google.protobuf.BytesValue":  "bytes",
	".google.protobuf.Duration":    "string",
	".google.protobuf.Timestamp":   "string",
}
