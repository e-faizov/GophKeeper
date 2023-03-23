package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/e-faizov/GophKeeper/internal/models"
)

var jwt string

var (
	unAuthErr = errors.New("unauthorized")
)

type HttpsRequests struct {
	Url string
}

func (h *HttpsRequests) Registration(login, pass string) error {
	usr := models.User{
		Login:    login,
		Password: pass,
	}

	data, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	resp, err := http.Post(h.Url+"/api/user/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusConflict {
		return errors.New("try another login")
	}

	if resp.StatusCode != http.StatusCreated || resp.Cookies() == nil {
		return unAuthErr
	}

	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == "jwt" {
			jwt = c.Value
			return nil
		}
	}

	return unAuthErr
}

func (h *HttpsRequests) Auth(login, pass string) error {
	usr := models.User{
		Login:    login,
		Password: pass,
	}

	data, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	resp, err := http.Post(h.Url+"/api/user/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK || resp.Cookies() == nil {
		return unAuthErr
	}

	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == "jwt" {
			jwt = c.Value
			fmt.Println(jwt)
			return nil
		}
	}

	return unAuthErr
}

func (h *HttpsRequests) NewSecret(s models.Secret) error {
	if len(jwt) == 0 {
		return unAuthErr
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	resp, err := postJsonData(h.Url+"/api/secrets/new", data)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New("secret does not created")
	}
	return nil
}

func (h *HttpsRequests) EditSecret(s models.Secret) error {
	if len(jwt) == 0 {
		return unAuthErr
	}

	if s.Version == 0 {
		return errors.New("empty version")
	}
	if len(s.Uid) == 0 {
		return errors.New("empty uid")
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	resp, err := postJsonData(h.Url+"/api/secrets/edit", data)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("secret does not changed")
	}
	return nil
}

func (h *HttpsRequests) GetSecretsList() ([]models.Secret, error) {
	if len(jwt) == 0 {
		return nil, unAuthErr
	}

	resp, err := getData(h.Url + "/api/secrets/getAll")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("wrong http status")
	}

	var res []models.Secret

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *HttpsRequests) GetSecret(uid string, version int) (models.Secret, error) {
	if len(jwt) == 0 {
		return models.Secret{}, unAuthErr
	}

	data, err := json.Marshal(models.Secret{
		Uid:     uid,
		Version: version,
	})
	if err != nil {
		return models.Secret{}, err
	}

	resp, err := postJsonData(h.Url+"/api/secrets/get", data)
	if err != nil {
		return models.Secret{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.Secret{}, errors.New("secret does not geted")
	}

	var res models.Secret
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return models.Secret{}, err
	}

	return res, nil
}

func (h *HttpsRequests) RemoveSecret(uid string, version int) error {
	if len(jwt) == 0 {
		return unAuthErr
	}

	data, err := json.Marshal(models.Secret{
		Uid:     uid,
		Version: version,
	})
	if err != nil {
		return err
	}

	resp, err := postJsonData(h.Url+"/api/secrets/remove", data)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("secret does not remove")
	}
	return nil
}

func postJsonData(url string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	return client.Do(req)
}

func getData(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	return client.Do(req)
}
