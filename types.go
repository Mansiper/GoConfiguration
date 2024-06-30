package configuration

import (
	"reflect"
)

type configReader struct {
	sources []configSource
	options ConfigOptions
	data    configData
}
type configSource struct {
	name     string
	value    string
	ft       formatType
	fromFile bool
}
type configData struct {
	currentLine int
	currentPos  int
	currentFile string

	initErrors []string
}

type Parser func(string) (interface{}, error)

type ConfigOptions struct {
	// Rewrite values (not for slice values), default is true
	RewriteValues bool
	// Custom parsers for user types (key - parser name, value - parser)
	Parsers map[string]Parser
}

type intermediateTree map[string][]intermediateData
type intermediateData struct {
	source    int
	value     interface{}
	valueType valueType
}

type formatType int
type valueType int

type structInfo struct {
	fieldName  string
	fieldType  reflect.Type
	field      reflect.Value
	keyName    string
	defValue   string
	isRequired bool
	useParser  bool
	separator  string
	separator2 string //map separator
	isSlice    bool
	isMap      bool
	isPointer  bool
	append     bool
	size       int
}

type jsonTempData struct {
	prefix     string
	isRoot     bool
	isObject   bool
	parseState int
	foundInfo  structInfo
}
