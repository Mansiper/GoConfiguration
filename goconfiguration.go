package configuration

import (
	"errors"
	"strings"
)

// Create a new config reader instance
// files - list of files to read configuration from
func NewConfigReader(files ...string) *configReader {
	config := &configReader{
		sources: []configSource{},
		options: ConfigOptions{
			RewriteValues: true,
			Parsers:       make(map[string]Parser),
		},
		data: configData{
			currentLine: 0,
			currentPos:  0,
			currentFile: "",
			initErrors:  []string{},
		},
	}

	for _, file := range files {
		config.AddFile(file)
	}

	return config
}

// Add configuration file
// file - relative or absolute path to the file
// supported file types: json, ini, env
func (cr *configReader) AddFile(file string) *configReader {
	fileType := cr.getFileType(file)
	if fileType == ftUnknown {
		cr.data.initErrors = append(cr.data.initErrors, unsupportedFileTypeError(file))
		return cr
	}

	source := configSource{
		value:    file,
		ft:       fileType,
		fromFile: true,
	}
	cr.sources = append(cr.sources, source)
	return cr
}

// Use environment variables as a configuration source
func (cr *configReader) AddEnvironment() *configReader {
	source := configSource{
		value:    "",
		ft:       ftEnvironment,
		fromFile: false,
	}
	cr.sources = append(cr.sources, source)
	return cr
}

// Add configuration as a string
// values - configuration string
// formatType - format of the configuration string
func (cr *configReader) AddString(values string, formatType formatType, name string) *configReader {
	source := configSource{
		name:     name,
		value:    values,
		ft:       formatType,
		fromFile: false,
	}
	cr.sources = append(cr.sources, source)
	return cr
}

// Set configuration reader options
// options - configuration reader options
func (cr *configReader) WithOptions(options ConfigOptions) *configReader {
	if options.Parsers == nil {
		options.Parsers = make(map[string]Parser)
	}
	cr.options = options

	return cr
}

// Set whether to rewrite values (not for slice values)
// rewrite - rewrite values or not
func (cr *configReader) RewriteValues(rewrite bool) *configReader {
	cr.options.RewriteValues = rewrite
	return cr
}

// Add custom parser for user types
// fieldName - field name
// parser - custom parser for user type
func (cr *configReader) WithParser(envName string, parser Parser) *configReader {
	if cr.options.Parsers == nil {
		cr.options.Parsers = make(map[string]Parser)
	}
	cr.options.Parsers[envName] = parser
	return cr
}

// Ensure that there are no errors during configuration reading
// panic if there are errors
func (cr *configReader) EnsureHasNoErrors() *configReader {
	if len(cr.data.initErrors) > 0 {
		panic(strings.Join(cr.data.initErrors, "\n"))
	}
	return cr
}

// Get errors that occurred during preparation of the configuration reader
func (cr *configReader) GetErrors() []error {
	if len(cr.data.initErrors) > 0 {
		errs := make([]error, len(cr.data.initErrors))
		for i, e := range cr.data.initErrors {
			errs[i] = errors.New(e)
		}
		return errs
	}
	return nil
}

// Read configuration into the user config struct
// userConfig - pointer to the user config struct
func (cr *configReader) ReadConfig(userConfig interface{}) error {
	si, err := cr.getStructInfo(userConfig, "", "")
	if err != nil {
		return err
	}

	it := intermediateTree{}

	for i, source := range cr.sources {
		var err error = nil
		if source.ft == ftEnvironment {
			cr.readEnvironment(it, si, i)
		} else if source.fromFile {
			err = cr.readConfigFile(source, it, si, i)
		} else {
			err = cr.readConfigString(source, it, si, i)
		}
		if err != nil {
			return err
		}
	}

	err = cr.setValues(it, si)
	if err != nil {
		return err
	}

	return nil
}
