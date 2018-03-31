package util

import (
	"encoding/json"
	"errors"
	"fmt"
)

func Deserialize(serializedBytes []byte, object interface{}) error {
	if len(serializedBytes) == 0 {
		return nil
	}
	err := json.Unmarshal(serializedBytes, object)
	if err != nil {
		panic(fmt.Sprintf("Error decoding : %s", err))
	}
	return err
}

func Serialize(object interface{}) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error encoding : %s", err))
	}
	return data, nil
}
