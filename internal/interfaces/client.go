package interfaces

import "github.com/e-faizov/GophKeeper/internal/models"

type Requests interface {
	Registration(login, pass string) error
	Auth(login, pass string) error
	NewSecret(s models.Secret) error
	EditSecret(s models.Secret) error
	GetSecretsList() ([]models.Secret, error)
	GetSecret(uid string, version int) (models.Secret, error)
	RemoveSecret(uid string, version int) error
}
