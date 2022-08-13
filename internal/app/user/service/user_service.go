package service

import (
	"errors"

	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Create(user domain.User) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	FindByID(ID int) (domain.User, error)
	Update(user domain.User) (domain.User, error)
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

	return user, nil
}

func (s *userService) Login(user domain.LoginParam) (domain.User, error) {
	newUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		return newUser, err
	}

	if newUser.ID == 0 {
		return newUser, errors.New("No User found in that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(newUser.PasswordHash), []byte(user.Password))
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *userService) IsEmailAvailable(input domain.CheckEmailInput) (bool, error) {
	email := input.Email
	user, err := s.userRepo.FindByEmail(email)

	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}
	return false, nil
}

func (s *userService) SaveAvatar(ID int, fileLocation string) (domain.User, error) {
	user, err := s.userRepo.FindByID(ID)

	if err != nil {
		return user, err
	}
	user.AvatarFileName = fileLocation
	updatedUser, err := s.userRepo.Update(user)

	if err != nil {
		return updatedUser, err
	}
	return updatedUser, nil
}
