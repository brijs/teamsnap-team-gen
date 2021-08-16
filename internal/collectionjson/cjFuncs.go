package cj

import (
	"encoding/json"
	"log"
)

func ReadCollectionJson(inputData interface{}) (cj CollectionJsonType, err error) {
	var buf []byte

	switch inputData.(type) {
	case string:
		buf = []byte(inputData.(string))
	case []byte:
		buf = inputData.([]byte)
	default:
		log.Fatal("Unsupported Collection JSON data encounters")
	}

	err = json.Unmarshal(buf, &cj)

	return
}
