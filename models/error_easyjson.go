// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	strfmt "github.com/go-openapi/strfmt"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonE34310f8DecodeGithubComElasticApmServerModels(in *jlexer.Lexer, out *Error) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "context":
			if in.IsNull() {
				in.Skip()
				out.Context = nil
			} else {
				if out.Context == nil {
					out.Context = new(Context)
				}
				(*out.Context).UnmarshalEasyJSON(in)
			}
		case "culprit":
			out.Culprit = string(in.String())
		case "exception":
			if in.IsNull() {
				in.Skip()
				out.Exception = nil
			} else {
				if out.Exception == nil {
					out.Exception = new(Exception)
				}
				easyjsonE34310f8DecodeGithubComElasticApmServerModels1(in, &*out.Exception)
			}
		case "id":
			(out.ID).UnmarshalEasyJSON(in)
		case "log":
			if in.IsNull() {
				in.Skip()
				out.Log = nil
			} else {
				if out.Log == nil {
					out.Log = new(LogRecord)
				}
				easyjsonE34310f8DecodeGithubComElasticApmServerModels2(in, &*out.Log)
			}
		case "timestamp":
			if in.IsNull() {
				in.Skip()
				out.Timestamp = nil
			} else {
				if out.Timestamp == nil {
					out.Timestamp = new(strfmt.DateTime)
				}
				(*out.Timestamp).UnmarshalEasyJSON(in)
			}
		case "transaction":
			if in.IsNull() {
				in.Skip()
				out.Transaction = nil
			} else {
				if out.Transaction == nil {
					out.Transaction = new(ErrorTransaction)
				}
				easyjsonE34310f8DecodeGithubComElasticApmServerModels3(in, &*out.Transaction)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE34310f8EncodeGithubComElasticApmServerModels(out *jwriter.Writer, in Error) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Context != nil {
		const prefix string = ",\"context\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Context).MarshalEasyJSON(out)
	}
	if in.Culprit != "" {
		const prefix string = ",\"culprit\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Culprit))
	}
	if in.Exception != nil {
		const prefix string = ",\"exception\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjsonE34310f8EncodeGithubComElasticApmServerModels1(out, *in.Exception)
	}
	if in.ID != "" {
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.ID).MarshalEasyJSON(out)
	}
	if in.Log != nil {
		const prefix string = ",\"log\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjsonE34310f8EncodeGithubComElasticApmServerModels2(out, *in.Log)
	}
	{
		const prefix string = ",\"timestamp\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Timestamp == nil {
			out.RawString("null")
		} else {
			(*in.Timestamp).MarshalEasyJSON(out)
		}
	}
	if in.Transaction != nil {
		const prefix string = ",\"transaction\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjsonE34310f8EncodeGithubComElasticApmServerModels3(out, *in.Transaction)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Error) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE34310f8EncodeGithubComElasticApmServerModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Error) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE34310f8EncodeGithubComElasticApmServerModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Error) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE34310f8DecodeGithubComElasticApmServerModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Error) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE34310f8DecodeGithubComElasticApmServerModels(l, v)
}
func easyjsonE34310f8DecodeGithubComElasticApmServerModels3(in *jlexer.Lexer, out *ErrorTransaction) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			(out.ID).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE34310f8EncodeGithubComElasticApmServerModels3(out *jwriter.Writer, in ErrorTransaction) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.ID).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}
func easyjsonE34310f8DecodeGithubComElasticApmServerModels2(in *jlexer.Lexer, out *LogRecord) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "level":
			if in.IsNull() {
				in.Skip()
				out.Level = nil
			} else {
				if out.Level == nil {
					out.Level = new(string)
				}
				*out.Level = string(in.String())
			}
		case "logger_name":
			if in.IsNull() {
				in.Skip()
				out.LoggerName = nil
			} else {
				if out.LoggerName == nil {
					out.LoggerName = new(string)
				}
				*out.LoggerName = string(in.String())
			}
		case "message":
			if in.IsNull() {
				in.Skip()
				out.Message = nil
			} else {
				if out.Message == nil {
					out.Message = new(string)
				}
				*out.Message = string(in.String())
			}
		case "param_message":
			out.ParamMessage = string(in.String())
		case "stacktrace":
			if in.IsNull() {
				in.Skip()
				out.Stacktrace = nil
			} else {
				in.Delim('[')
				if out.Stacktrace == nil {
					if !in.IsDelim(']') {
						out.Stacktrace = make(LogRecordStacktrace, 0, 8)
					} else {
						out.Stacktrace = LogRecordStacktrace{}
					}
				} else {
					out.Stacktrace = (out.Stacktrace)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *StacktraceFrame
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(StacktraceFrame)
						}
						easyjsonE34310f8DecodeGithubComElasticApmServerModels4(in, &*v1)
					}
					out.Stacktrace = append(out.Stacktrace, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE34310f8EncodeGithubComElasticApmServerModels2(out *jwriter.Writer, in LogRecord) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Level != nil {
		const prefix string = ",\"level\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(*in.Level))
	}
	if in.LoggerName != nil {
		const prefix string = ",\"logger_name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(*in.LoggerName))
	}
	{
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Message == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Message))
		}
	}
	if in.ParamMessage != "" {
		const prefix string = ",\"param_message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ParamMessage))
	}
	{
		const prefix string = ",\"stacktrace\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Stacktrace == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Stacktrace {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					easyjsonE34310f8EncodeGithubComElasticApmServerModels4(out, *v3)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
func easyjsonE34310f8DecodeGithubComElasticApmServerModels4(in *jlexer.Lexer, out *StacktraceFrame) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "abs_path":
			out.AbsPath = string(in.String())
		case "colno":
			out.Colno = int64(in.Int64())
		case "context_line":
			out.ContextLine = string(in.String())
		case "filename":
			if in.IsNull() {
				in.Skip()
				out.Filename = nil
			} else {
				if out.Filename == nil {
					out.Filename = new(string)
				}
				*out.Filename = string(in.String())
			}
		case "function":
			out.Function = string(in.String())
		case "library_frame":
			out.LibraryFrame = bool(in.Bool())
		case "lineno":
			if in.IsNull() {
				in.Skip()
				out.Lineno = nil
			} else {
				if out.Lineno == nil {
					out.Lineno = new(int64)
				}
				*out.Lineno = int64(in.Int64())
			}
		case "module":
			out.Module = string(in.String())
		case "post_context":
			if in.IsNull() {
				in.Skip()
				out.PostContext = nil
			} else {
				in.Delim('[')
				if out.PostContext == nil {
					if !in.IsDelim(']') {
						out.PostContext = make([]string, 0, 4)
					} else {
						out.PostContext = []string{}
					}
				} else {
					out.PostContext = (out.PostContext)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.PostContext = append(out.PostContext, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "pre_context":
			if in.IsNull() {
				in.Skip()
				out.PreContext = nil
			} else {
				in.Delim('[')
				if out.PreContext == nil {
					if !in.IsDelim(']') {
						out.PreContext = make([]string, 0, 4)
					} else {
						out.PreContext = []string{}
					}
				} else {
					out.PreContext = (out.PreContext)[:0]
				}
				for !in.IsDelim(']') {
					var v5 string
					v5 = string(in.String())
					out.PreContext = append(out.PreContext, v5)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "vars":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Vars = make(map[string]interface{})
				} else {
					out.Vars = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v6 interface{}
					if m, ok := v6.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v6.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v6 = in.Interface()
					}
					(out.Vars)[key] = v6
					in.WantComma()
				}
				in.Delim('}')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE34310f8EncodeGithubComElasticApmServerModels4(out *jwriter.Writer, in StacktraceFrame) {
	out.RawByte('{')
	first := true
	_ = first
	if in.AbsPath != "" {
		const prefix string = ",\"abs_path\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.AbsPath))
	}
	if in.Colno != 0 {
		const prefix string = ",\"colno\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Colno))
	}
	if in.ContextLine != "" {
		const prefix string = ",\"context_line\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ContextLine))
	}
	{
		const prefix string = ",\"filename\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Filename == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Filename))
		}
	}
	if in.Function != "" {
		const prefix string = ",\"function\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Function))
	}
	if in.LibraryFrame {
		const prefix string = ",\"library_frame\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.LibraryFrame))
	}
	{
		const prefix string = ",\"lineno\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Lineno == nil {
			out.RawString("null")
		} else {
			out.Int64(int64(*in.Lineno))
		}
	}
	if in.Module != "" {
		const prefix string = ",\"module\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Module))
	}
	{
		const prefix string = ",\"post_context\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.PostContext == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v7, v8 := range in.PostContext {
				if v7 > 0 {
					out.RawByte(',')
				}
				out.String(string(v8))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"pre_context\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.PreContext == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v9, v10 := range in.PreContext {
				if v9 > 0 {
					out.RawByte(',')
				}
				out.String(string(v10))
			}
			out.RawByte(']')
		}
	}
	if len(in.Vars) != 0 {
		const prefix string = ",\"vars\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v11First := true
			for v11Name, v11Value := range in.Vars {
				if v11First {
					v11First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v11Name))
				out.RawByte(':')
				if m, ok := v11Value.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v11Value.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v11Value))
				}
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}
func easyjsonE34310f8DecodeGithubComElasticApmServerModels1(in *jlexer.Lexer, out *Exception) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "attributes":
			if m, ok := out.Attributes.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Attributes.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Attributes = in.Interface()
			}
		case "code":
			if m, ok := out.Code.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Code.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Code = in.Interface()
			}
		case "handled":
			out.Handled = bool(in.Bool())
		case "message":
			if in.IsNull() {
				in.Skip()
				out.Message = nil
			} else {
				if out.Message == nil {
					out.Message = new(string)
				}
				*out.Message = string(in.String())
			}
		case "module":
			out.Module = string(in.String())
		case "stacktrace":
			if in.IsNull() {
				in.Skip()
				out.Stacktrace = nil
			} else {
				in.Delim('[')
				if out.Stacktrace == nil {
					if !in.IsDelim(']') {
						out.Stacktrace = make(ExceptionStacktrace, 0, 8)
					} else {
						out.Stacktrace = ExceptionStacktrace{}
					}
				} else {
					out.Stacktrace = (out.Stacktrace)[:0]
				}
				for !in.IsDelim(']') {
					var v12 *StacktraceFrame
					if in.IsNull() {
						in.Skip()
						v12 = nil
					} else {
						if v12 == nil {
							v12 = new(StacktraceFrame)
						}
						easyjsonE34310f8DecodeGithubComElasticApmServerModels4(in, &*v12)
					}
					out.Stacktrace = append(out.Stacktrace, v12)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "type":
			out.Type = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE34310f8EncodeGithubComElasticApmServerModels1(out *jwriter.Writer, in Exception) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Attributes != nil {
		const prefix string = ",\"attributes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if m, ok := in.Attributes.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Attributes.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Attributes))
		}
	}
	if in.Code != nil {
		const prefix string = ",\"code\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if m, ok := in.Code.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Code.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Code))
		}
	}
	if in.Handled {
		const prefix string = ",\"handled\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Handled))
	}
	{
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Message == nil {
			out.RawString("null")
		} else {
			out.String(string(*in.Message))
		}
	}
	if in.Module != "" {
		const prefix string = ",\"module\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Module))
	}
	{
		const prefix string = ",\"stacktrace\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Stacktrace == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v13, v14 := range in.Stacktrace {
				if v13 > 0 {
					out.RawByte(',')
				}
				if v14 == nil {
					out.RawString("null")
				} else {
					easyjsonE34310f8EncodeGithubComElasticApmServerModels4(out, *v14)
				}
			}
			out.RawByte(']')
		}
	}
	if in.Type != "" {
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	out.RawByte('}')
}
