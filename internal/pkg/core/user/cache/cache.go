package cache

import (
	"gitlab.ozon.dev/DenisAleksandrovichM/masterclass-2/internal/pkg/core/user/models"
)

type Interface interface {
	Add(user models.User) error
	Read(login string) (models.User, error)
	Delete(login string) error
	List() []models.User
	Update(user models.User) error
}
