package configuration

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
)

var defaultJsonData = jsonTempData{prefix: "ROOT", isRoot: true, isObject: true}

func (cr *configReader) getFileType(name string) formatType {
	split := strings.Split(name, ".")
	if len(split) > 1 {
		switch strings.ToLower(split[len(split)-1]) {
		case "json":
			return FtJson
		// case "yaml":
		// return FtYaml
		case "ini":
			return FtIni
		case "env":
			return FtEnv
		}
	}
	return ftUnknown
}

// https://stackoverflow.com/questions/37135193/how-to-set-default-values-in-go-structs
// https://go.dev/play/p/rFql2x0Klm4
func (cr *configReader) getStructInfo(userConfig interface{}, fieldPrefix, namePrefix string) ([]structInfo, error) {
	if userConfig == nil {
		return nil, errors.New("user config is nil")
	}
	if reflect.TypeOf(userConfig).Kind() != reflect.Ptr {
		return nil, errors.New("pass your config struct as a pointer")
	}

	info := []structInfo{}

	v := reflect.ValueOf(userConfig).Elem()
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return nil, errors.New("user config must be a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		var field = t.Field(i)
		var fieldType = field.Type

		if fieldType.Kind() == reflect.Struct &&
			fieldType.String() != "time.Time" {
			var keyName string
			var err error
			if envData, ok := cr.hasEnvData(field.Tag.Get("env")); ok {
				useParser := false
				keyName, _, _, useParser, err = cr.getTagData(envData, field)
				if err != nil {
					return nil, err
				}
				if useParser {
					info, err = cr.appendStructInfo(info, fieldType, field, v, i, envData, namePrefix, fieldPrefix)
					if err != nil {
						return nil, err
					} else {
						continue
					}
				}
			}

			subInfo, err := cr.getStructInfo(v.Field(i).Addr().Interface(), fieldPrefix+field.Name+".", namePrefix+keyName+".")
			if err != nil {
				return nil, err
			}
			info = append(info, subInfo...)
		} else if envData, ok := cr.hasEnvData(field.Tag.Get("env")); ok {
			if !v.Field(i).CanSet() {
				continue
			}
			var err error
			info, err = cr.appendStructInfo(info, fieldType, field, v, i, envData, namePrefix, fieldPrefix)
			if err != nil {
				return nil, err
			}
		}
	}

	return info, nil
}

func (cr *configReader) appendStructInfo(info []structInfo,
	fieldType reflect.Type, field reflect.StructField, v reflect.Value,
	i int, envData, namePrefix, fieldPrefix string) ([]structInfo, error) {

	keyName, isRequired, appendToSlice, useParser, err := cr.getTagData(envData, field)
	if err != nil {
		return nil, err
	}

	isPointer := fieldType.Kind() == reflect.Ptr

	sep := field.Tag.Get("sep")
	if sep == "" {
		sep = sepDefault
	}
	sep2 := field.Tag.Get("sep2")
	if sep2 == "" {
		sep2 = sep2Default
	}

	// todo: pointer to struct, slice of struct
	isSlice := fieldType.Kind() == reflect.Slice || (isPointer && strings.Contains(fieldType.String(), "[]"))
	isArray := fieldType.Kind() == reflect.Array
	arraySize := 0
	isMap := fieldType.Kind() == reflect.Map

	if isSlice || isArray {
		if isArray {
			arraySize = fieldType.Len()
		}
		fieldType = fieldType.Elem()
		isPointer = isPointer || fieldType.Kind() == reflect.Ptr
		if isPointer {
			fieldType = fieldType.Elem()
		}
	} else if isMap {
		fieldType = fieldType.Elem()
	}

	if isPointer && !isSlice && !isArray {
		fieldType, err = getPointerFieldType(fieldType)
		if err != nil {
			return nil, err
		}
	}

	def, ok := field.Tag.Lookup("def")
	if !ok {
		if isPointer || isSlice {
			def = nilDefault
		} else {
			def = ""
		}
	}

	info = append(info, structInfo{
		fieldName:  fieldPrefix + field.Name,
		fieldType:  fieldType,
		field:      v.Field(i),
		keyName:    strings.ToLower(namePrefix + keyName),
		defValue:   def,
		isRequired: isRequired,
		useParser:  useParser,
		separator:  sep,
		separator2: sep2,
		isSlice:    isSlice || isArray,
		isMap:      isMap,
		isPointer:  isPointer,
		append:     appendToSlice,
		size:       arraySize,
	})

	return info, nil
}

func (cr *configReader) hasEnvData(envData string) (string, bool) {
	return envData, envData != ignoreField
}

func (cr *configReader) getTagData(envData string, field reflect.StructField) (string, bool, bool, bool, error) {
	keyName := ""
	isRequired := false
	appendToSlice := false
	useParser := false
	if envData == "" {
		keyName = field.Name
	} else {
		split := strings.Split(strings.ToLower(envData), ",")
		if len(split) == 0 || strings.TrimSpace(split[0]) == "" {
			return "", false, false, false, errors.New("env tag is empty for field " + field.Name)
		}
		if strings.ContainsAny(split[0], ".") {
			return "", false, false, false, errors.New("env tag contains invalid characters for field " + field.Name)
		}
		keyName = split[0]
		for _, s := range split[1:] {
			if s == "required" {
				isRequired = true
			} else if s == "append" {
				appendToSlice = true
			} else if s == "useparser" {
				useParser = true
			}
		}
	}

	return keyName, isRequired, appendToSlice, useParser, nil
}

func (cr *configReader) setValues(it intermediateTree, si []structInfo) error {
	for _, info := range si {
		data, ok := it[info.keyName]
		if ok || info.defValue != "" {
			str := ""
			strSlice := []string{}
			strMap := map[string]string{}
			isEMpty := false
			var vType valueType = vtAny
			if ok {
				var value interface{}
				if info.isSlice {
					if info.append {
						for _, d := range data {
							strSlice = append(strSlice, d.value.([]string)...)
							vType = d.valueType
						}
					} else if len(data) > 0 {
						last := data[len(data)-1]
						strSlice = last.value.([]string)
						vType = last.valueType
					}
					if len(strSlice) == 0 && info.isRequired && (info.defValue == "" || info.defValue == nilDefault) {
						return errors.New("required field " + info.fieldName + " is empty")
					}
					if len(strSlice) == 0 || (len(strSlice) == 1 && strSlice[0] == "") {
						isEMpty = true
						if info.defValue == "" {
							strSlice = []string{}
						} else if info.defValue == nilDefault {
							strSlice = nil
						} else {
							strSlice = strings.Split(info.defValue, info.separator)
						}
					}
				} else if info.isMap {
					if info.append {
						for _, d := range data {
							for k, v := range d.value.(map[string]string) {
								strMap[k] = v
								vType = d.valueType
							}
						}
					} else if len(data) > 0 {
						last := data[len(data)-1]
						strMap = last.value.(map[string]string)
						vType = last.valueType
					}
					if len(strMap) == 0 && info.isRequired && (info.defValue == "" || info.defValue == nilDefault) {
						return errors.New("required field " + info.fieldName + " is empty")
					}
					if len(strMap) == 0 || vType == vtNull {
						isEMpty = true
						if info.defValue == "" {
							strMap = map[string]string{}
						} else if info.defValue == nilDefault {
							strMap = nil
						} else {
							split := strings.Split(info.defValue, info.separator)
							for _, s := range split {
								kv := strings.Split(s, info.separator2)
								if len(kv) == 2 {
									strMap[kv[0]] = kv[1]
								} else {
									return errors.New("invalid default value for field " + info.fieldName)
								}
							}
						}
					}
				} else {
					if !cr.options.RewriteValues {
						value = data[0].value
						vType = data[0].valueType
					} else if len(data) > 0 {
						last := data[len(data)-1]
						value = last.value
						vType = last.valueType
					}
					if (value == nil || value.(string) == "") && info.isRequired && info.defValue == "" {
						return errors.New("required field " + info.fieldName + " is empty")
					}
					str = value.(string)
					if str == "" || vType == vtNull {
						isEMpty = true
						str = info.defValue
					}
				}
			} else {
				isEMpty = true
				if info.isSlice {
					if info.isRequired && (info.defValue == "" || info.defValue == nilDefault) {
						return errors.New("required field " + info.fieldName + " value is missing")
					}
					if info.defValue == nilDefault {
						strSlice = nil
					} else if info.defValue == "" {
						strSlice = []string{}
					} else {
						strSlice = strings.Split(info.defValue, info.separator)
					}
				} else if info.isMap {
					if info.isRequired && (info.defValue == "" || info.defValue == nilDefault) {
						return errors.New("required field " + info.fieldName + " value is missing")
					}
					if info.defValue == nilDefault {
						strMap = nil
					} else if info.defValue == "" {
						strMap = map[string]string{}
					} else {
						split := strings.Split(info.defValue, info.separator)
						for _, s := range split {
							kv := strings.Split(s, info.separator2)
							if len(kv) == 2 {
								strMap[kv[0]] = kv[1]
							} else {
								return errors.New("invalid default value for field " + info.fieldName)
							}
						}
					}
				} else {
					str = info.defValue
				}
			}

			if (info.isPointer || info.isSlice) && info.size == 0 && info.defValue == nilDefault && isEMpty ||
				(info.isPointer || info.isSlice) && info.size == 0 && vType == vtNull {
				continue
			} else if info.isSlice && info.size > 0 && isEMpty {
				if len(strSlice) == 0 {
					continue
				} else if len(strSlice) > info.size {
					return errors.New("field " + info.fieldName + " has more values than allowed: " + strings.Join(strSlice, ","))
				}
			} else if info.isMap && isEMpty && len(strMap) == 0 && (info.defValue == "" || info.defValue == nilDefault) {
				continue
			}

			err := cr.setFieldValue(info, str, strSlice, strMap, vType)
			if err != nil {
				return err
			}
		} else if info.isRequired {
			return errors.New("required field " + info.fieldName + " value is missing")
		} else {
			if info.isSlice || info.isMap || info.useParser {
				info.field.Set(reflect.Zero(info.field.Type()))
			}
		}
	}
	return nil
}

func (cr *configReader) readConfigFile(source configSource, it intermediateTree, si []structInfo, sourceId int) error {
	b, err := os.ReadFile(source.value)
	if err != nil {
		return err
	}

	cr.data.currentLine = 1
	cr.data.currentPos = 0
	cr.data.currentFile = source.value
	r := bufio.NewReader(bytes.NewReader(b))

	switch source.ft {
	case FtEnv:
		err = cr.parseEnvData(r, it, si, sourceId)
	case FtJson:
		err = cr.parseJsonData(r, it, si, defaultJsonData, sourceId)
	// case FtYaml:
	// err = cr.parseYamlData(r, it, si, defaultYamlData, sourceId)
	case FtIni:
		err = cr.parseIniData(r, it, si, sourceId)
	default:
		return errors.New(unsupportedFileTypeError(source.value))
	}

	if err != nil {
		err = cr.processNamedError(err, "file")
	}
	return err
}

func (cr *configReader) readConfigString(source configSource, it intermediateTree, si []structInfo, sourceId int) error {
	cr.data.currentLine = 1
	cr.data.currentPos = 0
	cr.data.currentFile = source.name
	r := bufio.NewReader(bytes.NewReader([]byte(source.value)))

	var err error = nil
	switch source.ft {
	case FtEnv:
		err = cr.parseEnvData(r, it, si, sourceId)
	case FtJson:
		err = cr.parseJsonData(r, it, si, defaultJsonData, sourceId)
	// case FtYaml:
	// err = cr.parseYamlData(r, it, si, defaultYamlData, sourceId)
	case FtIni:
		err = cr.parseIniData(r, it, si, sourceId)
	}

	if err != nil {
		err = cr.processNamedError(err, "string")
	}
	return err
}

func (cr *configReader) checkNewLine(rn rune) {
	if rn == '\n' {
		cr.data.currentLine++
		cr.data.currentPos = 0
	}
}

func (cr *configReader) readToNextLine(r *bufio.Reader, allowMultiline bool) error {
	backslash := false
	readToBol := false

	rn := ' '
	var err error = nil
	for {
		if rn, _, err = r.ReadRune(); err == nil {
			cr.data.currentPos++
			if readToBol {
				if rn == '\r' || rn == '\n' {
					cr.checkNewLine(rn)
					continue
				} else {
					readToBol = false
				}
			}

			if rn == '\r' || rn == '\n' {
				cr.checkNewLine(rn)
				if backslash {
					backslash = false
					readToBol = true
					continue
				}
				if rn == '\n' {
					break
				}
			} else if rn == '\\' && allowMultiline {
				backslash = true
				continue
			}
		} else if err != io.EOF {
			return err
		} else {
			break
		}
	}

	return err
}

func (cr *configReader) findFieldByName(r *bufio.Reader, si []structInfo, name string, allowMultiline bool) (bool, structInfo, bool, error) {
	found := false
	foundInfo := structInfo{}
	continue_ := false

	for _, s := range si {
		if s.keyName == name {
			found = true
			foundInfo = s
			break
		}
	}
	if !found {
		if err := cr.readToNextLine(r, allowMultiline); err != nil {
			if err == io.EOF {
				continue_ = false
			}
			return found, foundInfo, false, err
		} else {
			continue_ = true
		}
	}

	return found, foundInfo, continue_, nil
}

func (cr *configReader) findFieldByJsonName(si []structInfo, name string) (bool, structInfo) {
	found := false
	foundInfo := structInfo{}

	for _, s := range si {
		if s.keyName == name || strings.HasPrefix(s.keyName, name+".") {
			found = true
			foundInfo = s
			break
		}
	}

	return found, foundInfo
}

func (cr *configReader) readComment(r *bufio.Reader, singleLine bool) error {
	rn := ' '
	prev := ' '
	var err error = nil
	for {
		if rn, _, err = r.ReadRune(); err == nil {
			cr.data.currentPos++
			if singleLine {
				if rn == '\n' || rn == '\r' {
					cr.checkNewLine(rn)
					if rn == '\n' {
						break
					}
				}
			} else {
				if rn == '/' && prev == '*' {
					break
				} else if rn == '\n' || rn == '\r' {
					cr.checkNewLine(rn)
				}
				prev = rn
			}
		} else {
			return err
		}
	}
	return nil
}

func (cr *configReader) checkDuplicates(name string, it intermediateTree, sourceId int) bool {
	if name == "{" || name == "[" {
		return false
	}

	if _, ok := it[name]; ok {
		for _, d := range it[name] {
			if d.source == sourceId {
				return true
			}
		}
	}
	return false
}

//todo!: json allows name "{" or "[" and now it can be a problem
//todo: useparser for all types
