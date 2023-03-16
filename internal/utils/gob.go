package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

// go binary encoder
func ToGOB64(m interface{}) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// go binary decoder
func FromGOB64(str string) (interface{}, error) {
	var m interface{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
