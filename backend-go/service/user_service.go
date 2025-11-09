package service

import (
	"backend-go/models"
	"backend-go/repository"
)

type UserService interface {
	GetFirstUser() (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetFirstUser() (*models.User, error) {
	return s.repo.First()
}
