package configuration

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

type jsonReadValueResult struct {
	value    string
	isOpener bool
	isString bool
	divider  rune
	err      error
}

const (
	psJsonName = iota
	psJsonValue
	psJsonDivider
)

func (cr *configReader) parseJsonData(r *bufio.Reader, it intermediateTree, si []structInfo, data jsonTempData, sourceId int) error {
	name := ""
	divider := ' '
	doBreak := false
	notComma := false
	found := false
	foundInfo := structInfo{}
	var err error = nil

	if !data.isObject {
		data.parseState = psJsonValue
	}

	for {
		switch data.parseState {
		case psJsonName:
			name, err = cr.readJsonName(r)
			if err != nil {
				return cr.processEofError(err)
			}

			isDuplicate := cr.checkDuplicates(data.prefix+name, it, sourceId)
			if isDuplicate {
				return errors.New("json file contains duplicate key: " + data.prefix + name + " " + cr.currentPointInfo())
			}

			doReturn, continue_, err := cr.setJsonTempDataFromName(&data, name)
			if doReturn {
				if data.prefix == "" {
					return cr.readJsonTillEnd(r)
				}
				return cr.processEofError(err)
			} else if continue_ {
				continue
			}

			if data.foundInfo.keyName != "" && data.foundInfo.isMap {
				found = true
				foundInfo = data.foundInfo
			} else {
				found, foundInfo = cr.findFieldByJsonName(si, data.prefix+name)
			}

			data.parseState = psJsonValue

		case psJsonValue:
			valueResult := cr.readJsonValue(r)
			if valueResult.err != nil {
				return cr.processEofError(valueResult.err)
			} else if valueResult.divider == ',' {
				if notComma && valueResult.value == "" && !valueResult.isString {
					return cr.invalidCharacterError()
				}
				notComma = true
			} else if valueResult.divider == '}' && !data.isObject || valueResult.divider == ']' && data.isObject {
				return cr.invalidCharacterError()
			}
			if valueResult.isOpener {
				newData, doReturn, continue_, break_, err := cr.setJsonTempDataFromValue(data, valueResult.value, name)
				if doReturn {
					return cr.processEofError(err)
				} else if continue_ {
					newData.foundInfo = foundInfo
					newData.foundInfo.isSlice = valueResult.value == "["
					err = cr.parseJsonData(r, it, si, newData, sourceId)
					if err != nil {
						return err
					}
					divider = ' '
					data.parseState = psJsonDivider
					continue
				} else if break_ {
					doBreak = true
					break
				}
			}
			divider = valueResult.divider

			foundInfo = data.foundInfo
			if foundInfo.keyName != "" {
				found = true
				if foundInfo.isSlice || foundInfo.isMap {
					data.prefix = strings.TrimSuffix(data.prefix, ".")
				}
			}
			if found {
				vType := vtAny
				if valueResult.isString {
					vType = vtString
				} else {
					vType = cr.getJsonValueType(valueResult.value)
				}
				if foundInfo.isMap {
					err = cr.addJsonValue(foundInfo, it, data.prefix, valueResult.value, name, vType, sourceId)
				} else if !(vType == vtEmpty && !valueResult.isString && valueResult.divider != ' ') {
					err = cr.addJsonValue(foundInfo, it, data.prefix+name, valueResult.value, "", vType, sourceId)
				}
				if err != nil {
					return err
				}
			}

			if valueResult.divider == '}' || valueResult.divider == ']' {
				return nil
			}

			if !data.isObject && foundInfo.isSlice && valueResult.divider == ',' {
				data.parseState = psJsonValue
			} else {
				data.parseState = psJsonDivider
			}

		case psJsonDivider:
			isComma := false
			isEnd := false
			if divider == ' ' {
				isComma, isEnd, divider, err = cr.readJsonDivider(r)
			} else {
				isComma = divider == ','
				notComma = isComma
				isEnd = divider == '}' || divider == ']'
			}
			if err != nil {
				return cr.processEofError(err)
			}

			_, _, err := cr.setJsonTempDataFromName(&data, string(divider))
			if err != nil {
				return cr.processEofError(err)
			}

			if isComma {
				if data.isObject {
					data.parseState = psJsonName
				} else {
					data.parseState = psJsonValue
				}
				continue
			} else if isEnd {
				doBreak = true
				break
			}

			data.parseState = psJsonName
		}

		if doBreak {
			break
		}
	}

	return nil
}

func (cr *configReader) readJsonName(r *bufio.Reader) (string, error) {
	var buffer bytes.Buffer
	started := false
	finished := false
	readToColon := false
	expectComment := false
	escape := false

	rn := ' '
	var err error = nil
	for {
		if rn, _, err = r.ReadRune(); err == nil {
			cr.data.currentPos++
			if rn == '\n' || rn == '\r' {
				cr.checkNewLine(rn)
			}
			if !started && containsRune([]rune{'{', '}', '[', ']'}, rn) {
				return string(rn), nil
			} else if (!started || readToColon) && containsRune([]rune{' ', '\t', '\r', '\n'}, rn) {
				continue
			} else if (!started || finished) && rn == '/' || (expectComment && (rn == '/' || rn == '*')) {
				if expectComment {
					err = cr.readComment(r, rn == '/')
					if err != nil && err != io.EOF {
						return "", err
					}
					expectComment = false
					continue
				}
				expectComment = true
				continue
			} else if readToColon {
				if rn == ':' {
					break
				} else {
					return "", cr.invalidCharacterError()
				}
			} else if rn == '\\' && started && !finished && !escape {
				escape = true
				continue
			} else if rn == '"' {
				if escape {
					buffer.WriteRune(rn)
					escape = false
				} else if started {
					finished = true
					readToColon = true
				} else {
					started = true
				}
				continue
			} else if !started {
				return "", cr.invalidCharacterError()
			}

			if escape {
				buffer.WriteRune('\\')
				escape = false
			}
			buffer.WriteRune(rn)
		} else if err != io.EOF {
			return "", err
		} else {
			if started && !finished {
				return "", errors.New("unexpected end of file")
			}
			break
		}
	}

	name := strings.Trim(buffer.String(), " \t")
	if name == "" && (err == nil || err == io.EOF) {
		err = errors.New("wrong format: can't read name")
	} else {
		name = strings.ToLower(name)
	}

	return name, err
}

func (cr *configReader) readJsonValue(r *bufio.Reader) jsonReadValueResult {
	var buffer bytes.Buffer
	started := false
	finished := false
	isQuoted := false
	expectComment := false
	escape := false
	divider := ' '

	rn := ' '
	var err error = nil
	for {
		if rn, _, err = r.ReadRune(); err == nil {
			cr.data.currentPos++
			if rn == '\n' || rn == '\r' {
				cr.checkNewLine(rn)
			}
			if !isQuoted && rn == '/' || (expectComment && (rn == '/' || rn == '*')) {
				if expectComment {
					err = cr.readComment(r, rn == '/')
					if err != nil && err != io.EOF {
						return jsonReadValueResult{"", false, false, ' ', err}
					}
					expectComment = false
					continue
				}
				expectComment = true
				finished = true
				continue
			} else if !started {
				if containsRune([]rune{' ', '\t', '\r', '\n'}, rn) {
					continue
				} else if rn == '{' || rn == '[' {
					return jsonReadValueResult{string(rn), true, false, ' ', nil}
				} else if rn == ',' || rn == '}' || rn == ']' {
					if rn == '}' {
						return jsonReadValueResult{"", false, false, rn, cr.invalidCharacterError()}
					}
					return jsonReadValueResult{"", false, false, rn, nil}
				} else {
					started = true
					if rn == '"' {
						isQuoted = true
						continue
					}
				}
			} else if isQuoted && rn == '\\' && !escape {
				escape = true
				continue
			} else if rn == '"' {
				if escape && isQuoted {
					buffer.WriteRune(rn)
					escape = false
					continue
				} else if isQuoted {
					break
				} else {
					return jsonReadValueResult{"", false, false, ' ', cr.invalidCharacterError()}
				}
			} else if !isQuoted && containsRune([]rune{',', '}', ']', ' ', '\t'}, rn) {
				if rn == ',' || rn == '}' || rn == ']' {
					divider = rn
				}
				break
			}

			if escape {
				buffer.WriteRune('\\')
				escape = false
			}
			buffer.WriteRune(rn)
		} else if err != io.EOF {
			return jsonReadValueResult{"", false, false, ' ', err}
		} else {
			if started && !finished {
				return jsonReadValueResult{"", false, false, ' ', errors.New("unexpected end of file")}
			}
			break
		}
	}

	if isQuoted {
		name := strings.Trim(buffer.String(), " \t")
		return jsonReadValueResult{name, false, true, divider, nil}
	}

	name := strings.Trim(buffer.String(), " \t")
	name = strings.ToLower(name)
	if name == "true" || name == "false" || name == "null" {
		return jsonReadValueResult{name, false, false, divider, nil}
	}
	if _, errParse := strconv.ParseFloat(name, 64); errParse != nil {
		return jsonReadValueResult{name, false, false, divider, errors.New("wrong value format: " + name + " " + cr.currentPointInfo())}
	}

	return jsonReadValueResult{name, false, false, divider, nil}
}

func (cr *configReader) readJsonDivider(r *bufio.Reader) (bool, bool, rune, error) {
	isComma := false
	isEnd := false
	expectComment := false

	rn := ' '
	var err error = nil
	for {
		if rn, _, err = r.ReadRune(); err == nil {
			cr.data.currentPos++
			if rn == '/' || (expectComment && (rn == '/' || rn == '*')) {
				if expectComment {
					err = cr.readComment(r, rn == '/')
					if err != nil && err != io.EOF {
						return false, false, ' ', err
					}
					expectComment = false
					continue
				}
				expectComment = true
				continue
			} else if containsRune([]rune{' ', '\t', '\r', '\n'}, rn) {
				if rn == '\n' || rn == '\r' {
					cr.checkNewLine(rn)
				}
				continue
			} else if rn == ',' {
				isComma = true
				break
			} else if rn == '}' || rn == ']' {
				isEnd = true
				break
			} else {
				return false, false, ' ', cr.invalidCharacterError()
			}
		} else if err != io.EOF {
			return false, false, ' ', err
		} else {
			return false, false, ' ', errors.New("unexpected end of file")
		}
	}

	return isComma, isEnd, rn, nil
}

func (cr *configReader) readJsonTillEnd(r *bufio.Reader) error {
	rn := ' '
	var err error = nil
	for {
		if rn, _, err = r.ReadRune(); err == nil {
			cr.data.currentPos++
			if rn == '\n' || rn == '\r' {
				cr.checkNewLine(rn)
				continue
			} else if rn == ' ' || rn == '\t' {
				continue
			} else {
				return cr.invalidCharacterError()
			}
		} else if err != io.EOF {
			return err
		} else {
			return nil
		}
	}
}

func (cr *configReader) getJsonValueType(value string) valueType {
	if value == "" {
		return vtEmpty
	} else {
		value = strings.Trim(value, " \t")
		value = strings.ToLower(value)
	}
	if value == "true" || value == "false" {
		return vtBool
	} else if value == "null" {
		return vtNull
	} else if value[0] == '"' {
		return vtString
	}
	return vtNumber
}

func (cr *configReader) setJsonTempDataFromName(data *jsonTempData, name string) (bool, bool, error) {
	if data.prefix == "ROOT" {
		if name == "{" {
			data.prefix = ""
			data.isObject = true
			return false, true, nil
		} else {
			return true, false, errors.New("wrong format: root element must be an object")
		}
	} else if data.isObject && name == "}" {
		return true, false, nil
	} else if !data.isObject && name == "]" {
		return true, false, nil
	} else if data.isObject && name == "]" {
		return true, false, errors.New("wrong format: object can't be closed with ']'")
	} else if !data.isObject && name == "}" {
		return true, false, errors.New("wrong format: array can't be closed with '}'")
	}
	return false, false, nil
}

func (cr *configReader) setJsonTempDataFromValue(data jsonTempData, value, name string) (jsonTempData, bool, bool, bool, error) {
	newData := jsonTempData{}
	continue_ := false
	if value == "{" {
		newData.prefix = cr.getJsonPrefix(data.prefix, name)
		newData.isObject = true
		continue_ = true
	} else if value == "[" {
		newData.prefix = cr.getJsonPrefix(data.prefix, name)
		newData.isObject = false
		continue_ = true
	} else if value == "}" || value == "]" {
		return newData, true, continue_, true, nil
	}
	return newData, false, continue_, false, nil
}
