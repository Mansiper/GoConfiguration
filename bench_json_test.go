package configuration

import (
	"encoding/json"
	"testing"
)

type benchmarkConfig struct {
	F1  int                `env:"f1" json:"f1"`
	F2  bool               `env:"f2" json:"f2"`
	F3  uint               `env:"f3" json:"f3"`
	F4  float32            `env:"f4" json:"f4"`
	F5  string             `env:"f5" json:"f5"`
	F6  []int              `env:"f6" json:"f6"`
	F7  [2]bool            `env:"f7" json:"f7"`
	F8  map[string]string  `env:"f8" json:"f8"`
	F9  subBenchmarkConfig `env:"f9" json:"f9"`
	F10 *int               `env:"f10" json:"f10"`
}
type subBenchmarkConfig struct {
	F1 int    `env:"f1" json:"f1"`
	F2 string `env:"f2" json:"f2"`
}

const benchmarkJsonString = `{
	"f1": 1,
	"f2": true,
	"f3": 3,
	"f4": 4.4,
	"f5": "5",
	"f6": [6, 6],
	"f7": [true, false],
	"f8": {"8": "8"},
	"f9": {"f1": 9, "f2": "9"},
	"f10": 10
}`

func BenchmarkBasicJson(b *testing.B) {
	config := &benchmarkConfig{}
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal([]byte(benchmarkJsonString), config)
	}
}

func BenchmarkMyJson(b *testing.B) {
	cr := NewConfigReader().
		AddString(benchmarkJsonString, FtJson, "json")
	b.ResetTimer()

	config := &benchmarkConfig{}
	for i := 0; i < b.N; i++ {
		_ = cr.ReadConfig(config)
	}
}

//go test -benchmem -run=^$ -bench ^BenchmarkBasicJson$ configuration -benchmem
