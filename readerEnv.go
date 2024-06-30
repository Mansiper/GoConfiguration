package configuration

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

func (cr *configReader) parseEnvData(r *bufio.Reader, it intermediateTree, si []structInfo, sourceId int) error {
	for {
		name, err := cr.readEnvName(r)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
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

		found, foundInfo, continue_, err := cr.findFieldByName(r, si, name, true)
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

		value, err := cr.readEnvValue(r)
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

func (cr *configReader) readEnvName(r *bufio.Reader) (string, error) {
	var buffer bytes.Buffer
	started := false
	finished := false

	rn := ' '
	var err error = nil
	for {
		if rn, _, err = r.ReadRune(); err == nil {
			cr.data.currentPos++
			if rn == '#' {
				if !started {
					err = cr.readComment(r, true)
					if err != nil {
						return "", err
					}
					continue
				} else {
					return "", cr.invalidCharacterError()
				}
			} else if rn == '\r' || rn == '\n' {
				if !started {
					cr.checkNewLine(rn)
					continue
				} else if started {
					return "", cr.invalidCharacterError()
				}
			} else if rn == ' ' || rn == '\t' {
				if !started {
					continue
				} else {
					finished = true
					continue
				}
			} else if rn != ' ' && rn != '\t' && rn != '=' && finished {
				return "", cr.invalidCharacterError()
			} else if rn == '=' {
				break
			}

			buffer.WriteRune(rn)
			started = true
		} else if err != io.EOF {
			return "", err
		} else {
			if started && !finished {
				return "", errors.New("unexpected end of file")
			}
			break
		}
	}

	name := buffer.String()
	if started && name == "" {
		return "", errors.New("wrong format: can't read name")
	}
	name = strings.ToLower(name)

	return name, err
}

func (cr *configReader) readEnvValue(r *bufio.Reader) (string, error) {
	var buffer bytes.Buffer
	started := false
	finished := false
	inQuote := false
	quote := ' '
	backslash := false
	readToBol := false
	escape := false

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
					buffer.WriteRune('\n')
				}
			}

			if backslash {
				if rn == '\r' || rn == '\n' {
					cr.checkNewLine(rn)
					backslash = false
					readToBol = true
					continue
				} else {
					return "", cr.invalidCharacterError()
				}
			} else if rn == ' ' || rn == '\t' {
				if !started || finished {
					continue
				} else if started && !finished && !inQuote {
					finished = true
					continue
				}
			} else if rn == '\\' {
				if started && !inQuote || !started {
					started = true
					backslash = true
					continue
				} else if started && inQuote && !finished && !escape {
					escape = true
					continue
				}
			} else if rn == '#' && !inQuote {
				if started {
					finished = true
				}
				err = cr.readComment(r, true)
				if err != nil && err != io.EOF {
					return "", err
				}
				break
			} else if rn == '"' || rn == '\'' {
				if !started && !inQuote {
					inQuote = true
					quote = rn
					continue
				} else if !started && inQuote {
					inQuote = false
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
			} else if !containsRune([]rune{' ', '\t', '#', '\r', '\n'}, rn) && finished {
				return "", cr.invalidCharacterError()
			} else if rn == '\r' || rn == '\n' {
				if !started {
					return "", nil
				} else if started && !finished && inQuote {
					return "", cr.invalidCharacterError()
				} else if started && !finished && !inQuote || finished {
					cr.checkNewLine(rn)
					if rn == '\n' {
						break
					} else {
						continue
					}
				}
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
