package configuration

import (
	"bufio"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseIniData_success_big(t *testing.T) {
	// Arrange
	allCases := `# comment
; comment
env 1=value 1
unknown = value

[ rOOt ]
env_2	=	"va\"lue_\\2"
 [sub] 
env_3 = value_3 # comment

env_4=value_4\
line_2 \
line 3
  env_5="value_5 # comment"  
[.sub2]
env_6=value_6\
#comment
[sub.sub2]
env_1='value "1" #'

[root]
arr1[] = 1
arr1[] = "2"
arr2 = "1,2,3"
arr2 = 4,5
arr2[] = 6

map[key1] = '1'
map[key2] = 2
`
	si := []structInfo{
		{keyName: "env 1"},
		{keyName: "env_2"},
		{keyName: "sub.env_3"},
		{keyName: "sub.env_4"},
		{keyName: "sub.env_5"},
		{keyName: "sub.sub2.env_6"},
		{keyName: "sub.sub2.env_1"},
		{keyName: "arr1", isSlice: true},
		{keyName: "arr2", isSlice: true, separator: ","},
		{keyName: "map", isMap: true},
	}
	r := bufio.NewReader(strings.NewReader(allCases))
	it := make(intermediateTree)
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}

	// Act
	err := cr.parseIniData(r, it, si, 0)

	// Assert
	require.Nil(t, err)
	require.Equal(t, "value 1", it["env 1"][0].value)
	require.Equal(t, "va\"lue_\\\\2", it["env_2"][0].value)
	require.Equal(t, "value_3", it["sub.env_3"][0].value)
	require.Equal(t, "value_4\nline_2 \nline 3", it["sub.env_4"][0].value)
	require.Equal(t, "value_5 # comment", it["sub.env_5"][0].value)
	require.Equal(t, "value_6\n", it["sub.sub2.env_6"][0].value)
	require.Equal(t, "value \"1\" #", it["sub.sub2.env_1"][0].value)
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

type testCaseIniSuccess struct {
	data    string
	name    string
	value   string
	isSlice bool
	isMap   bool
	sep     string
}

func Test_parseIniData_success_cases(t *testing.T) {
	// Arrange
	cases := []testCaseIniSuccess{
		{data: "env 1=value", name: "env 1", value: "value"},
		{data: "env_1='value'", name: "env_1", value: "value"},
		{data: "env_1='value # x'", name: "env_1", value: "value # x"},
		{data: "env_1='value' # x'", name: "env_1", value: "value"},
		{data: "env_1=\"value\"", name: "env_1", value: "value"},
		{data: "env_1=\"va\\\"lue\"", name: "env_1", value: "va\"lue"},
		{data: "env_1='va\\'lue'", name: "env_1", value: "va'lue"},
		{data: "env_1=value\n", name: "env_1", value: "value"},
		{data: " env_1 = value ", name: "env_1", value: "value"},
		{data: "\tenv_1\t=\tvalue\t", name: "env_1", value: "value"},
		{data: "\r\nenv_1=value\r\n", name: "env_1", value: "value"},
		{data: "env_1=value\\\n\r#l2", name: "env_1", value: "value\n"},
		{data: "env_1=value\\\r\n\\#l2\n", name: "env_1", value: "value\n\n"},
		{data: "env_1=value\\\nl2\\\nl3", name: "env_1", value: "value\nl2\nl3"},
		{data: "env_1=value\\\n\\#l2\ne=v", name: "env_1", value: "value\n\ne=v"},
		{data: "env_1=''", name: "env_1", value: ""},
		{data: "env_1=# comment", name: "env_1", value: ""},
		{data: "env_1=", name: "env_1", value: ""},
		{data: "env_1=\n", name: "env_1", value: ""},
		{data: "env_1=#value", name: "env_1", value: ""},
		{data: "[ROOT]\nenv_1=value", name: "env_1", value: "value"},
		{data: "[ root\t]\nenv_1=value", name: "env_1", value: "value"},
		{data: "[sub]\nenv_1=value", name: "sub.env_1", value: "value"},
		{data: "[sub.sub 2]\nenv_1=value", name: "sub.sub 2.env_1", value: "value"},
		{data: "[sub]\n[.sub2]\nenv_1=value", name: "sub.sub2.env_1", value: "value"},
		{data: "[sub]\n[.sub2.sub3]\nenv_1=value", name: "sub.sub2.sub3.env_1", value: "value"},
		{data: "[sub]# comment\nenv_1=value", name: "sub.env_1", value: "value"},
		{data: "; comment\n[sub]\nenv_1=value", name: "sub.env_1", value: "value"},
		{data: "[sub;#]\nenv_1=value", name: "sub;#.env_1", value: "value"},
		{data: "[\"sub\"]\nenv_1=value", name: "\"sub\".env_1", value: "value"},
		{"env_1[]=value", "env_1", "value", true, false, ","},
		{"env_1[]=", "env_1", "", true, false, ","},
		{"env_1[]=\"1;2;3\"", "env_1", "1;2;3", true, false, ";"},
		{"env_1[]='1;2;3'\n", "env_1", "1;2;3", true, false, ";"},
		{"env_1='1;2;3'\n", "env_1", "1,2,3", true, false, ";"},
		{"env_1='1;2;3'\nenv_1=4", "env_1", "1,2,3,4", true, false, ";"},
		{"env_1[key]=a\n", "env_1", "a", false, true, ","},
		{"env_1[key]='a b'", "env_1", "a b", false, true, ","},
		{"env_]1[=value", "env_]1[", "value", false, false, ""},
		{"env_1[=value", "env_1[", "value", false, false, ""},
	}

	// Act & Assert
	for i, c := range cases {
		t.Log("Test case "+strconv.Itoa(i+1)+":", c.data)
		test_parseIniData_success_cases(t, c)
	}
}

func test_parseIniData_success_cases(t *testing.T, testCase testCaseIniSuccess) {
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
	err := cr.parseIniData(r, it, si, 0)

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

type testCaseIniWrongName struct {
	data  string
	check string
}

func Test_parseIniData_wrongName_cases(t *testing.T) {
	// Arrange
	cases := []testCaseIniWrongName{
		{"env_\n1=value_1", "(1:5)"},
		{"env_1", "end of file"},
		{"env_1\n", "(1:6)"},
		{"env_1 ", "end of file"},
		{"[]\nenv_1=value_1", "(1:3)"},
		{"[ROOT] s", "(1:8)"},
		{"[sub.]\nenv=", "(2:0)"},
		{"[[sub]]", "(1:2)"},
		{"[env_1]=value", "(1:8)"},
		{"[e]nv_1=value", "(1:4)"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data, c.check)
		test_parseIniData_wrongName(t, c)
	}
}

func test_parseIniData_wrongName(t *testing.T, testCase testCaseIniWrongName) {
	// Arrange
	r := bufio.NewReader(strings.NewReader(testCase.data))
	it := make(intermediateTree)
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}
	cr.data.currentLine = 1

	// Act
	err := cr.parseIniData(r, it, []structInfo{}, 0)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), testCase.check)
}

type testCaseIniWrongValue struct {
	data  string
	check string
}

func Test_parseIniData_wrongValue_cases(t *testing.T) {
	// Arrange
	cases := []testCaseIniWrongValue{
		{"env_1='value", "(1:12)"},
		{"env_1=\"value", "(1:12)"},
		{"env_1=\"value'", "(1:13)"},
		{"env_1='value\"", "(1:13)"},
		{"env_1=value\\ ", "(1:13)"},
		{"env_1='value\\\nl2", "(1:14)"},
		{"env_1=value\\\nl2\"", "(2:3)"},
		{"env_1='val'ue'", "(1:12)"},
		{"env_1=\"val\"ue\"", "(1:12)"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data, c.check)
		test_parseIniData_wrongValue(t, c)
	}
}

func test_parseIniData_wrongValue(t *testing.T, testCase testCaseIniWrongValue) {
	// Arrange
	r := bufio.NewReader(strings.NewReader(testCase.data))
	it := make(intermediateTree)
	si := []structInfo{{keyName: "env_1"}}
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}
	cr.data.currentLine = 1

	// Act
	err := cr.parseIniData(r, it, si, 0)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), testCase.check)
}
