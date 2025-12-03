package xutil

// TraceIDKey
// W3C Trace Context 规范
//
//	Key: traceparent
//	Format: {version}-{trace-id}-{parent-id}-{trace-flags}
//	Value: 00-4bf92f08537468a83416fd253b2d1840-00f067aa0ba902b7-01
//
// B3 规范 (Single-Header)
//
//	Key: b3
//	Format: {TraceID}-{SpanID}-{Sampled}-{ParentSpanID}
//	Value: 80f198ee56343ba864fe2592ccaa2846-64fe2592ccaa2846-1-343ba864fe2592cc
var TraceIDKey = "traceparent"

type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}
