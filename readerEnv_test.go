package configuration

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseEnvData_success_big(t *testing.T) {
	// Arrange
	allCases := `# comment
env_1=value_1

env_2	=	"val\"ue_\\2"
env_3 = value_3 # comment
unknown = value

env_4=value_4\
line_2\
line_3
  env_5="value_5 # comment"  
env_6=value_6\
#comment
sub.env_1='value "1" #'

arr1[] = 1
arr1[] = "2"
arr2 = "1,2,3"
arr2 = 4,5
arr2[] = 6

map[key1] = '1'
map[key2] = 2
`
	si := []structInfo{
		{keyName: "env_1"},
		{keyName: "env_2"},
		{keyName: "env_3"},
		{keyName: "env_4"},
		{keyName: "env_5"},
		{keyName: "env_6"},
		{keyName: "sub.env_1"},
		{keyName: "arr1", isSlice: true},
		{keyName: "arr2", isSlice: true, separator: ","},
		{keyName: "map", isMap: true},
	}
	r := bufio.NewReader(strings.NewReader(allCases))
	it := make(intermediateTree)
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}

	// Act
	err := cr.parseEnvData(r, it, si, 0)

	// Assert
	require.Nil(t, err)
	require.Equal(t, "value_1", it["env_1"][0].value)
	require.Equal(t, "val\"ue_\\\\2", it["env_2"][0].value)
	require.Equal(t, "value_3", it["env_3"][0].value)
	require.Equal(t, "value_4\nline_2\nline_3", it["env_4"][0].value)
	require.Equal(t, "value_5 # comment", it["env_5"][0].value)
	require.Equal(t, "value_6\n", it["env_6"][0].value)
	require.Equal(t, "value \"1\" #", it["sub.env_1"][0].value)
	require.Equal(t, "1", it["arr1"][0].value.([]string)[0])
	require.Equal(t, "2", it["arr1"][0].value.([]string)[1])
	require.Equal(t, "1", it["arr2"][0].value.([]string)[0])
	require.Equal(t, "2", it["arr2"][0].value.([]string)[1])
	require.Equal(t, "3", it["arr2"][0].value.([]string)[2])
	require.Equal(t, "4", it["arr2"][0].value.([]string)[3])
	require.Equal(t, "5", it["arr2"][0].value.([]string)[4])
	require.Equal(t, "6", it["arr2"][0].value.([]string)[5])
	require.Equal(t, "1", it["map"][0].value.(map[string]string)["key1"])
	require.Equal(t, "2", it["map"][0].value.(map[string]string)["key2"])
}

type testCaseEnvSuccess struct {
	data    string
	name    string
	value   string
	isSlice bool
	isMap   bool
	sep     string
}

func Test_parseEnvData_success_cases(t *testing.T) {
	// Arrange
	cases := []testCaseEnvSuccess{
		{"env_1=value", "env_1", "value", false, false, ""},
		{"env_1='value'", "env_1", "value", false, false, ""},
		{"env_1='value # x'", "env_1", "value # x", false, false, ""},
		{"env_1='value' # x'", "env_1", "value", false, false, ""},
		{"env_1=\"value\"", "env_1", "value", false, false, ""},
		{"env_1=value\n", "env_1", "value", false, false, ""},
		{" env_1 = value ", "env_1", "value", false, false, ""},
		{"env_1='v\\'alue'", "env_1", "v'alue", false, false, ""},
		{"env_1=\"v\\\"alue\"", "env_1", "v\"alue", false, false, ""},
		{"\tenv_1\t=\tvalue\t", "env_1", "value", false, false, ""},
		{"\r\nenv_1=value\r\n", "env_1", "value", false, false, ""},
		{"env_1=value\\\r\n#l2", "env_1", "value\n", false, false, ""},
		{"env_1=value\\\nl2\\\nl3", "env_1", "value\nl2\nl3", false, false, ""},
		{"env_1=\\\nvalue", "env_1", "\nvalue", false, false, ""},
		{"env_1=''", "env_1", "", false, false, ""},
		{"env_1=# comment", "env_1", "", false, false, ""},
		{"env_1=", "env_1", "", false, false, ""},
		{"env_1=\n", "env_1", "", false, false, ""},
		{"env_1=#value", "env_1", "", false, false, ""},
		{"env_1[]=value", "env_1", "value", true, false, ","},
		{"env_1[]=", "env_1", "", true, false, ","},
		{"env_1[]=1;2;3", "env_1", "1;2;3", true, false, ";"},
		{"env_1[]='1;2;3'\n", "env_1", "1;2;3", true, false, ";"},
		{"env_1='1;2;3'\n", "env_1", "1,2,3", true, false, ";"},
		{"env_1='1;2;3'\nenv_1=4", "env_1", "1,2,3,4", true, false, ";"},
		{"env_1[key]=a\n", "env_1", "a", false, true, ","},
		{"env_1[key]='a b'", "env_1", "a b", false, true, ","},
		{"env_]1[=value", "env_]1[", "value", false, false, ""},
		{"env_1[=value", "env_1[", "value", false, false, ""},
		{"[env_1]=value", "[env_1]", "value", false, false, ""},
		{"[e]nv_1=value", "[e]nv_1", "value", false, false, ""},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data)
		test_parseEnvData_success_cases(t, c)
	}
}

func test_parseEnvData_success_cases(t *testing.T, testCase testCaseEnvSuccess) {
	// Arrange
	r := bufio.NewReader(strings.NewReader(testCase.data))
	it := make(intermediateTree)
	si := []structInfo{{
		keyName:   testCase.name,
		isSlice:   testCase.isSlice,
		isMap:     testCase.isMap,
		separator: testCase.sep,
	}}
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}
	cr.data.currentLine = 1

	// Act
	err := cr.parseEnvData(r, it, si, 0)

	// Assert
	assert.Nil(t, err)
	if testCase.isSlice {
		require.Equal(t, testCase.value, strings.Join(it[testCase.name][0].value.([]string), ","))
	} else if testCase.isMap {
		require.Equal(t, testCase.value, it[testCase.name][0].value.(map[string]string)["key"])
	} else {
		require.Equal(t, testCase.value, it[testCase.name][0].value)
	}
}

type testCaseEnvWrongName struct {
	data  string
	check string
}

func Test_parseEnvData_wrongName_cases(t *testing.T) {
	// Arrange
	cases := []testCaseEnvWrongName{
		{"env 1=value_1", "(1:5)"},
		{"env_\n1=value_1", "(1:5)"},
		{"env_1", "end of file"},
		{"env_1\n", "(1:6)"},
		{"env#1=value", "(1:4)"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data, c.check)
		test_parseEnvData_wrongName(t, c)
	}
}

func test_parseEnvData_wrongName(t *testing.T, testCase testCaseEnvWrongName) {
	// Arrange
	r := bufio.NewReader(strings.NewReader(testCase.data))
	it := make(intermediateTree)
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}
	cr.data.currentLine = 1

	// Act
	err := cr.parseEnvData(r, it, []structInfo{}, 0)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), testCase.check)
}

type testCaseEnvWrongValue struct {
	data  string
	check string
}

func Test_parseEnvData_wrongValue_cases(t *testing.T) {
	// Arrange
	cases := []testCaseEnvWrongValue{
		{"env_1=value 1", "(1:13)"},
		{"env_1='value", "(1:12)"},
		{"env_1=\"value", "(1:12)"},
		{"env_1=\"value'", "(1:13)"},
		{"env_1='value\"", "(1:13)"},
		{"env_1=value\\ ", "(1:13)"},
		{"env_1='value\\\nl2", "(1:14)"},
		{"env_1=value\\\nl2\"", "(2:3)"},
		{"env_1=value\\\nl 2\n", "(2:3)"},
		{"env_1='val'ue'", "(1:12)"},
		{"env_1=\"val\"ue\"", "(1:12)"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data, c.check)
		test_parseEnvData_wrongValue(t, c)
	}
}

func test_parseEnvData_wrongValue(t *testing.T, testCase testCaseEnvWrongValue) {
	// Arrange
	r := bufio.NewReader(strings.NewReader(testCase.data))
	it := make(intermediateTree)
	si := []structInfo{{keyName: "env_1"}}
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}
	cr.data.currentLine = 1

	// Act
	err := cr.parseEnvData(r, it, si, 0)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), testCase.check)
}
