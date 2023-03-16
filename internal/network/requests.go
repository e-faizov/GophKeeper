package network

import (
	"errors"

	"github.com/e-faizov/GophKeeper/internal/models"
)

var jwt string

var (
	unAuthErr = errors.New("unauthorized")
)

func Registration(login, pass string) error {
	jwt = "jwt"
	return nil
}

func Auth(login, pass string) error {
	jwt = "jwt"
	return errors.New("wrong")
}

func NewSecret(s models.Secret) error {
	if len(jwt) == 0 {
		return unAuthErr
	}

	return nil
}

func EditSecret(s models.Secret) error {
	if len(jwt) == 0 {
		return unAuthErr
	}

	if s.Version == 0 {
		return errors.New("empty version")
	}
	if len(s.Uid) == 0 {
		return errors.New("empty uid")
	}
	return nil
}

func GetSecretsList() ([]models.Secret, error) {
	if len(jwt) == 0 {
		return nil, unAuthErr
	}

	return []models.Secret{
		{
			Data1:   "data1",
			Version: 3,
			Type:    1,
		},
		{
			Data1:   "data2",
			Version: 3,
			Type:    1,
		},
	}, nil
}

func GetSecret(uid string, version int) (models.Secret, error) {
	if len(jwt) == 0 {
		return models.Secret{}, unAuthErr
	}

	return models.Secret{
		Data1:   "data1",
		Data2:   "data2",
		Version: version,
		Uid:     uid,
		Type:    1,
	}, nil
}

func RemoveSecret(uid string, version int) error {
	if len(jwt) == 0 {
		return unAuthErr
	}

	return nil
}
