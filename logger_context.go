package log

import "context"

const (
	_defaultFieldName = "_context"
)

// ContextParser want to to parse context into `context` field and log into fields.
// it would be called only when `ctx` it not nil
type ContextParser interface {
	// Parse contains the logic code which indicates how to parse the context.
	Parse(ctx context.Context) interface{}

	// FieldName is field which will be used in log fields.
	FieldName() string
}

// NewContextParserFunc parseFunc `func(ctx context.Context) interface{}` as the primary parameter,
// fieldName `string` as secondary parameter.
func NewContextParserFunc(f func(ctx context.Context) interface{}, fieldName string) ContextParser {
	if fieldName == "" {
		fieldName = _defaultFieldName
	}

	return funcContextParser{
		parse:     f,
		fieldName: fieldName,
	}
}

// DefaultContextParserFunc use default field name as output context field name.
func DefaultContextParserFunc(parseFunc func(ctx context.Context) interface{}) ContextParser {
	return NewContextParserFunc(parseFunc, _defaultFieldName)
}

type funcContextParser struct {
	parse     func(ctx context.Context) interface{}
	fieldName string
}

func (f funcContextParser) Parse(ctx context.Context) interface{} {
	return f.parse(ctx)
}

func (f funcContextParser) FieldName() string {
	return f.fieldName
}

// nonParser .
func nonParser(ctx context.Context) interface{} {
	return "non action"
}
