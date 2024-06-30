package configuration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewConfigReader_success_cases(t *testing.T) {
	// Arrange
	cases := [][]string{
		{},
		{"config.env"},
		{"config.env", "config.ini", "config.json"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", strings.Join(c, ","))
		test_NewConfigReader_success(t, c)
	}
}

func test_NewConfigReader_success(t *testing.T, files []string) {
	// Arrange
	result := NewConfigReader(files...)
	errs := result.GetErrors()

	// Assert
	require.NotNil(t, result)
	require.Empty(t, errs)
}

type testCaseAddFile struct {
	path          string
	expectedType  formatType
	expectedError string
}

func Test_AddFile_success_cases(t *testing.T) {
	// Arrange
	cases := []testCaseAddFile{
		{"config.env", FtEnv, ""},
		{"config.json", FtJson, ""},
		// {"config.yaml", FtYaml, ""},
		{"config.ini", FtIni, ""},
		{"config.unknown", ftUnknown, unsupportedFileTypeError("config.unknown")},
		{"C:\\config.ini", FtIni, ""},
		{"/config.env", FtEnv, ""},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.path)
		test_AddFile_success(t, c)
	}
}

func test_AddFile_success(t *testing.T, testCase testCaseAddFile) {
	// Arrange
	cr := NewConfigReader()

	// Act
	cr.AddFile(testCase.path)

	// Assert
	if testCase.expectedError == "" {
		require.Empty(t, cr.GetErrors())
		require.Len(t, cr.sources, 1)
		require.Equal(t, testCase.expectedType, cr.sources[0].ft)
	} else {
		require.Equal(t, []string{testCase.expectedError}, cr.data.initErrors)
		require.Len(t, cr.sources, 0)
	}
}

func Test_AddEnvironment_success(t *testing.T) {
	// Arrange
	cr := NewConfigReader()

	// Act
	cr.AddEnvironment()

	// Assert
	require.Empty(t, cr.GetErrors())
	require.Len(t, cr.sources, 1)
	require.Equal(t, ftEnvironment, cr.sources[0].ft)
	require.False(t, cr.sources[0].fromFile)
	require.Empty(t, cr.sources[0].value)
}

func Test_AddString_success_cases(t *testing.T) {
	// Arrange
	cases := []formatType{
		FtEnv,
		FtJson,
		// FtYaml,
		FtIni,
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c)
		test_AddString_success(t, c)
	}
}

func test_AddString_success(t *testing.T, ft formatType) {
	// Arrange
	const data = "conig data"
	cr := NewConfigReader()

	// Act
	cr.AddString(data, ft, "test")

	// Assert
	require.Empty(t, cr.GetErrors())
	require.Len(t, cr.sources, 1)
	require.Equal(t, ft, cr.sources[0].ft)
	require.False(t, cr.sources[0].fromFile)
	require.Equal(t, cr.sources[0].value, data)
	require.Equal(t, cr.sources[0].name, "test")
}

func Test_EnsureHasNoErrors_success(t *testing.T) {
	// Arrange
	cr := NewConfigReader()

	// Act
	cr.EnsureHasNoErrors()

	// Assert
	require.Empty(t, cr.GetErrors())
}

func Test_EnsureHasNoErrors_panic(t *testing.T) {
	// Arrange
	cr := NewConfigReader()
	cr.data.initErrors = append(cr.data.initErrors, "error1", "error2")

	// Act & Assert
	require.Panicsf(t, func() {
		cr.EnsureHasNoErrors()
	}, "error1\nerror2")
}

func Test_GetErrors_success_cases(t *testing.T) {
	// Arrange
	cases := [][]string{
		{},
		{"error1"},
		{"error1", "error2"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", strings.Join(c, ","))
		test_GetErrors_success(t, c)
	}
}

func test_GetErrors_success(t *testing.T, errs []string) {
	// Arrange
	cr := NewConfigReader()
	cr.data.initErrors = errs

	// Act
	result := cr.GetErrors()

	// Assert
	require.Len(t, result, len(errs))
	for i, e := range errs {
		require.Equal(t, e, result[i].Error())
	}
}
