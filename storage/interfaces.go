// Filename: storage/interfaces.go
package storage

import "github.com/mrjvadi/BackendPanelVpn/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
}
