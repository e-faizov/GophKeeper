package interfaces

import (
	"context"

	"github.com/e-faizov/GophKeeper/internal/models"
)

type CryptoLogic interface {
	NewSecret(ctx context.Context, userID string, s models.Secret) error
	EditSecret(ctx context.Context, userID string, s models.Secret) error
	RemoveSecret(ctx context.Context, userID string, s models.Secret) error
	GetSecret(ctx context.Context, userID string, s models.Secret) (models.Secret, error)
	GetSecretsList(ctx context.Context, userID string) ([]models.Secret, error)
}

type CryptoStore interface {
	NewSecret(ctx context.Context, userID, uid string, s models.Secret) error
	GetSecret(ctx context.Context, userID string, s models.Secret) (models.Secret, error)
	EditSecret(ctx context.Context, userID string, s models.Secret) error
	RemoveSecret(ctx context.Context, userID string, s models.Secret) error
	GetSecretsList(ctx context.Context, userID string) ([]models.Secret, error)
}
