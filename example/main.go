package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	goc "github.com/mansiper/goconfiguration"
)

type config struct {
	Field_1  int                `env:"f1,required" def:"10"`
	Field_2  []string           `env:"f2,append" def:"a|b" sep:"|"`
	Field_3  [3]time.Time       `env:"f3" def:"now,now"`
	Field_4  map[string]float64 `env:"f4" def:"a?1.1,b?2.2" sep2:"?"`
	Field_5  subConfig          `env:"f5"`
	Fiend_6  string             `env:"GOPATH"`
	Fiend_7  *float64           `env:"-"`
	Field_8  []*int8            `env:"f8" def:"*nil,1,2,3,4,5"`
	FIell_9  [2]*float32        `env:"f9,required" def:"1.23"`
	Field_10 parsableType       `env:"f10,useparser" def:"f1_1"`
}
type subConfig struct {
	Field_1 *uint32 `env:"sf1" def:"*nil"`
}
type parsableType struct {
	Field_1 string
	Field_2 int
}

func main() {
	config := &config{}
	cr := goc.NewConfigReader().
		RewriteValues(true).
		WithParser("f10", ParseSubConfig).
		AddEnvironment().
		AddFile(".env").
		AddFile("config.ini").
		AddString("f5.sf1 = 10", goc.FtEnv, "config 1").
		AddFile("config.json").
		EnsureHasNoErrors()
	err := cr.ReadConfig(config)
	if err != nil {
		panic(err)
	}

	js, _ := json.Marshal(config)
	fmt.Println(string(js))

	/* Result example:
	{"Field_1":-12345,"Field_2":["aa","bb"],"Field_3":["2024-07-01T12:30:01Z","2024-06-28T12:54:35.2429406-03:00","0001-01-01T00:00:00Z"],
	"Field_4":{"e":2.7182,"pi":3.1415},"Field_5":{"Field_1":10},"Fiend_6":"C:\\Code\\Go","Fiend_7":null,"Field_8":[null,1,2,null],
	"FIell_9":[1.1,2.2],"Field_10":{"Field_1":"answer","Field_2":42}}
	*/
}

func ParseSubConfig(s string) (interface{}, error) {
	split := strings.Split(s, "_")
	i, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, err
	}

	return parsableType{Field_1: split[0], Field_2: i}, nil
}
