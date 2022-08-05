package service

import (
	"errors"

	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Create(user domain.User) (domain.User, error)
}

type userService struct {
	userRepo User
}

func NewUserService(usr User) *userService {
	return &userService{
		userRepo: usr,
	}
}

func (s *userService) Register(input domain.UserParam) (domain.User, error) {
	hashedpassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	param := domain.User{
		Name:         input.Name,
		Occupation:   input.Occupation,
		Email:        input.Email,
		PasswordHash: string(hashedpassword),
		Role:         "user",
	}

	user, err := s.userRepo.Create(param)
	if err != nil {
		return user, errors.New("failed create user")
	}

	return user, err
}
