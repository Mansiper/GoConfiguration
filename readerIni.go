package configuration

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

func (cr *configReader) parseIniData(r *bufio.Reader, it intermediateTree, si []structInfo, sourceId int) error {
	prefix := ""
	for {
		str, isName, err := cr.readIniNameOrSection(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if !isName {
			if strings.ToUpper(str) == "ROOT" {
				prefix = ""
				continue
			}

			if str == "" {
				return errors.New("wrong format: specify section or subsection name " + cr.currentPointInfo())
			}
			if strings.HasSuffix(str, ".") {
				return errors.New("wrong format: section or subsection must not end with '.' " + cr.currentPointInfo())
			}

			if strings.HasPrefix(str, ".") {
				if prefix == "" {
					return errors.New("wrong format: subsection for root must not start with '.' " + cr.currentPointInfo())
				}
				prefix = prefix + str[1:] + "."
			} else {
				prefix = str + "."
			}
			continue
		}

		name := prefix + str
		isSlice := strings.HasSuffix(name, "[]")
		if isSlice {
			name = name[:len(name)-2]
		}
		openIndex := strings.Index(name, "[")
		isMap := len(name) >= 4 && !isSlice && openIndex > 0 && openIndex < len(name)-2 && strings.HasSuffix(name, "]")
		key := ""
		if isMap {
			key = name[openIndex+1 : len(name)-1]
			name = name[:openIndex]
		}

		found, foundInfo, continue_, err := cr.findFieldByName(r, si, name, false)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		} else if !found && continue_ {
			continue
		} else if !found && !continue_ {
			break
		}

		value, err := cr.readIniValue(r)
		if err != nil && err != io.EOF {
			return err
		}

		addValue(foundInfo, it, name, value, key, isSlice, sourceId)
		if err == io.EOF {
			break
		}
	}

	return nil
}

func (cr *configReader) readIniNameOrSection(r *bufio.Reader) (string, bool, error) {
	var buffer bytes.Buffer
	started := false
	finished := false
	readToEol := false

	rn := ' '
	isName := true
	var err error = nil
	for {
		if rn, _, err = r.ReadRune(); err == nil {
			cr.data.currentPos++
			if readToEol {
				if !containsRune([]rune{'\r', '\n', ' ', '\t', ';', '#'}, rn) {
					return "", false, cr.invalidCharacterError()
				} else {
					cr.checkNewLine(rn)
					if rn == '\n' {
						break
					} else if rn == '\r' {
						continue
					} else if rn == ';' || rn == '#' {
						err = cr.readComment(r, true)
						if err != nil && err != io.EOF {
							return "", false, err
						}
						break
					}
				}
			} else if (rn == '#' || rn == ';') && !started {
				err = cr.readComment(r, true)
				if err != nil && err != io.EOF {
					return "", false, err
				}
				continue
			} else if rn == '[' {
				if !isName {
					return "", false, cr.invalidCharacterError()
				} else if !started {
					isName = false
					continue
				}
			} else if rn == ']' && started && !isName {
				readToEol = true
				continue
			} else if rn == '\r' || rn == '\n' {
				if !started {
					cr.checkNewLine(rn)
					continue
				} else if started {
					return "", false, cr.invalidCharacterError()
				}
			} else if (rn == ' ' || rn == '\t') && !started {
				continue
			} else if rn != ' ' && rn != '\t' && rn != '=' && finished {
				return "", false, cr.invalidCharacterError()
			} else if rn == '=' {
				if isName {
					break
				} else {
					return "", false, cr.invalidCharacterError()
				}
			}
			buffer.WriteRune(rn)
			started = true
		} else if err != io.EOF {
			return "", false, err
		} else {
			if started && !finished {
				return "", false, errors.New("unexpected end of file")
			}
			break
		}
	}

	str := buffer.String()
	if str == "" && started {
		return "", false, errors.New("wrong format: can't read name")
	}
	str = strings.Trim(str, " \t")
	str = strings.ToLower(str)

	return str, isName, err
}

func (cr *configReader) readIniValue(r *bufio.Reader) (string, error) {
	var buffer bytes.Buffer
	started := false
	finished := false
	inQuote := false
	quote := ' '
	backslash := false
	readToBol := false
	addNewLine := false
	newLineCntr := 0
	escape := false

	rn := ' '
	var err error = nil
	for {
		if rn, _, err = r.ReadRune(); err == nil {
			cr.data.currentPos++
			if rn == '\n' {
				newLineCntr += 1
				if newLineCntr > 1 {
					break
				}
			} else if rn != '\r' {
				newLineCntr = 0
			}
			if readToBol {
				if rn == '\r' || rn == '\n' {
					cr.checkNewLine(rn)
					if rn == '\n' && addNewLine {
						buffer.WriteRune('\n')
						addNewLine = false
					}
					continue
				} else {
					readToBol = false
				}
			}

			if (rn == ' ' || rn == '\t') && (!started || finished) {
				continue
			} else if rn == '\\' {
				if started && !inQuote {
					backslash = true
					continue
				} else if started && inQuote && !finished && !escape {
					escape = true
					continue
				}
			} else if backslash {
				if rn == '\r' || rn == '\n' {
					cr.checkNewLine(rn)
					if rn == '\n' {
						buffer.WriteRune('\n')
					} else {
						addNewLine = true
					}
					backslash = false
					readToBol = true
					continue
				} else if rn == ';' || rn == '#' {
					err = cr.readComment(r, true)
					if err != nil && err != io.EOF {
						return "", err
					}
					buffer.WriteRune('\n')
					backslash = false
					continue
				} else {
					return "", cr.invalidCharacterError()
				}
			} else if rn == '#' || rn == ';' {
				if !inQuote || finished {
					finished = true
					err = cr.readComment(r, true)
					if err != nil && err != io.EOF {
						return "", err
					}
					break
				}
			} else if rn == '\r' || rn == '\n' {
				if inQuote {
					return "", cr.invalidCharacterError()
				} else {
					cr.checkNewLine(rn)
					if rn == '\n' {
						break
					} else {
						continue
					}
				}
			} else if rn == '"' || rn == '\'' {
				if !started && !inQuote {
					inQuote = true
					quote = rn
					started = true
					continue
				} else if started && inQuote && quote == rn {
					if escape {
						buffer.WriteRune(rn)
						escape = false
						continue
					} else {
						inQuote = false
						finished = true
						continue
					}
				} else if started && !inQuote {
					return "", cr.invalidCharacterError()
				}
			} else if !containsRune([]rune{' ', '\t', '#', ';', '\r', '\n'}, rn) && finished {
				return "", cr.invalidCharacterError()
			}

			if escape {
				buffer.WriteRune('\\')
				escape = false
			}
			buffer.WriteRune(rn)
			started = true
		} else if err != io.EOF {
			return "", err
		} else {
			if inQuote {
				return "", cr.invalidCharacterError()
			}
			break
		}
	}
	return buffer.String(), err
}
