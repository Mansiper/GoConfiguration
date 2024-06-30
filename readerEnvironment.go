package configuration

import (
	"os"
)

func (cr *configReader) readEnvironment(it intermediateTree, si []structInfo, sourceId int) {
	for _, s := range si {
		if val, ok := os.LookupEnv(s.keyName); ok {
			if _, ok := it[s.keyName]; !ok {
				it[s.keyName] = []intermediateData{{source: sourceId, value: val, valueType: vtAny}}
			} else {
				it[s.keyName] = append(it[s.keyName], intermediateData{source: sourceId, value: val, valueType: vtAny})
			}
		}
	}
}
