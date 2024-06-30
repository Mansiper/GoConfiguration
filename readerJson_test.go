package configuration

import (
	"bufio"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseJsonData_success_big(t *testing.T) {
	// Arrange
	allCases := `{
	//comment
	"env 1": "value 1","env\"_\\2": "value\"_\\2",
	/*comm
	ent*/
	"//env_3/**/": "//value_3/**/",
	//"e": "v",
	"env_4"//: "value_4",
	: "value_4",
	"env_5": //"value_55",
	"value_5",
	"env_6": "value_6"//,
	,
	"env_7": true,
	"env_8": false,
	"env_9": 3.14,
	"env_10": -123,
	"env_11": null,
	"arr_1": [1,2,3],
	"arr_2": [true,false,],
	"map_1": {"key1": 1,"key2": 2,},
	"sub_1": {
		"env_1": "value_1",
		"env_2": [1,2]
	},
	"sub_2": {
		"sub_2_1": {"env_1": "value_1"},
		/*comme
				nt*/
	},
}`
	si := []structInfo{
		{keyName: "env 1"},
		{keyName: "env\"_\\\\2"},
		{keyName: "//env_3/**/"},
		{keyName: "env_4"},
		{keyName: "env_5"},
		{keyName: "env_6"},
		{keyName: "env_7"},
		{keyName: "env_8"},
		{keyName: "env_9"},
		{keyName: "env_10"},
		{keyName: "env_11"},
		{keyName: "arr_1", isSlice: true},
		{keyName: "arr_2", isSlice: true},
		{keyName: "map_1", isMap: true},
		{keyName: "sub_1.env_1"},
		{keyName: "sub_1.env_2", isSlice: true},
		{keyName: "sub_2.sub_2_1.env_1"},
		{keyName: "sub_2.arr_1", isSlice: true},
	}
	r := bufio.NewReader(strings.NewReader(allCases))
	it := make(intermediateTree)
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}

	// Act
	err := cr.parseJsonData(r, it, si, defaultJsonData, 0)

	require.Nil(t, err)
	require.Equal(t, "value 1", it["env 1"][0].value)
	require.Equal(t, "value\"_\\\\2", it["env\"_\\\\2"][0].value)
	require.Equal(t, "//value_3/**/", it["//env_3/**/"][0].value)
	require.Equal(t, "value_4", it["env_4"][0].value)
	require.Equal(t, "value_5", it["env_5"][0].value)
	require.Equal(t, "value_6", it["env_6"][0].value)
	require.Equal(t, "true", it["env_7"][0].value)
	require.Equal(t, "false", it["env_8"][0].value)
	require.Equal(t, "3.14", it["env_9"][0].value)
	require.Equal(t, "-123", it["env_10"][0].value)
	require.Equal(t, vtNull, it["env_11"][0].valueType)
	require.Equal(t, "1", it["arr_1"][0].value.([]string)[0])
	require.Equal(t, "2", it["arr_1"][0].value.([]string)[1])
	require.Equal(t, "3", it["arr_1"][0].value.([]string)[2])
	require.Equal(t, "true", it["arr_2"][0].value.([]string)[0])
	require.Equal(t, "false", it["arr_2"][0].value.([]string)[1])
	require.Equal(t, "1", it["map_1"][0].value.(map[string]string)["key1"])
	require.Equal(t, "2", it["map_1"][0].value.(map[string]string)["key2"])
	require.Equal(t, "value_1", it["sub_1.env_1"][0].value)
	require.Equal(t, "1", it["sub_1.env_2"][0].value.([]string)[0])
	require.Equal(t, "2", it["sub_1.env_2"][0].value.([]string)[1])
	require.Equal(t, "value_1", it["sub_2.sub_2_1.env_1"][0].value)
}

type testCaseJsonSuccess struct {
	data    string
	name    string
	value   string
	vType   valueType
	isSlice bool
	isMap   bool
	sep     string
}

func Test_parseJsonData_success_cases(t *testing.T) {
	// Arrange
	cases := []testCaseJsonSuccess{
		{"{\"env 1\":\"value\"}", "env 1", "value", vtString, false, false, ""},
		{"\n{\n\"env 1\"\n:\n\"value\"\n}\n", "env 1", "value", vtString, false, false, ""},
		{"{\"//env 1/**/\":\"//value/**/\"}", "//env 1/**/", "//value/**/", vtString, false, false, ""},
		{" //\n/**/\n{/*w\nw*/\"env_1\"//\n://\n\"value\"//comment\n} /**/ ", "env_1", "value", vtString, false, false, ""},
		{"{\"env_1\":tRUe}", "env_1", "true", vtBool, false, false, ""},
		{"{\"env_1\":FALSe}", "env_1", "false", vtBool, false, false, ""},
		{"{\"env_1\":null}", "env_1", "null", vtNull, false, false, ""},
		{"{\"env_1\":3.14}", "env_1", "3.14", vtNumber, false, false, ""},
		{"{\"env_1\":-123}", "env_1", "-123", vtNumber, false, false, ""},
		{"{\"env_1\":[1.1,2.2,]}", "env_1", "1.1,2.2", vtNumber, true, false, ""},
		{"{\"env_1\":{\"key\":\"value\"}}", "env_1", "value", vtString, false, true, ""},
		{"{\"env_1\":{\"key\":123,},}", "env_1", "123", vtNumber, false, true, ""},
		{"{\"env_1\":/**/[//\n1//\n,//\n2/*\n*/,/**/]}", "env_1", "1,2", vtNumber, true, false, ""},
		{"{\"env_1\":\"va\\\"lue\"}", "env_1", "va\"lue", vtString, false, false, ""},
		{"{\"env_1\":\"va\\\\lue\"}", "env_1", "va\\\\lue", vtString, false, false, ""},
		{"{\"env_1\":[\"v1\",\"v2\"],}", "env_1", "v1,v2", vtString, true, false, ""},
	}

	// Act & Assert
	for i, c := range cases {
		t.Log("Test case "+strconv.Itoa(i+1)+":", c.data)
		test_parseJsonData_success_cases(t, c)
	}
}

func test_parseJsonData_success_cases(t *testing.T, testCase testCaseJsonSuccess) {
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
	err := cr.parseJsonData(r, it, si, defaultJsonData, 0)

	// Assert
	assert.Nil(t, err)
	if testCase.isSlice {
		require.Equal(t, testCase.value, strings.Join(it[testCase.name][0].value.([]string), ","))
	} else if testCase.isMap {
		require.Equal(t, testCase.value, it[testCase.name][0].value.(map[string]string)["key"])
	} else {
		require.Equal(t, testCase.value, it[testCase.name][0].value)
	}
	require.Equal(t, testCase.vType, it[testCase.name][0].valueType)
}

type testCaseJsonWrongName struct {
	data  string
	check string
}

func Test_parseJsonData_wrongName_cases(t *testing.T) {
	// Arrange
	cases := []testCaseJsonWrongName{
		{"{env:1}", "(1:2)"},
		{"{'env':1}", "(1:2)"},
		{"{\"env:1}", "end of file"},
		{"[{\"env\":1}]", "must be an object"},
		{"{\"env\";1}", "(1:7)"},
		{"{\"env\"=1}", "(1:7)"},
		{"{\"env\"\"value\"}", "(1:7)"},
		{"{\"env\"::\"value\"}", "(1:9)"},
		{"{\"env\"::123}", "(1:12)"},
		{"{\"e\"nv\":1}", "(1:5)"},
		{"{\"env", "end of file"},
		{"{\"env\":1", "end of file"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data, c.check)
		test_parseJsonData_wrongName(t, c)
	}
}

func test_parseJsonData_wrongName(t *testing.T, testCase testCaseJsonWrongName) {
	// Arrange
	r := bufio.NewReader(strings.NewReader(testCase.data))
	it := make(intermediateTree)
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}
	cr.data.currentLine = 1

	// Act
	err := cr.parseJsonData(r, it, []structInfo{}, defaultJsonData, 0)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), testCase.check)
}

type testCaseJsonWrongValue struct {
	data  string
	check string
}

func Test_parseJsonData_wrongValue_cases(t *testing.T) {
	// Arrange
	cases := []testCaseJsonWrongValue{
		{"{\"env_1\":'value'}", "(1:17)"},
		{"{\"env_1\":\"value}", "end of file"},
		{"{\"env_1\":\"value'}", "end of file"},
		{"{\"env_1\":'value\"}", "(1:16)"},
		{"{\"env_1\":\"va\"lue\"}", "(1:14)"},
		{"{\"env_1\":tr }", "(1:12)"},
		{"{\"env_1\":flse}", "(1:14)"},
		{"{\"env_1\":nll}", "(1:13)"},
		{"{\"env_1\":1.2.3}", "(1:15)"},
		{"{\"env_1\":[1,2}}", "(1:14)"},
		{"{\"env_1\":{1,2]}", "(1:11)"},
		{"{\"env_1\":{\"e\":1]}", "(1:16)"},
		{"{\"env_1\":1\"env_2\":2}", "(1:11)"},
		{"{\"env_1\":[\"e\":1]}", "(1:14)"},
		{"{\"env_1\":{\"e\",\"v\"}}", "(1:14)"},
		{"{\"env_1\":{1,2}}", "(1:11)"},
		{"{\"env_1\":[1,\"2\"]}", "(1:15)"},
		{"{\"env_1\":[true,1]}", "(1:17)"},
		{"{\"env_1\":[1,,2]}", "(1:13)"},
		{"{\"env_1\":[\"s\"\"s\"]}", "(1:14)"},
		{"{\"env_1\":\"value\"]", "closed with ']'"},
		{"{\"env_1\":}", "(1:10)"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data, c.check)
		test_parseJsonData_wrongValue(t, c)
	}
}

func test_parseJsonData_wrongValue(t *testing.T, testCase testCaseJsonWrongValue) {
	// Arrange
	r := bufio.NewReader(strings.NewReader(testCase.data))
	it := make(intermediateTree)
	si := []structInfo{{keyName: "env_1"}}
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}
	cr.data.currentLine = 1

	// Act
	err := cr.parseJsonData(r, it, si, defaultJsonData, 0)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), testCase.check)
}

type testCaseJsonDuplicateKey struct {
	data    string
	keyName string
}

func Test_parseJsonData_duplicateKey_cases(t *testing.T) {
	// Arrange
	cases := []testCaseJsonDuplicateKey{
		{"{\"env_1\":1,\"env_1\":2}", "env_1"},
		{"{\"env_1\":{\"e1\":1,\"e1\":\"s\"}}", "env_1.e1"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data)
		test_parseJsonData_duplicateKey(t, c)
	}
}
func test_parseJsonData_duplicateKey(t *testing.T, testCase testCaseJsonDuplicateKey) {
	// Arrange
	r := bufio.NewReader(strings.NewReader(testCase.data))
	it := make(intermediateTree)
	si := []structInfo{{keyName: testCase.keyName}}
	cr := &configReader{options: ConfigOptions{RewriteValues: true}}
	cr.data.currentLine = 1

	// Act
	err := cr.parseJsonData(r, it, si, defaultJsonData, 0)

	// Assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "duplicate key: "+testCase.keyName)
}
