package configuration

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func addr[T any](i T) *T {
	return &i
}

type testCaseSetValueTest struct {
	F1  int
	F2  int8
	F3  int16
	F4  int32
	F5  int64
	F6  uint
	F7  uint8
	F8  uint16
	F9  uint32
	F10 uint64
	F11 float32
	F12 float64
	F13 string
	F14 bool
	F15 time.Time
	F16 time.Duration
	F17 struct {
		F1 int
	}

	F1s  []int
	F2s  []int8
	F3s  []int16
	F4s  []int32
	F5s  []int64
	F6s  []uint
	F7s  []uint8
	F8s  []uint16
	F9s  []uint32
	F10s []uint64
	F11s []float32
	F12s []float64
	F13s []string
	F14s []bool
	F15s []time.Time
	F16s []time.Duration
	F17s struct {
		F1s []int
	}

	F1p  *int
	F2p  *int8
	F3p  *int16
	F4p  *int32
	F5p  *int64
	F6p  *uint
	F7p  *uint8
	F8p  *uint16
	F9p  *uint32
	F10p *uint64
	F11p *float32
	F12p *float64
	F13p *string
	F14p *bool
	F15p *time.Time
	F16p *time.Duration
	F17p struct {
		F1p *int
	}
	F18p *int

	F1a  [3]int
	F2a  [3]int8
	F3a  [3]int16
	F4a  [3]int32
	F5a  [3]int64
	F6a  [3]uint
	F7a  [3]uint8
	F8a  [3]uint16
	F9a  [3]uint32
	F10a [3]uint64
	F11a [3]float32
	F12a [3]float64
	F13a [3]string
	F14a [3]bool
	F15a [3]time.Time
	F16a [3]time.Duration
	F17a struct {
		F1a [3]int
	}

	F1m  map[string]int
	F2m  map[string]int8
	F3m  map[string]int16
	F4m  map[string]int32
	F5m  map[string]int64
	F6m  map[string]uint
	F7m  map[string]uint8
	F8m  map[string]uint16
	F9m  map[string]uint32
	F10m map[string]uint64
	F11m map[string]float32
	F12m map[string]float64
	F13m map[string]string
	F14m map[string]bool
	F15m map[string]time.Time
	F16m map[string]time.Duration
	F17m struct {
		F1m map[string]int
	}

	F1sp  []*int
	F2sp  []*int8
	F3sp  []*int16
	F4sp  []*int32
	F5sp  []*int64
	F6sp  []*uint
	F7sp  []*uint8
	F8sp  []*uint16
	F9sp  []*uint32
	F10sp []*uint64
	F11sp []*float32
	F12sp []*float64
	F13sp []*string
	F14sp []*bool
	F15sp []*time.Time
	F16sp []*time.Duration
	F17sp struct {
		F1sp []*int
	}
	F18sp []*int

	F1ap  [3]*int
	F2ap  [3]*int8
	F3ap  [3]*int16
	F4ap  [3]*int32
	F5ap  [3]*int64
	F6ap  [3]*uint
	F7ap  [3]*uint8
	F8ap  [3]*uint16
	F9ap  [3]*uint32
	F10ap [3]*uint64
	F11ap [3]*float32
	F12ap [3]*float64
	F13ap [3]*string
	F14ap [3]*bool
	F15ap [3]*time.Time
	F16ap [3]*time.Duration
	F17ap struct {
		F1ap [3]*int
	}
	F18ap [3]*int
}

func Test_getStructInfo_setValues_success(t *testing.T) {
	// Arrange
	cr := &configReader{}
	config := &testCaseSetValueTest{}
	it := intermediateTree{
		"f1":         []intermediateData{{value: "1", source: 0, valueType: vtAny}},
		"f2":         []intermediateData{{value: "2", source: 0, valueType: vtAny}},
		"f3":         []intermediateData{{value: "3", source: 0, valueType: vtAny}},
		"f4":         []intermediateData{{value: "4", source: 0, valueType: vtAny}},
		"f5":         []intermediateData{{value: "5", source: 0, valueType: vtAny}},
		"f6":         []intermediateData{{value: "6", source: 0, valueType: vtAny}},
		"f7":         []intermediateData{{value: "7", source: 0, valueType: vtAny}},
		"f8":         []intermediateData{{value: "8", source: 0, valueType: vtAny}},
		"f9":         []intermediateData{{value: "9", source: 0, valueType: vtAny}},
		"f10":        []intermediateData{{value: "10", source: 0, valueType: vtAny}},
		"f11":        []intermediateData{{value: "11.1", source: 0, valueType: vtAny}},
		"f12":        []intermediateData{{value: "12.1", source: 0, valueType: vtAny}},
		"f13":        []intermediateData{{value: "13", source: 0, valueType: vtAny}},
		"f14":        []intermediateData{{value: "true", source: 0, valueType: vtAny}},
		"f15":        []intermediateData{{value: "2020-01-01T00:00:00Z", source: 0, valueType: vtAny}},
		"f16":        []intermediateData{{value: "1h", source: 0, valueType: vtAny}},
		"f17.f1":     []intermediateData{{value: "1", source: 0, valueType: vtAny}},
		"f1s":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f2s":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f3s":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f4s":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f5s":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f6s":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f7s":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f8s":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f9s":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f10s":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f11s":       []intermediateData{{value: []string{"1.1", "2.2"}, source: 0, valueType: vtAny}},
		"f12s":       []intermediateData{{value: []string{"1.1", "2.2"}, source: 0, valueType: vtAny}},
		"f13s":       []intermediateData{{value: []string{"a", "b"}, source: 0, valueType: vtAny}},
		"f14s":       []intermediateData{{value: []string{"true", "false"}, source: 0, valueType: vtAny}},
		"f15s":       []intermediateData{{value: []string{"2020-01-01T00:00:00Z", "2020-01-01T00:00:00Z"}, source: 0, valueType: vtAny}},
		"f16s":       []intermediateData{{value: []string{"1h", "1m"}, source: 0, valueType: vtAny}},
		"f17s.f1s":   []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f1p":        []intermediateData{{value: "1", source: 0, valueType: vtAny}},
		"f2p":        []intermediateData{{value: "2", source: 0, valueType: vtAny}},
		"f3p":        []intermediateData{{value: "3", source: 0, valueType: vtAny}},
		"f4p":        []intermediateData{{value: "4", source: 0, valueType: vtAny}},
		"f5p":        []intermediateData{{value: "5", source: 0, valueType: vtAny}},
		"f6p":        []intermediateData{{value: "6", source: 0, valueType: vtAny}},
		"f7p":        []intermediateData{{value: "7", source: 0, valueType: vtAny}},
		"f8p":        []intermediateData{{value: "8", source: 0, valueType: vtAny}},
		"f9p":        []intermediateData{{value: "9", source: 0, valueType: vtAny}},
		"f10p":       []intermediateData{{value: "10", source: 0, valueType: vtAny}},
		"f11p":       []intermediateData{{value: "11.1", source: 0, valueType: vtAny}},
		"f12p":       []intermediateData{{value: "12.1", source: 0, valueType: vtAny}},
		"f13p":       []intermediateData{{value: "13", source: 0, valueType: vtAny}},
		"f14p":       []intermediateData{{value: "true", source: 0, valueType: vtAny}},
		"f15p":       []intermediateData{{value: "2020-01-01T00:00:00Z", source: 0, valueType: vtAny}},
		"f16p":       []intermediateData{{value: "1h", source: 0, valueType: vtAny}},
		"f17p.f1p":   []intermediateData{{value: "1", source: 0, valueType: vtAny}},
		"f18p":       []intermediateData{{value: nilDefault, source: 0, valueType: vtAny}},
		"f1a":        []intermediateData{{value: []string{}, source: 0, valueType: vtAny}},
		"f2a":        []intermediateData{{value: []string{"1"}, source: 0, valueType: vtAny}},
		"f3a":        []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f4a":        []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f5a":        []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f6a":        []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f7a":        []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f8a":        []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f9a":        []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f10a":       []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f11a":       []intermediateData{{value: []string{"1.1", "2.2", "3.3"}, source: 0, valueType: vtAny}},
		"f12a":       []intermediateData{{value: []string{"1.1", "2.2", "3.3"}, source: 0, valueType: vtAny}},
		"f13a":       []intermediateData{{value: []string{"a", "b", "c"}, source: 0, valueType: vtAny}},
		"f14a":       []intermediateData{{value: []string{"true", "false", "true"}, source: 0, valueType: vtAny}},
		"f15a":       []intermediateData{{value: []string{"2020-01-01T00:00:00Z", "2020-01-01T00:00:00Z", "2020-01-01T00:00:00Z"}, source: 0, valueType: vtAny}},
		"f16a":       []intermediateData{{value: []string{"1h", "1m", "1h"}, source: 0, valueType: vtAny}},
		"f17a.f1a":   []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f1m":        []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f2m":        []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f3m":        []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f4m":        []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f5m":        []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f6m":        []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f7m":        []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f8m":        []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f9m":        []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f10m":       []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f11m":       []intermediateData{{value: map[string]string{"a": "1.1"}, source: 0, valueType: vtAny}},
		"f12m":       []intermediateData{{value: map[string]string{"a": "2.2"}, source: 0, valueType: vtAny}},
		"f13m":       []intermediateData{{value: map[string]string{"a": "abc"}, source: 0, valueType: vtAny}},
		"f14m":       []intermediateData{{value: map[string]string{"a": "true"}, source: 0, valueType: vtAny}},
		"f15m":       []intermediateData{{value: map[string]string{"a": "2020-01-01T00:00:00Z"}, source: 0, valueType: vtAny}},
		"f16m":       []intermediateData{{value: map[string]string{"a": "1h"}, source: 0, valueType: vtAny}},
		"f17m.f1m":   []intermediateData{{value: map[string]string{"a": "1"}, source: 0, valueType: vtAny}},
		"f1sp":       []intermediateData{{value: []string{"1", "*nil"}, source: 0, valueType: vtAny}},
		"f2sp":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f3sp":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f4sp":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f5sp":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f6sp":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f7sp":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f8sp":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f9sp":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f10sp":      []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f11sp":      []intermediateData{{value: []string{"1.1", "2.2"}, source: 0, valueType: vtAny}},
		"f12sp":      []intermediateData{{value: []string{"1.1", "2.2"}, source: 0, valueType: vtAny}},
		"f13sp":      []intermediateData{{value: []string{"a", "b"}, source: 0, valueType: vtAny}},
		"f14sp":      []intermediateData{{value: []string{"true", "false"}, source: 0, valueType: vtAny}},
		"f15sp":      []intermediateData{{value: []string{"2020-01-01T00:00:00Z", "2020-01-01T00:00:00Z"}, source: 0, valueType: vtAny}},
		"f16sp":      []intermediateData{{value: []string{"1h", "1m"}, source: 0, valueType: vtAny}},
		"f17sp.f1sp": []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f18sp":      []intermediateData{{value: []string{nilDefault}, source: 0, valueType: vtAny}},
		"f1ap":       []intermediateData{{value: []string{}, source: 0, valueType: vtAny}},
		"f2ap":       []intermediateData{{value: []string{"1"}, source: 0, valueType: vtAny}},
		"f3ap":       []intermediateData{{value: []string{"1", "2"}, source: 0, valueType: vtAny}},
		"f4ap":       []intermediateData{{value: []string{"1", "2", "*nil"}, source: 0, valueType: vtAny}},
		"f5ap":       []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f6ap":       []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f7ap":       []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f8ap":       []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f9ap":       []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f10ap":      []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f11ap":      []intermediateData{{value: []string{"1.1", "2.2", "3.3"}, source: 0, valueType: vtAny}},
		"f12ap":      []intermediateData{{value: []string{"1.1", "2.2", "3.3"}, source: 0, valueType: vtAny}},
		"f13ap":      []intermediateData{{value: []string{"a", "b", "c"}, source: 0, valueType: vtAny}},
		"f14ap":      []intermediateData{{value: []string{"true", "false", "true"}, source: 0, valueType: vtAny}},
		"f15ap":      []intermediateData{{value: []string{"2020-01-01T00:00:00Z", "2020-01-01T00:00:00Z", "2020-01-01T00:00:00Z"}, source: 0, valueType: vtAny}},
		"f16ap":      []intermediateData{{value: []string{"1h", "1m", "1h"}, source: 0, valueType: vtAny}},
		"f17ap.f1ap": []intermediateData{{value: []string{"1", "2", "3"}, source: 0, valueType: vtAny}},
		"f18ap":      []intermediateData{{value: []string{nilDefault, nilDefault, nilDefault}, source: 0, valueType: vtAny}},
	}
	var err error = nil
	var err2 error = nil

	// Act
	si, err := cr.getStructInfo(config, "", "")
	if err == nil {
		err2 = cr.setValues(it, si)
	}

	// Assert
	require.Nil(t, err)
	require.Nil(t, err2)
	require.Equal(t, 1, config.F1)
	require.Equal(t, int8(2), config.F2)
	require.Equal(t, int16(3), config.F3)
	require.Equal(t, int32(4), config.F4)
	require.Equal(t, int64(5), config.F5)
	require.Equal(t, uint(6), config.F6)
	require.Equal(t, uint8(7), config.F7)
	require.Equal(t, uint16(8), config.F8)
	require.Equal(t, uint32(9), config.F9)
	require.Equal(t, uint64(10), config.F10)
	require.Equal(t, float32(11.1), config.F11)
	require.Equal(t, float64(12.1), config.F12)
	require.Equal(t, "13", config.F13)
	require.Equal(t, true, config.F14)
	require.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), config.F15)
	require.Equal(t, time.Hour, config.F16)
	require.Equal(t, 1, config.F17.F1)
	require.Equal(t, []int{1, 2}, config.F1s)
	require.Equal(t, []int8{1, 2}, config.F2s)
	require.Equal(t, []int16{1, 2}, config.F3s)
	require.Equal(t, []int32{1, 2}, config.F4s)
	require.Equal(t, []int64{1, 2}, config.F5s)
	require.Equal(t, []uint{1, 2}, config.F6s)
	require.Equal(t, []uint8{1, 2}, config.F7s)
	require.Equal(t, []uint16{1, 2}, config.F8s)
	require.Equal(t, []uint32{1, 2}, config.F9s)
	require.Equal(t, []uint64{1, 2}, config.F10s)
	require.Equal(t, []float32{1.1, 2.2}, config.F11s)
	require.Equal(t, []float64{1.1, 2.2}, config.F12s)
	require.Equal(t, []string{"a", "b"}, config.F13s)
	require.Equal(t, []bool{true, false}, config.F14s)
	require.Equal(t, []time.Time{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}, config.F15s)
	require.Equal(t, []time.Duration{time.Hour, time.Minute}, config.F16s)
	require.Equal(t, []int{1, 2}, config.F17s.F1s)
	require.Equal(t, 1, *config.F1p)
	require.Equal(t, int8(2), *config.F2p)
	require.Equal(t, int16(3), *config.F3p)
	require.Equal(t, int32(4), *config.F4p)
	require.Equal(t, int64(5), *config.F5p)
	require.Equal(t, uint(6), *config.F6p)
	require.Equal(t, uint8(7), *config.F7p)
	require.Equal(t, uint16(8), *config.F8p)
	require.Equal(t, uint32(9), *config.F9p)
	require.Equal(t, uint64(10), *config.F10p)
	require.Equal(t, float32(11.1), *config.F11p)
	require.Equal(t, float64(12.1), *config.F12p)
	require.Equal(t, "13", *config.F13p)
	require.Equal(t, true, *config.F14p)
	require.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), *config.F15p)
	require.Equal(t, time.Hour, *config.F16p)
	require.Equal(t, 1, *config.F17p.F1p)
	require.Nil(t, config.F18p)
	require.Equal(t, [3]int{0, 0, 0}, config.F1a)
	require.Equal(t, [3]int8{1, 0, 0}, config.F2a)
	require.Equal(t, [3]int16{1, 2, 0}, config.F3a)
	require.Equal(t, [3]int32{1, 2, 3}, config.F4a)
	require.Equal(t, [3]int64{1, 2, 3}, config.F5a)
	require.Equal(t, [3]uint{1, 2, 3}, config.F6a)
	require.Equal(t, [3]uint8{1, 2, 3}, config.F7a)
	require.Equal(t, [3]uint16{1, 2, 3}, config.F8a)
	require.Equal(t, [3]uint32{1, 2, 3}, config.F9a)
	require.Equal(t, [3]uint64{1, 2, 3}, config.F10a)
	require.Equal(t, [3]float32{1.1, 2.2, 3.3}, config.F11a)
	require.Equal(t, [3]float64{1.1, 2.2, 3.3}, config.F12a)
	require.Equal(t, [3]string{"a", "b", "c"}, config.F13a)
	require.Equal(t, [3]bool{true, false, true}, config.F14a)
	require.Equal(t, [3]time.Time{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}, config.F15a)
	require.Equal(t, [3]time.Duration{time.Hour, time.Minute, time.Hour}, config.F16a)
	require.Equal(t, [3]int{1, 2, 3}, config.F17a.F1a)
	require.Equal(t, map[string]int{"a": 1}, config.F1m)
	require.Equal(t, map[string]int8{"a": 1}, config.F2m)
	require.Equal(t, map[string]int16{"a": 1}, config.F3m)
	require.Equal(t, map[string]int32{"a": 1}, config.F4m)
	require.Equal(t, map[string]int64{"a": 1}, config.F5m)
	require.Equal(t, map[string]uint{"a": 1}, config.F6m)
	require.Equal(t, map[string]uint8{"a": 1}, config.F7m)
	require.Equal(t, map[string]uint16{"a": 1}, config.F8m)
	require.Equal(t, map[string]uint32{"a": 1}, config.F9m)
	require.Equal(t, map[string]uint64{"a": 1}, config.F10m)
	require.Equal(t, map[string]float32{"a": 1.1}, config.F11m)
	require.Equal(t, map[string]float64{"a": 2.2}, config.F12m)
	require.Equal(t, map[string]string{"a": "abc"}, config.F13m)
	require.Equal(t, map[string]bool{"a": true}, config.F14m)
	require.Equal(t, map[string]time.Time{"a": time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}, config.F15m)
	require.Equal(t, map[string]time.Duration{"a": time.Hour}, config.F16m)
	require.Equal(t, map[string]int{"a": 1}, config.F17m.F1m)
	require.Equal(t, []*int{addr(1), nil}, config.F1sp)
	require.Equal(t, []*int8{addr(int8(1)), addr(int8(2))}, config.F2sp)
	require.Equal(t, []*int16{addr(int16(1)), addr(int16(2))}, config.F3sp)
	require.Equal(t, []*int32{addr(int32(1)), addr(int32(2))}, config.F4sp)
	require.Equal(t, []*int64{addr(int64(1)), addr(int64(2))}, config.F5sp)
	require.Equal(t, []*uint{addr(uint(1)), addr(uint(2))}, config.F6sp)
	require.Equal(t, []*uint8{addr(uint8(1)), addr(uint8(2))}, config.F7sp)
	require.Equal(t, []*uint16{addr(uint16(1)), addr(uint16(2))}, config.F8sp)
	require.Equal(t, []*uint32{addr(uint32(1)), addr(uint32(2))}, config.F9sp)
	require.Equal(t, []*uint64{addr(uint64(1)), addr(uint64(2))}, config.F10sp)
	require.Equal(t, []*float32{addr(float32(1.1)), addr(float32(2.2))}, config.F11sp)
	require.Equal(t, []*float64{addr(float64(1.1)), addr(float64(2.2))}, config.F12sp)
	require.Equal(t, []*string{addr("a"), addr("b")}, config.F13sp)
	require.Equal(t, []*bool{addr(true), addr(false)}, config.F14sp)
	require.Equal(t, []*time.Time{addr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)), addr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))}, config.F15sp)
	require.Equal(t, []*time.Duration{addr(time.Hour), addr(time.Minute)}, config.F16sp)
	require.Equal(t, []*int{addr(1), addr(2)}, config.F17sp.F1sp)
	require.Equal(t, []*int{nil}, config.F18sp)
	require.Equal(t, [3]*int{nil, nil, nil}, config.F1ap)
	require.Equal(t, [3]*int8{addr(int8(1)), nil, nil}, config.F2ap)
	require.Equal(t, [3]*int16{addr(int16(1)), addr(int16(2)), nil}, config.F3ap)
	require.Equal(t, [3]*int32{addr(int32(1)), addr(int32(2)), nil}, config.F4ap)
	require.Equal(t, [3]*int64{addr(int64(1)), addr(int64(2)), addr(int64(3))}, config.F5ap)
	require.Equal(t, [3]*uint{addr(uint(1)), addr(uint(2)), addr(uint(3))}, config.F6ap)
	require.Equal(t, [3]*uint8{addr(uint8(1)), addr(uint8(2)), addr(uint8(3))}, config.F7ap)
	require.Equal(t, [3]*uint16{addr(uint16(1)), addr(uint16(2)), addr(uint16(3))}, config.F8ap)
	require.Equal(t, [3]*uint32{addr(uint32(1)), addr(uint32(2)), addr(uint32(3))}, config.F9ap)
	require.Equal(t, [3]*uint64{addr(uint64(1)), addr(uint64(2)), addr(uint64(3))}, config.F10ap)
	require.Equal(t, [3]*float32{addr(float32(1.1)), addr(float32(2.2)), addr(float32(3.3))}, config.F11ap)
	require.Equal(t, [3]*float64{addr(1.1), addr(2.2), addr(3.3)}, config.F12ap)
	require.Equal(t, [3]*string{addr("a"), addr("b"), addr("c")}, config.F13ap)
	require.Equal(t, [3]*bool{addr(true), addr(false), addr(true)}, config.F14ap)
	require.Equal(t, [3]*time.Time{addr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)), addr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)), addr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))}, config.F15ap)
	require.Equal(t, [3]*time.Duration{addr(time.Hour), addr(time.Minute), addr(time.Hour)}, config.F16ap)
	require.Equal(t, [3]*int{addr(1), addr(2), addr(3)}, config.F17ap.F1ap)
	require.Equal(t, [3]*int{nil, nil, nil}, config.F18ap)
}

type testCaseSetValueDefaulSuccess struct {
	F1  int            `def:"-1"`
	F2  uint           `def:"1"`
	F3  float32        `def:"2.2"`
	F4  string         `def:"a"`
	F5  bool           `def:"true"`
	F6  time.Time      `def:"2020-01-01T12:34:56Z"`
	F7  time.Duration  `def:"1h"`
	F8  []int          `def:"1,2,3"`
	F9  []string       `def:"a;b;c" sep:";"`
	F10 []float32      `def:"1.1,2.2,3.3"`
	F11 time.Time      `def:"now"`
	F12 *int8          `def:"-1"`
	F13 *uint16        `def:"2"`
	F14 *string        `def:"a"`
	F15 *bool          `def:"true"`
	F16 *float64       `def:"1.1"`
	F17 *time.Time     `def:"2020-01-01T12:34:56Z"`
	F18 *time.Duration `def:"1h"`
	F19 *int           `def:"*nil"`
	F20 *string        `def:"*nil"`
	F21 string         `def:""`
	F22 string         `env:"F22,required" def:"a"`
	F23 []string
	F24 []string `def:""`
	F25 []string `def:"*nil"`
	F26 [2]string
	F27 [2]int     `def:"1,2"`
	F28 [2]float32 `def:"1.1"`
	F29 [2]float32 `def:"1.1|2.2" sep:"|"`
	F30 map[string]int
	F31 map[string]int `def:""`
	F32 map[string]int `def:"*nil"`
	F33 map[string]int `def:"a:1"`
	F34 map[string]int `def:"a:1,b:2"`
	F35 []*int         `def:"1,2"`
	F36 []*int         `def:"*nil,2,*nil"`
	F37 []*string      `def:"a;b;*nil" sep:";"`
	F38 []*string
	F39 []*string
	F40 []*string `def:""`
	F41 [2]*string
	F42 [2]*int     `def:"1,*nil"`
	F43 [2]*float64 `def:"1.1"`
	F44 [2]*float64 `def:"*nil|2.2" sep:"|"`
}

func Test_setValue_default_success(t *testing.T) {
	// Arrange
	cr := &configReader{}
	config := testCaseSetValueDefaulSuccess{}
	it := intermediateTree{}
	for i := 1; i <= 29; i++ {
		name := "f" + strconv.Itoa(i)
		if i == 8 || i == 9 || i == 10 || i == 35 || i == 36 || i == 39 || i == 40 {
			it[name] = []intermediateData{{value: []string{}, source: 0, valueType: vtAny}}
		} else if i == 23 || i == 24 || i == 25 || i == 37 || i == 38 {
			it[name] = []intermediateData{{value: []string{""}, source: 0, valueType: vtAny}}
		} else if i == 26 || i == 27 || i == 28 || i == 29 || i == 41 || i == 42 || i == 43 || i == 44 {
			it[name] = []intermediateData{{value: []string{""}, source: 0, valueType: vtAny}}
		} else if i == 30 || i == 31 || i == 32 || i == 33 || i == 34 {
			it[name] = []intermediateData{{value: map[string]string{}, source: 0, valueType: vtAny}}
		} else {
			it[name] = []intermediateData{{value: "", source: 0, valueType: vtAny}}
		}
	}

	// Act
	si, err := cr.getStructInfo(&config, "", "")
	if err == nil {
		err = cr.setValues(it, si)
	}

	// Assert
	require.Nil(t, err)
	require.Equal(t, -1, config.F1)
	require.Equal(t, uint(1), config.F2)
	require.Equal(t, float32(2.2), config.F3)
	require.Equal(t, "a", config.F4)
	require.Equal(t, true, config.F5)
	require.Equal(t, time.Date(2020, 1, 1, 12, 34, 56, 0, time.UTC), config.F6)
	require.Equal(t, time.Hour, config.F7)
	require.Equal(t, []int{1, 2, 3}, config.F8)
	require.Equal(t, []string{"a", "b", "c"}, config.F9)
	require.Equal(t, []float32{1.1, 2.2, 3.3}, config.F10)
	require.Equal(t, time.Now().Year(), config.F11.Year())
	require.Equal(t, int8(-1), *config.F12)
	require.Equal(t, uint16(2), *config.F13)
	require.Equal(t, "a", *config.F14)
	require.Equal(t, true, *config.F15)
	require.Equal(t, 1.1, *config.F16)
	require.Equal(t, time.Date(2020, 1, 1, 12, 34, 56, 0, time.UTC), *config.F17)
	require.Equal(t, time.Hour, *config.F18)
	require.Nil(t, config.F19)
	require.Nil(t, config.F20)
	require.Equal(t, "", config.F21)
	require.Equal(t, "a", config.F22)
	require.Nil(t, config.F23)
	require.Equal(t, []string{}, config.F24)
	require.Nil(t, config.F25)
	require.Equal(t, [2]string{"", ""}, config.F26)
	require.Equal(t, [2]int{1, 2}, config.F27)
	require.Equal(t, [2]float32{1.1, 0}, config.F28)
	require.Equal(t, [2]float32{1.1, 2.2}, config.F29)
	require.Nil(t, config.F30)
	require.Nil(t, config.F31)
	require.Nil(t, config.F32)
	require.Equal(t, map[string]int{"a": 1}, config.F33)
	require.Equal(t, map[string]int{"a": 1, "b": 2}, config.F34)
	require.Equal(t, []*int{addr(1), addr(2)}, config.F35)
	require.Equal(t, []*int{nil, addr(2), nil}, config.F36)
	require.Equal(t, []*string{addr("a"), addr("b"), nil}, config.F37)
	require.Nil(t, config.F38)
	require.Nil(t, config.F39)
	require.Nil(t, config.F40)
	require.Equal(t, [2]*string{nil, nil}, config.F41)
	require.Equal(t, [2]*int{addr(1), nil}, config.F42)
	require.Equal(t, [2]*float64{addr(1.1), nil}, config.F43)
	require.Equal(t, [2]*float64{nil, addr(2.2)}, config.F44)
}

type testCaseSetValueDefaulError struct {
	config  interface{}
	isSLice bool
	isMap   bool
	err     string
}

func Test_setValue_default_error_cases(t *testing.T) {
	// Arrange
	c1 := struct {
		Field int `env:"f,required" def:""`
	}{}
	c2 := struct {
		Field int `env:"f,required"`
	}{}
	c3 := struct {
		Field []int `env:"f,required" def:""`
	}{}
	c4 := struct {
		Field []int `env:"f,required"`
	}{}
	c5 := struct {
		Field []int `env:"unknown,required"`
	}{}
	c6 := struct {
		Field [1]int `env:"f,required"`
	}{}
	c7 := struct {
		Field map[string]int `env:"f,required"`
	}{}
	c8 := struct {
		Field map[string]int `env:"f,required" def:""`
	}{}
	c9 := struct {
		Field [2]int `def:"1,2,3"`
	}{}
	c10 := struct {
		Field []*int `env:"f,required" def:""`
	}{}
	c11 := struct {
		Field [1]*int `env:"f,required"`
	}{}
	cases := []testCaseSetValueDefaulError{
		{&c1, false, false, "required field Field is empty"},
		{&c2, false, false, "required field Field is empty"},
		{&c3, true, false, "required field Field is empty"},
		{&c4, true, false, "required field Field is empty"},
		{&c5, true, false, "required field Field value is missing"},
		{&c6, true, false, "required field Field is empty"},
		{&c7, false, true, "required field Field is empty"},
		{&c8, false, true, "required field Field is empty"},
		{&c9, false, false, "field Field has more values than allowed: 1,2,3"},
		{&c10, true, false, "required field Field is empty"},
		{&c11, true, false, "required field Field is empty"},
	}

	// Act & Assert
	for i, c := range cases {
		t.Log("Test case:", i)
		test_setValue_default_error(t, c)
	}
}

func test_setValue_default_error(t *testing.T, testCase testCaseSetValueDefaulError) {
	// Arrange
	cr := &configReader{}
	it := intermediateTree{}
	if testCase.isSLice {
		it["f"] = []intermediateData{{value: []string{}, source: 0}}
	} else if testCase.isMap {
		it["f"] = []intermediateData{{value: map[string]string{}, source: 0}}
	} else {
		it["f"] = []intermediateData{{value: "", source: 0}}
	}

	// Act
	si, err := cr.getStructInfo(testCase.config, "", "")
	if err == nil {
		err = cr.setValues(it, si)
	}

	// Assert
	require.NotNil(t, err)
	require.Equal(t, testCase.err, err.Error())
}

func Test_setValue_useParser(t *testing.T) {
	// Arrange
	type subType struct {
		Field int
	}
	type rootType struct {
		Field subType `env:"sub,useparser"`
	}
	cr := &configReader{}
	cr.options.Parsers = make(map[string]Parser)
	cr.options.Parsers["sub"] =
		func(value string) (interface{}, error) {
			i, err := strconv.Atoi(value)
			return subType{Field: i}, err
		}
	config := &rootType{}
	it := intermediateTree{
		"sub": []intermediateData{{value: "1", source: 0, valueType: vtAny}},
	}

	// Act
	si, err := cr.getStructInfo(config, "", "")
	if err == nil {
		err = cr.setValues(it, si)
	}

	// Assert
	require.Nil(t, err)
	require.Equal(t, 1, config.Field.Field)
}

type testCaseGetStructInfoTagSuccess struct {
	config interface{}
	def    string
	sep    string
	sep2   string
}

func Test_getStructInfo_success_tag_cases(t *testing.T) {
	// Arrange
	c1 := struct {
		Field string `env:"e" def:"value" sep:";"`
	}{}
	c2 := struct {
		Field string `def:"value" sep:";"`
	}{}
	c3 := struct {
		Field []string `sep:";"`
	}{}
	c4 := struct {
		Field []string `def:"value"`
	}{}
	c5 := struct {
		Field [1]string `def:"value" sep:"|"`
	}{}
	c6 := struct {
		Field [1]string `sep:"|"`
	}{}
	c7 := struct {
		Field map[string]string `def:"a?b!c?d" sep:"!" sep2:"?"`
	}{}
	c8 := struct {
		Field time.Time `def:"now"`
	}{}
	c9 := struct {
		Field []*int `def:"1" sep:";"`
	}{}
	c10 := struct {
		Field [1]*string `def:"value" sep:";"`
	}{}
	cases := []testCaseGetStructInfoTagSuccess{
		{&c1, "value", ";", ":"},
		{&c2, "value", ";", ":"},
		{&c3, nilDefault, ";", ":"},
		{&c4, "value", ",", ":"},
		{&c5, "value", "|", ":"},
		{&c6, "", "|", ":"},
		{&c7, "a?b!c?d", "!", "?"},
		{&c8, nowTime, ",", ":"},
		{&c9, "1", ";", ":"},
		{&c10, "value", ";", ":"},
	}

	// Act & Assert
	for i, c := range cases {
		t.Log("Test case:", i)
		test_getStructInfo_tag_success(t, c)
	}
}

func test_getStructInfo_tag_success(t *testing.T, testCase testCaseGetStructInfoTagSuccess) {
	// Arrange
	cr := &configReader{}

	// Act
	si, err := cr.getStructInfo(testCase.config, "", "")

	// Assert
	require.Nil(t, err)
	require.Len(t, si, 1)
	require.Equal(t, testCase.def, si[0].defValue)
	require.Equal(t, testCase.sep, si[0].separator)
}

type testCaseGetStructInfoError struct {
	config interface{}
	err    string
}

func Test_getStructInfo_error_cases(t *testing.T) {
	// Arrange
	var c1 interface{} = nil
	c2 := 1
	c3 := "value"
	c4 := testCaseSetValueTest{}
	c5 := []int{}
	c6 := make([]int, 1)
	c7 := map[string]string{}
	c8 := []*int{}
	c9 := make([]*int, 1)
	cases := []testCaseGetStructInfoError{
		{c1, "user config is nil"},
		{c2, "pass your config struct as a pointer"},
		{&c2, "user config must be a struct"},
		{c3, "pass your config struct as a pointer"},
		{&c3, "user config must be a struct"},
		{c4, "pass your config struct as a pointer"},
		{c5, "pass your config struct as a pointer"},
		{&c5, "user config must be a struct"},
		{c6, "pass your config struct as a pointer"},
		{&c6, "user config must be a struct"},
		{c7, "pass your config struct as a pointer"},
		{&c7, "user config must be a struct"},
		{c8, "pass your config struct as a pointer"},
		{&c8, "user config must be a struct"},
		{c9, "pass your config struct as a pointer"},
		{&c9, "user config must be a struct"},
	}

	// Act & Assert
	for i, c := range cases {
		t.Log("Test case:", i)
		test_getStructInfo_error(t, c)
	}
}

func test_getStructInfo_error(t *testing.T, testCase testCaseGetStructInfoError) {
	// Arrange
	cr := &configReader{}

	// Act
	_, err := cr.getStructInfo(testCase.config, "", "")

	// Assert
	require.Equal(t, testCase.err, err.Error())
}

func Test_getStructInfo_ignoreNotSettable_cases(t *testing.T) {
	// Arrange
	c1 := struct{ field int }{}
	c2 := struct{ int }{}
	cases := []interface{}{&c1, &c2}

	// Act & Assert
	for i, c := range cases {
		t.Log("Test case:", i)
		test_getStructInfo_ignoreNotSettable(t, c)
	}
}

func test_getStructInfo_ignoreNotSettable(t *testing.T, config interface{}) {
	// Arrange
	cr := &configReader{}

	// Act
	si, err := cr.getStructInfo(config, "", "")

	// Assert
	require.Nil(t, err)
	require.Empty(t, si)
}

type testCaseGetTagDataSuccess struct {
	data         string
	field        reflect.StructField
	expName      string
	expRequired  bool
	expAppend    bool
	expUseParser bool
}

func Test_getTagData_success_cases(t *testing.T) {
	// Arrange
	str := struct{ Field int }{}
	field := reflect.TypeOf(str).Field(0)
	cases := []testCaseGetTagDataSuccess{
		{"", field, "Field", false, false, false},
		{"env_1", field, "env_1", false, false, false},
		{"env_1,required", field, "env_1", true, false, false},
		{"env_1,append", field, "env_1", false, true, false},
		{"env_1,required,append", field, "env_1", true, true, false},
		{"env_1,append,required,useparser", field, "env_1", true, true, true},
		{"append,required", field, "append", true, false, false},
		{"required,required", field, "required", true, false, false},
		{"append,append", field, "append", false, true, false},
		{"required,append", field, "required", false, true, false},
		{"required,required,append", field, "required", true, true, false},
		{"append,required,append", field, "append", true, true, false},
		{"useparser,useparser", field, "useparser", false, false, true},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data)
		test_getTagData_success_cases(t, c)
	}
}

func test_getTagData_success_cases(t *testing.T, testCase testCaseGetTagDataSuccess) {
	// Arrange
	cr := &configReader{}

	// Act
	name, required, append, useParser, err := cr.getTagData(testCase.data, testCase.field)

	// Assert
	require.Nil(t, err)
	require.Equal(t, testCase.expName, name)
	require.Equal(t, testCase.expRequired, required)
	require.Equal(t, testCase.expAppend, append)
	require.Equal(t, testCase.expUseParser, useParser)
}

type testCaseGetTagDataError struct {
	data  string
	field reflect.StructField
	err   string
}

func Test_getTagData_error_cases(t *testing.T) {
	// Arrange
	str := struct{ Field int }{}
	field := reflect.TypeOf(str).Field(0)
	cases := []testCaseGetTagDataError{
		{" ", field, "env tag is empty for field Field"},
		{"e.nv", field, "env tag contains invalid characters for field Field"},
	}

	// Act & Assert
	for _, c := range cases {
		t.Log("Test case:", c.data)
		test_getTagData_error_cases(t, c)
	}
}

func test_getTagData_error_cases(t *testing.T, testCase testCaseGetTagDataError) {
	// Arrange
	cr := &configReader{}

	// Act
	_, _, _, _, err := cr.getTagData(testCase.data, testCase.field)

	// Assert
	require.Equal(t, testCase.err, err.Error())
}

type testCaseReadConfigString struct {
	data     string
	ft       formatType
	expName  string
	expValue string
}

func Test_readConfigString_success_cases(t *testing.T) {
	// Arrange
	cases := []testCaseReadConfigString{
		{"key=value", FtEnv, "key", "value"},
		{"k ey=value", FtIni, "k ey", "value"},
		{"{\"key\": \"value\"}", FtJson, "key", "value"},
	}

	// Act & Assert
	for i, c := range cases {
		t.Log("Test case:", i)
		test_readConfigString_success(t, c)
	}
}

func test_readConfigString_success(t *testing.T, testCase testCaseReadConfigString) {
	// Arrange
	cr := &configReader{}
	it := intermediateTree{}
	si := []structInfo{{keyName: testCase.expName}}
	source := configSource{value: testCase.data, ft: testCase.ft}

	// Act
	err := cr.readConfigString(source, it, si, 0)

	// Assert
	require.Nil(t, err)
	require.Equal(t, testCase.expValue, it[testCase.expName][0].value)
}
