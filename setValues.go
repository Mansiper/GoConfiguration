package configuration

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

const (
	intName      = "an integer"
	uintName     = "an unsigned integer"
	floatName    = "a float"
	boolName     = "a boolean"
	timeName     = "a time"
	durationName = "a duration"
	stringName   = "a string"
)

func (cr *configReader) setFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	var err error = nil
	if info.isPointer && str == nilDefault {
		info.field.SetZero()
		return nil
	}
	if info.useParser {
		if str == nilDefault {
			info.field.Set(reflect.Zero(info.fieldType))
			return nil
		}
		if parser, ok := cr.options.Parsers[info.keyName]; ok {
			intfc, err := parser(str)
			if err != nil {
				return err
			}
			info.field.Set(reflect.ValueOf(intfc))
			return nil
		} else {
			return errors.New("parser not found for key " + info.keyName)
		}
	}

	switch info.fieldType.Kind() {
	case reflect.String:
		err = cr.setStringFieldValue(info, str, strSlice, strMap, vType)
	case reflect.Bool:
		err = cr.setBoolFieldValue(info, str, strSlice, strMap, vType)
	case reflect.Int:
		if vType != vtNumber && vType != vtAny && vType != vtNull {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (int, error) {
			return strconv.Atoi(s)
		}, intName)
	case reflect.Int8:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (int8, error) {
			i, err := strconv.ParseInt(s, 10, 8)
			return int8(i), err
		}, intName)
	case reflect.Int16:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (int16, error) {
			i, err := strconv.ParseInt(s, 10, 16)
			return int16(i), err
		}, intName)
	case reflect.Int32:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (int32, error) {
			i, err := strconv.ParseInt(s, 10, 32)
			return int32(i), err
		}, intName)
	case reflect.Int64:
		err = cr.setInt64FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Uint:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (uint, error) {
			i, err := strconv.ParseUint(s, 10, 64)
			return uint(i), err
		}, uintName)
	case reflect.Uint8:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (uint8, error) {
			i, err := strconv.ParseUint(s, 10, 8)
			return uint8(i), err
		}, uintName)
	case reflect.Uint16:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (uint16, error) {
			i, err := strconv.ParseUint(s, 10, 16)
			return uint16(i), err
		}, uintName)
	case reflect.Uint32:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (uint32, error) {
			i, err := strconv.ParseUint(s, 10, 32)
			return uint32(i), err
		}, uintName)
	case reflect.Uint64:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (uint64, error) {
			return strconv.ParseUint(s, 10, 64)
		}, uintName)
	case reflect.Float32:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, floatName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (float32, error) {
			f, err := strconv.ParseFloat(s, 32)
			return float32(f), err
		}, floatName)
	case reflect.Float64:
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, floatName)
		}
		err = setNumericField(info, str, strSlice, strMap, func(s string) (float64, error) {
			return strconv.ParseFloat(s, 64)
		}, floatName)
	case reflect.Struct:
		err = cr.setStructFieldValue(info, str, strSlice, strMap, vType)
	default:
		err = errors.New("unsupported field type " + info.fieldType.String())
	}
	return err
}

// setNumericField is a generic helper that handles dynamic slices, fixed arrays,
// maps, and scalar fields for any type T that can be parsed from a string.
func setNumericField[T any](info structInfo, str string, strSlice []string, strMap map[string]string, parse func(string) (T, error), typeName string) error {
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			slice := []T{}
			for index, s := range strSlice {
				v, err := parse(s)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, typeName)
				}
				slice = append(slice, v)
			}
			info.field.Set(reflect.ValueOf(slice))
		} else {
			slice := []*T{}
			for index, s := range strSlice {
				if s == nilDefault {
					slice = append(slice, nil)
				} else {
					v, err := parse(s)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, typeName)
					}
					slice = append(slice, &v)
				}
			}
			info.field.Set(reflect.ValueOf(slice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				v, err := parse(s)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, typeName)
				}
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&v))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(v))
				}
			}
		}
	} else if info.isMap {
		m := map[string]T{}
		for key, value := range strMap {
			v, err := parse(value)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			m[key] = v
		}
		info.field.Set(reflect.ValueOf(m))
	} else {
		v, err := parse(str)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, typeName)
		}
		if info.isPointer {
			info.field.Set(reflect.ValueOf(&v))
		} else {
			info.field.Set(reflect.ValueOf(v))
		}
	}
	return nil
}

func (cr *configReader) setStringFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtString && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, stringName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			info.field.Set(reflect.ValueOf(strSlice))
		} else {
			strSlicePtr := []*string{}
			for _, s := range strSlice {
				if s == nilDefault {
					strSlicePtr = append(strSlicePtr, nil)
				} else {
					strSlicePtr = append(strSlicePtr, &s)
				}
			}
			info.field.Set(reflect.ValueOf(strSlicePtr))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if !info.isPointer {
				info.field.Index(index).SetString(s)
			} else {
				if s == nilDefault {
					info.field.Index(index).SetZero()
				} else {
					info.field.Index(index).Set(reflect.ValueOf(&s))
				}
			}
		}
	} else if info.isMap {
		info.field.Set(reflect.ValueOf(strMap))
	} else {
		if info.isPointer {
			info.field.Set(reflect.ValueOf(&str))
		} else {
			info.field.SetString(str)
		}
	}
	return nil
}

func (cr *configReader) setBoolFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtBool && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, boolName)
	}
	return setNumericField(info, str, strSlice, strMap, strconv.ParseBool, boolName)
}

func (cr *configReader) setInt64FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if info.fieldType.String() == "time.Duration" {
		return cr.setDurationFieldValue(info, str, strSlice, strMap, vType)
	}
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	return setNumericField(info, str, strSlice, strMap, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	}, intName)
}

func (cr *configReader) setTimeFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtString && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, timeName)
	}
	parseTime := func(s string) (time.Time, error) {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil && s == nowTime {
			return time.Now(), nil
		}
		return t, err
	}
	return setNumericField(info, str, strSlice, strMap, parseTime, timeName)
}

func (cr *configReader) setDurationFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtString && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, timeName)
	}
	return setNumericField(info, str, strSlice, strMap, time.ParseDuration, durationName)
}

func (cr *configReader) setStructFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if info.fieldType.String() == "time.Time" {
		return cr.setTimeFieldValue(info, str, strSlice, strMap, vType)
	} else {
		return errors.New("unsupported type " + info.fieldType.String())
	}
}

func getValueIsNotTypeError(fieldName string, index int, typeName string) error {
	if index == -1 {
		return errors.New("field " + fieldName + " is not " + typeName)
	} else {
		return getValueIsNotTypeErrorByKey(fieldName, strconv.Itoa(index), typeName)
	}
}

func getValueIsNotTypeErrorByKey(fieldName, key, typeName string) error {
	return errors.New("field " + fieldName + "[" + key + "] is not " + typeName)
}
