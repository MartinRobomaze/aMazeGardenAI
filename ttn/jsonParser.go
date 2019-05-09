package ttn

import (
	"encoding/json"
)

type JsonParser struct {
	Json []byte
	JsonStructure interface{}
}

func (jsonParser JsonParser) JsonToData(v interface{}) error {
	err := json.Unmarshal(jsonParser.Json, &jsonParser.JsonStructure)

	v = jsonParser.JsonStructure

	if (err != nil) {
		return err
	}

	return nil
}
