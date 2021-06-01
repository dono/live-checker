package utils

import (
	"encoding/json"

	"github.com/mattn/go-jsonpointer"
)

func JpToString(obj interface{}, jp string) (string, error) {
	v, err := jsonpointer.Get(obj, jp)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	str := string(b)

	if str[0:1] == `"` && str[len(str)-1:] == `"` {
		str = str[1 : len(str)-1] // jsonの""を除去
	}

	return str, nil
}

// func JpToString(jsonBytes []byte, jp string) (string, error) {
// 	var obj interface{}
// 	err := json.Unmarshal(jsonBytes, &obj)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	v, err := jsonpointer.Get(obj, jp)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	b, err := json.Marshal(v)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	str := string(b)
// 	trimmed := str[1 : len(str)-1] // jsonの""を除去
//
// 	return trimmed, nil
// }
