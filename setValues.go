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
		err = cr.setIntFieldValue(info, str, strSlice, strMap, vType)
	case reflect.Int8:
		err = cr.setInt8FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Int16:
		err = cr.setInt16FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Int32:
		err = cr.setInt32FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Int64:
		err = cr.setInt64FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Uint:
		err = cr.setUintFieldValue(info, str, strSlice, strMap, vType)
	case reflect.Uint8:
		err = cr.setUint8FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Uint16:
		err = cr.setUint16FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Uint32:
		err = cr.setUint32FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Uint64:
		err = cr.setUint64FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Float32:
		err = cr.setFloat32FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Float64:
		err = cr.setFloat64FieldValue(info, str, strSlice, strMap, vType)
	case reflect.Struct:
		err = cr.setStructFieldValue(info, str, strSlice, strMap, vType)
	default:
		err = errors.New("unsupported field type " + info.fieldType.String())
	}
	return err
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
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			bSlice := []bool{}
			for index, s := range strSlice {
				b, err := strconv.ParseBool(s)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, boolName)
				}
				bSlice = append(bSlice, b)
			}
			info.field.Set(reflect.ValueOf(bSlice))
		} else {
			bSlice := []*bool{}
			for index, s := range strSlice {
				if s == nilDefault {
					bSlice = append(bSlice, nil)
				} else {
					b, err := strconv.ParseBool(s)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, boolName)
					}
					bSlice = append(bSlice, &b)
				}
			}
			info.field.Set(reflect.ValueOf(bSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				b, err := strconv.ParseBool(s)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, boolName)
				}
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&b))
				} else {
					info.field.Index(index).SetBool(b)
				}
			}
		}
	} else if info.isMap {
		bMap := map[string]bool{}
		for key, value := range strMap {
			b, err := strconv.ParseBool(value)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			bMap[key] = b
		}
		info.field.Set(reflect.ValueOf(bMap))
	} else {
		b, err := strconv.ParseBool(str)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, boolName)
		}
		if info.isPointer {
			info.field.Set(reflect.ValueOf(&b))
		} else {
			info.field.SetBool(b)
		}
	}
	return nil
}

func (cr *configReader) setIntFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny && vType != vtNull {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			iSlice := []int{}
			for index, s := range strSlice {
				i, err := strconv.Atoi(s)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				iSlice = append(iSlice, i)
			}
			info.field.Set(reflect.ValueOf(iSlice))
		} else {
			iSlice := []*int{}
			for index, s := range strSlice {
				if s == nilDefault {
					iSlice = append(iSlice, nil)
				} else {
					i, err := strconv.Atoi(s)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					iSlice = append(iSlice, &i)
				}
			}
			info.field.Set(reflect.ValueOf(iSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				i, err := strconv.Atoi(s)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&i))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(i))
				}
			}
		}
	} else if info.isMap {
		iMap := map[string]int{}
		for key, value := range strMap {
			i, err := strconv.Atoi(value)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			iMap[key] = i
		}
		info.field.Set(reflect.ValueOf(iMap))
	} else {
		i, err := strconv.Atoi(str)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		if info.isPointer {
			info.field.Set(reflect.ValueOf(&i))
		} else {
			info.field.Set(reflect.ValueOf(i))
		}
	}
	return nil
}

func (cr *configReader) setInt8FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			iSlice := []int8{}
			for index, s := range strSlice {
				i, err := strconv.ParseInt(s, 10, 8)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				iSlice = append(iSlice, int8(i))
			}
			info.field.Set(reflect.ValueOf(iSlice))
		} else {
			iSlice := []*int8{}
			for index, s := range strSlice {
				if s == nilDefault {
					iSlice = append(iSlice, nil)
				} else {
					i, err := strconv.ParseInt(s, 10, 8)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					i8 := int8(i)
					iSlice = append(iSlice, &i8)
				}
			}
			info.field.Set(reflect.ValueOf(iSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				i, err := strconv.ParseInt(s, 10, 8)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				i8 := int8(i)
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&i8))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(i8))
				}
			}
		}
	} else if info.isMap {
		iMap := map[string]int8{}
		for key, value := range strMap {
			i, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			iMap[key] = int8(i)
		}
		info.field.Set(reflect.ValueOf(iMap))
	} else {
		i, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		if info.isPointer {
			i8 := int8(i)
			info.field.Set(reflect.ValueOf(&i8))
		} else {
			info.field.SetInt(i)
		}
	}
	return nil
}

func (cr *configReader) setInt16FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			iSlice := []int16{}
			for index, s := range strSlice {
				i, err := strconv.ParseInt(s, 10, 16)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				iSlice = append(iSlice, int16(i))
			}
			info.field.Set(reflect.ValueOf(iSlice))
		} else {
			iSlice := []*int16{}
			for index, s := range strSlice {
				if s == nilDefault {
					iSlice = append(iSlice, nil)
				} else {
					i, err := strconv.ParseInt(s, 10, 16)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					i16 := int16(i)
					iSlice = append(iSlice, &i16)
				}
			}
			info.field.Set(reflect.ValueOf(iSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				i, err := strconv.ParseInt(s, 10, 16)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				i16 := int16(i)
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&i16))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(i16))
				}
			}
		}
	} else if info.isMap {
		iMap := map[string]int16{}
		for key, value := range strMap {
			i, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			iMap[key] = int16(i)
		}
		info.field.Set(reflect.ValueOf(iMap))
	} else {
		i, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		if info.isPointer {
			i16 := int16(i)
			info.field.Set(reflect.ValueOf(&i16))
		} else {
			info.field.SetInt(i)
		}
	}
	return nil
}

func (cr *configReader) setInt32FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			iSlice := []int32{}
			for index, s := range strSlice {
				i, err := strconv.ParseInt(s, 10, 32)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				iSlice = append(iSlice, int32(i))
			}
			info.field.Set(reflect.ValueOf(iSlice))
		} else {
			iSlice := []*int32{}
			for index, s := range strSlice {
				if s == nilDefault {
					iSlice = append(iSlice, nil)
				} else {
					i, err := strconv.ParseInt(s, 10, 32)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					i32 := int32(i)
					iSlice = append(iSlice, &i32)
				}
			}
			info.field.Set(reflect.ValueOf(iSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				i, err := strconv.ParseInt(s, 10, 32)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				i32 := int32(i)
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&i32))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(i32))
				}
			}
		}
	} else if info.isMap {
		iMap := map[string]int32{}
		for key, value := range strMap {
			i, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			iMap[key] = int32(i)
		}
		info.field.Set(reflect.ValueOf(iMap))
	} else {
		i, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		if info.isPointer {
			i32 := int32(i)
			info.field.Set(reflect.ValueOf(&i32))
		} else {
			info.field.SetInt(i)
		}
	}
	return nil
}

func (cr *configReader) setInt64FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if info.fieldType.String() == "time.Duration" {
		return cr.setDurationFieldValue(info, str, strSlice, strMap, vType)
	} else {
		if vType != vtNumber && vType != vtAny {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		if info.isSlice && info.size == 0 {
			if !info.isPointer {
				iSlice := []int64{}
				for index, s := range strSlice {
					i, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					iSlice = append(iSlice, i)
				}
				info.field.Set(reflect.ValueOf(iSlice))
			} else {
				iSlice := []*int64{}
				for index, s := range strSlice {
					if s == nilDefault {
						iSlice = append(iSlice, nil)
					} else {
						i, err := strconv.ParseInt(s, 10, 64)
						if err != nil {
							return getValueIsNotTypeError(info.fieldName, index, intName)
						}
						iSlice = append(iSlice, &i)
					}
				}
				info.field.Set(reflect.ValueOf(iSlice))
			}
		} else if info.isSlice && info.size > 0 {
			for index, s := range strSlice {
				if info.isPointer && s == nilDefault {
					info.field.Index(index).SetZero()
				} else {
					i, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					if info.isPointer {
						info.field.Index(index).Set(reflect.ValueOf(&i))
					} else {
						info.field.Index(index).Set(reflect.ValueOf(i))
					}
				}
			}
		} else if info.isMap {
			iMap := map[string]int64{}
			for key, value := range strMap {
				i, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
				}
				iMap[key] = i
			}
			info.field.Set(reflect.ValueOf(iMap))
		} else {
			i, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return getValueIsNotTypeError(info.fieldName, -1, intName)
			}
			if info.isPointer {
				info.field.Set(reflect.ValueOf(&i))
			} else {
				info.field.SetInt(i)
			}
		}
	}
	return nil
}

func (cr *configReader) setUintFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			iSlice := []uint{}
			for index, s := range strSlice {
				i, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, uintName)
				}
				iSlice = append(iSlice, uint(i))
			}
			info.field.Set(reflect.ValueOf(iSlice))
		} else {
			iSlice := []*uint{}
			for index, s := range strSlice {
				if s == nilDefault {
					iSlice = append(iSlice, nil)
				} else {
					i, err := strconv.ParseUint(s, 10, 64)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					ui := uint(i)
					iSlice = append(iSlice, &ui)
				}
			}
			info.field.Set(reflect.ValueOf(iSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				i, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, uintName)
				}
				ui := uint(i)
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&ui))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(ui))
				}
			}
		}
	} else if info.isMap {
		iMap := map[string]uint{}
		for key, value := range strMap {
			i, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			iMap[key] = uint(i)
		}
		info.field.Set(reflect.ValueOf(iMap))
	} else {
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, uintName)
		}
		if info.isPointer {
			ii := uint(i)
			info.field.Set(reflect.ValueOf(&ii))
		} else {
			info.field.SetUint(i)
		}
	}
	return nil
}

func (cr *configReader) setUint8FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			iSlice := []uint8{}
			for index, s := range strSlice {
				i, err := strconv.ParseUint(s, 10, 8)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, uintName)
				}
				iSlice = append(iSlice, uint8(i))
			}
			info.field.Set(reflect.ValueOf(iSlice))
		} else {
			iSlice := []*uint8{}
			for index, s := range strSlice {
				if s == nilDefault {
					iSlice = append(iSlice, nil)
				} else {
					i, err := strconv.ParseUint(s, 10, 8)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					i8 := uint8(i)
					iSlice = append(iSlice, &i8)
				}
			}
			info.field.Set(reflect.ValueOf(iSlice))

		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				i, err := strconv.ParseUint(s, 10, 8)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				i8 := uint8(i)
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&i8))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(i8))
				}
			}
		}
	} else if info.isMap {
		iMap := map[string]uint8{}
		for key, value := range strMap {
			i, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			iMap[key] = uint8(i)
		}
		info.field.Set(reflect.ValueOf(iMap))
	} else {
		i, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, uintName)
		}
		if info.isPointer {
			i8 := uint8(i)
			info.field.Set(reflect.ValueOf(&i8))
		} else {
			info.field.SetUint(i)
		}
	}
	return nil
}

func (cr *configReader) setUint16FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			iSlice := []uint16{}
			for index, s := range strSlice {
				i, err := strconv.ParseUint(s, 10, 16)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, uintName)
				}
				iSlice = append(iSlice, uint16(i))
			}
			info.field.Set(reflect.ValueOf(iSlice))
		} else {
			iSlice := []*uint16{}
			for index, s := range strSlice {
				if s == nilDefault {
					iSlice = append(iSlice, nil)
				} else {
					i, err := strconv.ParseUint(s, 10, 16)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					i16 := uint16(i)
					iSlice = append(iSlice, &i16)
				}
			}
			info.field.Set(reflect.ValueOf(iSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				i, err := strconv.ParseUint(s, 10, 16)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				i16 := uint16(i)
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&i16))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(i16))
				}
			}
		}
	} else if info.isMap {
		iMap := map[string]uint16{}
		for key, value := range strMap {
			i, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			iMap[key] = uint16(i)
		}
		info.field.Set(reflect.ValueOf(iMap))
	} else {
		i, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, uintName)
		}
		if info.isPointer {
			i16 := uint16(i)
			info.field.Set(reflect.ValueOf(&i16))
		} else {
			info.field.SetUint(i)
		}
	}
	return nil
}

func (cr *configReader) setUint32FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			iSlice := []uint32{}
			for index, s := range strSlice {
				i, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, uintName)
				}
				iSlice = append(iSlice, uint32(i))
			}
			info.field.Set(reflect.ValueOf(iSlice))
		} else {
			iSlice := []*uint32{}
			for index, s := range strSlice {
				if s == nilDefault {
					iSlice = append(iSlice, nil)
				} else {
					i, err := strconv.ParseUint(s, 10, 32)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					i32 := uint32(i)
					iSlice = append(iSlice, &i32)
				}
			}
			info.field.Set(reflect.ValueOf(iSlice))

		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				i, err := strconv.ParseUint(s, 10, 32)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				i32 := uint32(i)
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&i32))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(i32))
				}
			}
		}
	} else if info.isMap {
		iMap := map[string]uint32{}
		for key, value := range strMap {
			i, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			iMap[key] = uint32(i)
		}
		info.field.Set(reflect.ValueOf(iMap))
	} else {
		i, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, uintName)
		}
		if info.isPointer {
			i32 := uint32(i)
			info.field.Set(reflect.ValueOf(&i32))
		} else {
			info.field.SetUint(i)
		}
	}
	return nil
}

func (cr *configReader) setUint64FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, intName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			iSlice := []uint64{}
			for index, s := range strSlice {
				i, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, uintName)
				}
				iSlice = append(iSlice, i)
			}
			info.field.Set(reflect.ValueOf(iSlice))
		} else {
			iSlice := []*uint64{}
			for index, s := range strSlice {
				if s == nilDefault {
					iSlice = append(iSlice, nil)
				} else {
					i, err := strconv.ParseUint(s, 10, 64)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					iSlice = append(iSlice, &i)
				}
			}
			info.field.Set(reflect.ValueOf(iSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				i, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, uintName)
				}
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&i))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(i))
				}
			}
		}
	} else if info.isMap {
		iMap := map[string]uint64{}
		for key, value := range strMap {
			i, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			iMap[key] = i
		}
		info.field.Set(reflect.ValueOf(iMap))
	} else {
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, uintName)
		}
		if info.isPointer {
			info.field.Set(reflect.ValueOf(&i))
		} else {
			info.field.SetUint(i)
		}
	}
	return nil
}

func (cr *configReader) setFloat32FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, floatName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			fSlice := []float32{}
			for index, s := range strSlice {
				f, err := strconv.ParseFloat(s, 32)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, floatName)
				}
				fSlice = append(fSlice, float32(f))
			}
			info.field.Set(reflect.ValueOf(fSlice))
		} else {
			fSlice := []*float32{}
			for index, s := range strSlice {
				if s == nilDefault {
					fSlice = append(fSlice, nil)
				} else {
					f, err := strconv.ParseFloat(s, 32)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, floatName)
					}
					f32 := float32(f)
					fSlice = append(fSlice, &f32)
				}
			}
			info.field.Set(reflect.ValueOf(fSlice))

		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				f, err := strconv.ParseFloat(s, 32)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, floatName)
				}
				f32 := float32(f)
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&f32))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(f32))
				}
			}
		}
	} else if info.isMap {
		fMap := map[string]float32{}
		for key, value := range strMap {
			f, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			fMap[key] = float32(f)
		}
		info.field.Set(reflect.ValueOf(fMap))
	} else {
		f, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, floatName)
		}
		if info.isPointer {
			f32 := float32(f)
			info.field.Set(reflect.ValueOf(&f32))
		} else {
			info.field.SetFloat(f)
		}
	}
	return nil
}

func (cr *configReader) setFloat64FieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtNumber && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, floatName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			fSlice := []float64{}
			for index, s := range strSlice {
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, floatName)
				}
				fSlice = append(fSlice, f)
			}
			info.field.Set(reflect.ValueOf(fSlice))
		} else {
			fSlice := []*float64{}
			for index, s := range strSlice {
				if s == nilDefault {
					fSlice = append(fSlice, nil)
				} else {
					f, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, floatName)
					}
					fSlice = append(fSlice, &f)
				}
			}
			info.field.Set(reflect.ValueOf(fSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, floatName)
				}
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&f))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(f))
				}
			}
		}
	} else if info.isMap {
		fMap := map[string]float64{}
		for key, value := range strMap {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			fMap[key] = f
		}
		info.field.Set(reflect.ValueOf(fMap))
	} else {
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, floatName)
		}
		if info.isPointer {
			info.field.Set(reflect.ValueOf(&f))
		} else {
			info.field.SetFloat(f)
		}
	}
	return nil
}

func (cr *configReader) setTimeFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtString && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, timeName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			tSlice := []time.Time{}
			for index, s := range strSlice {
				t, err := time.Parse(time.RFC3339, s)
				if err != nil && s == nowTime {
					t = time.Now()
				} else if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, timeName)
				}
				tSlice = append(tSlice, t)
			}
			info.field.Set(reflect.ValueOf(tSlice))
		} else {
			tSlice := []*time.Time{}
			for index, s := range strSlice {
				if s == nilDefault {
					tSlice = append(tSlice, nil)
				} else {
					t, err := time.Parse(time.RFC3339, s)
					if err != nil && s == nowTime {
						t = time.Now()
					} else if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, timeName)
					}
					tSlice = append(tSlice, &t)
				}
			}
			info.field.Set(reflect.ValueOf(tSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				t, err := time.Parse(time.RFC3339, s)
				if err != nil && s == nowTime {
					t = time.Now()
				} else if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, timeName)
				}
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&t))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(t))
				}
			}
		}
	} else if info.isMap {
		tMap := map[string]time.Time{}
		for key, value := range strMap {
			t, err := time.Parse(time.RFC3339, value)
			if err != nil && value == nowTime {
				t = time.Now()
			} else if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			tMap[key] = t
		}
		info.field.Set(reflect.ValueOf(tMap))
	} else {
		t, err := time.Parse(time.RFC3339, str)
		if err != nil && str == nowTime {
			t = time.Now()
		} else if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, timeName)
		}
		if info.isPointer {
			info.field.Set(reflect.ValueOf(&t))
		} else {
			info.field.Set(reflect.ValueOf(t))
		}
	}
	return nil
}

func (cr *configReader) setDurationFieldValue(info structInfo, str string, strSlice []string, strMap map[string]string, vType valueType) error {
	if vType != vtString && vType != vtAny {
		return getValueIsNotTypeError(info.fieldName, -1, timeName)
	}
	if info.isSlice && info.size == 0 {
		if !info.isPointer {
			dSlice := []time.Duration{}
			for index, s := range strSlice {
				d, err := time.ParseDuration(s)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				dSlice = append(dSlice, d)
			}
			info.field.Set(reflect.ValueOf(dSlice))
		} else {
			dSlice := []*time.Duration{}
			for index, s := range strSlice {
				if s == nilDefault {
					dSlice = append(dSlice, nil)
				} else {
					d, err := time.ParseDuration(s)
					if err != nil {
						return getValueIsNotTypeError(info.fieldName, index, intName)
					}
					dSlice = append(dSlice, &d)
				}
			}
			info.field.Set(reflect.ValueOf(dSlice))
		}
	} else if info.isSlice && info.size > 0 {
		for index, s := range strSlice {
			if info.isPointer && s == nilDefault {
				info.field.Index(index).SetZero()
			} else {
				d, err := time.ParseDuration(s)
				if err != nil {
					return getValueIsNotTypeError(info.fieldName, index, intName)
				}
				if info.isPointer {
					info.field.Index(index).Set(reflect.ValueOf(&d))
				} else {
					info.field.Index(index).Set(reflect.ValueOf(d))
				}
			}
		}
	} else if info.isMap {
		dMap := map[string]time.Duration{}
		for key, value := range strMap {
			d, err := time.ParseDuration(value)
			if err != nil {
				return getValueIsNotTypeErrorByKey(info.fieldName, key, intName)
			}
			dMap[key] = d
		}
		info.field.Set(reflect.ValueOf(dMap))
	} else {
		d, err := time.ParseDuration(str)
		if err != nil {
			return getValueIsNotTypeError(info.fieldName, -1, intName)
		}
		if info.isPointer {
			info.field.Set(reflect.ValueOf(&d))
		} else {
			info.field.Set(reflect.ValueOf(d))
		}
	}
	return nil
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
