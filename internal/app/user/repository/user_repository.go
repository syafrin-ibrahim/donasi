package repository

import (
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"gorm.io/gorm"
)

type userDBRepository struct {
	db *gorm.DB
}

func NewUserDBRepository(db *gorm.DB) *userDBRepository {
	return &userDBRepository{
		db: db,
	}
}

func (usr *userDBRepository) Create(user domain.User) (domain.User, error) {
	err := usr.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func (usr *userDBRepository) FindByEmail(email string) (domain.User, error) {
	var user domain.User
	err := usr.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}