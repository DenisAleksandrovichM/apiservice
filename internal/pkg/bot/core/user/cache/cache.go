package cache

import (
	"context"
	"gitlab.ozon.dev/DenisAleksandrovichM/homework-1/internal/pkg/bot/core/user/models"
)

type Interface interface {
	Add(ctx context.Context, user models.User) (models.User, error)
	Read(ctx context.Context, login string) (models.User, error)
	Delete(ctx context.Context, login string) (models.User, error)
	List(ctx context.Context, queryParams map[string]interface{}) ([]models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
}
