package utils

import (
	"encoding/json"

	"github.com/mattn/go-jsonpointer"
)

func JpToString(jsonBytes []byte, jp string) (string, error) {
	var obj interface{}
	json.Unmarshal(jsonBytes, &obj)

	v, err := jsonpointer.Get(obj, jp)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	str := string(b)
	trimmed := str[1 : len(str)-1] // jsonの""を除去

	return trimmed, nil
}
