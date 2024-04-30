package services

import (
	"encoding/base64"
	"net/http"
	"runners/models"
	"runners/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	usersRepository *repositories.UsersRepository
}

func NewUsersService(usersRepository *repositories.UsersRepository) *UsersService {
	return &UsersService{usersRepository: usersRepository}
}

func (us UsersService) Login(username string, password string) (string, *models.ResponseError) {
	if username == "" || password == "" {
		return "", &models.ResponseError{
			Message: "Invalid username or password",
			Status:  http.StatusBadRequest,
		}
	}

	id, resErr := us.usersRepository.LoginUser(username, password)
	if resErr != nil {
		return "", resErr
	}

	if id == "" {
		return "", &models.ResponseError{
			Message: "Login failed",
			Status:  http.StatusUnauthorized,
		}
	}

	accessToken, resErr := generateAccessToken(username)
	if resErr != nil {
		return "", resErr
	}

	us.usersRepository.SetAccessToken(accessToken, id)
	return accessToken, nil
}

func (us UsersService) Logout(accessToken string) *models.ResponseError {
	if accessToken == "" {
		return &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}
	return us.usersRepository.RemoveAccessToken(accessToken)
}

func (us UsersService) AuthorizeUser(accessToken string, expectedRoles []string) (bool, *models.ResponseError) {
	if accessToken == "" {
		return false, &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}

	role, resErr := us.usersRepository.GetUserRole(accessToken)
	if resErr != nil {
		return false, resErr
	}

	if role == "" {
		return false, &models.ResponseError{
			Message: "Failed to authorize user",
			Status:  http.StatusUnauthorized,
		}
	}

	for _, exexpectedRole := range expectedRoles {
		if exexpectedRole == role {
			return true, nil
		}
	}

	return false, nil
}

func generateAccessToken(username string) (string, *models.ResponseError) {
	hash, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.ResponseError{
			Message: "Failed to generate token",
			Status:  http.StatusInternalServerError,
		}
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}
