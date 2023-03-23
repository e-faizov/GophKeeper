package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	"github.com/e-faizov/GophKeeper/internal/models"
)

// go binary encoder
func ConvertToGOB64(m models.User) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// go binary decoder
func ConvertFromGOB64(str string) (models.User, error) {
	var m models.User
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return models.User{}, err
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		return models.User{}, err
	}
	return m, nil
}
