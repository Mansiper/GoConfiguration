package configuration

import (
	"errors"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

func containsRune(runes []rune, r rune) bool {
	for _, rn := range runes {
		if r == rn {
			return true
		}
	}
	return false
}

func (cr *configReader) invalidCharacterError() error {
	source := cr.data.currentFile
	return errors.New("invalid character in \"" + source + "\" " + cr.currentPointInfo())
}

func (cr *configReader) currentPointInfo() string {
	return "(" + strconv.Itoa(cr.data.currentLine) + ":" + strconv.Itoa(cr.data.currentPos) + ")"
}

func unsupportedFileTypeError(fileName string) string {
	return "unsupported file type: " + fileName
}

func addValue(structInfo structInfo, it intermediateTree, name, value, key string, isSlice bool, sourceId int) {
	value = strings.Trim(value, " \t")
	if structInfo.isSlice {
		if _, ok := it[name]; !ok {
			it[name] = []intermediateData{{source: sourceId, value: []string{}, valueType: vtAny}}
		} else if !slices.ContainsFunc(it[name], func(data intermediateData) bool { return data.source == sourceId }) {
			it[name] = append(it[name], intermediateData{source: sourceId, value: []string{}, valueType: vtAny})
		}
		for i, v := range it[name] {
			if v.source == sourceId {
				if isSlice {
					it[name][i].value = append(it[name][i].value.([]string), value)
				} else {
					it[name][i].value = append(it[name][i].value.([]string), strings.Split(value, structInfo.separator)...)
				}
				break
			}
		}
	} else if structInfo.isMap && key != "" {
		if _, ok := it[name]; !ok {
			it[name] = []intermediateData{{source: sourceId, value: map[string]string{}, valueType: vtAny}}
		} else if !slices.ContainsFunc(it[name], func(data intermediateData) bool { return data.source == sourceId }) {
			it[name] = append(it[name], intermediateData{source: sourceId, value: map[string]string{}, valueType: vtAny})
		}
		for i, v := range it[name] {
			if v.source == sourceId {
				it[name][i].value.(map[string]string)[key] = value
			}
		}
	} else {
		if _, ok := it[name]; !ok {
			it[name] = []intermediateData{{source: sourceId, value: value, valueType: vtAny}}
		} else {
			it[name] = append(it[name], intermediateData{source: sourceId, value: value, valueType: vtAny})
		}
	}
}

func (cr *configReader) addJsonValue(structInfo structInfo, it intermediateTree, name, value, key string, vType valueType, sourceId int) error {
	if structInfo.isSlice {
		if vType == vtString && value == nilDefault {
			vType = vtNull
		}
		if _, ok := it[name]; !ok {
			it[name] = []intermediateData{{source: sourceId, value: []string{}, valueType: vType}}
		} else if !slices.ContainsFunc(it[name], func(data intermediateData) bool { return data.source == sourceId }) {
			it[name] = append(it[name], intermediateData{source: sourceId, value: []string{}, valueType: vType})
		}
		for i, v := range it[name] {
			if v.valueType == vtNull && vType != vtNull {
				it[name][i].valueType = vType
				v.valueType = vType
			}
			if v.valueType != vType && vType != vtNull {
				return errors.New("different value types in slice \"" + name + "\" " + cr.currentPointInfo())
			}
			if v.source == sourceId {
				it[name][i].value = append(it[name][i].value.([]string), value)
				break
			}
		}
	} else if structInfo.isMap {
		if _, ok := it[name]; !ok {
			it[name] = []intermediateData{{source: sourceId, value: map[string]string{}, valueType: vType}}
		} else if !slices.ContainsFunc(it[name], func(data intermediateData) bool { return data.source == sourceId }) {
			it[name] = append(it[name], intermediateData{source: sourceId, value: map[string]string{}, valueType: vType})
		}
		for i, v := range it[name] {
			if v.source == sourceId {
				it[name][i].value.(map[string]string)[key] = value
			}
		}
	} else {
		if _, ok := it[name]; !ok {
			it[name] = []intermediateData{{source: sourceId, value: value, valueType: vType}}
		} else {
			it[name] = append(it[name], intermediateData{source: sourceId, value: value, valueType: vType})
		}
	}

	return nil
}

func getPointerFieldType(fieldType reflect.Type) (reflect.Type, error) {
	switch fieldType.String() {
	case "*int":
		fieldType = reflect.TypeOf(int(0))
	case "*int8":
		fieldType = reflect.TypeOf(int8(0))
	case "*int16":
		fieldType = reflect.TypeOf(int16(0))
	case "*int32":
		fieldType = reflect.TypeOf(int32(0))
	case "*int64":
		fieldType = reflect.TypeOf(int64(0))
	case "*uint":
		fieldType = reflect.TypeOf(uint(0))
	case "*uint8":
		fieldType = reflect.TypeOf(uint8(0))
	case "*uint16":
		fieldType = reflect.TypeOf(uint16(0))
	case "*uint32":
		fieldType = reflect.TypeOf(uint32(0))
	case "*uint64":
		fieldType = reflect.TypeOf(uint64(0))
	case "*float32":
		fieldType = reflect.TypeOf(float32(0))
	case "*float64":
		fieldType = reflect.TypeOf(float64(0))
	case "*string":
		fieldType = reflect.TypeOf("")
	case "*bool":
		fieldType = reflect.TypeOf(false)
	case "*time.Time":
		fieldType = reflect.TypeOf(time.Time{})
	case "*time.Duration":
		fieldType = reflect.TypeOf(time.Duration(0))
	default:
		return nil, errors.New("unsupported type " + fieldType.String())
	}

	return fieldType, nil
}

func (cr *configReader) processEofError(err error) error {
	if err == io.EOF {
		return errors.New("unexpected end of file " + cr.currentPointInfo())
	}
	return err
}

func (cr *configReader) processNamedError(err error, source string) error {
	return errors.New("error in " + source + " \"" + cr.data.currentFile + "\": " + err.Error())
}

func (cr *configReader) getJsonPrefix(prefix, name string) string {
	if prefix == "" {
		return name + "."
	} else if strings.HasSuffix(prefix, ".") {
		return prefix + name + "."
	} else {
		return prefix + "." + name + "."
	}
}
