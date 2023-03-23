package logic

import (
	"context"

	"github.com/google/uuid"

	"github.com/e-faizov/GophKeeper/internal/interfaces"
	"github.com/e-faizov/GophKeeper/internal/models"
)

type CryptoLogicImpl struct {
	Store interfaces.CryptoStore
}

func (p *CryptoLogicImpl) NewSecret(ctx context.Context, userID string, s models.Secret) error {
	uid := uuid.New()
	return p.Store.NewSecret(ctx, userID, uid.String(), s)
}

func (p *CryptoLogicImpl) EditSecret(ctx context.Context, userID string, s models.Secret) error {
	return p.Store.EditSecret(ctx, userID, s)
}

func (p *CryptoLogicImpl) RemoveSecret(ctx context.Context, userID string, s models.Secret) error {
	return p.Store.RemoveSecret(ctx, userID, s)
}

func (p *CryptoLogicImpl) GetSecret(ctx context.Context, userID string, s models.Secret) (models.Secret, error) {
	return p.Store.GetSecret(ctx, userID, s)
}

func (p *CryptoLogicImpl) GetSecretsList(ctx context.Context, userID string) ([]models.Secret, error) {
	return p.Store.GetSecretsList(ctx, userID)
}
