package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"UserService/internal/entity"
	"UserService/internal/DTO"
	"UserService/internal/database"
)


func UpdateUser(req DTO.UpdateUserRequest) (string, error) {

	var user entity.User
	if err := database.DB.First(&user, "username = ?", req.Username).Error; err != nil {
		return "", errors.New("user not found")
	}

	if req.Username == ""  {
		return "", errors.New("username cannot be empty")
	}
	user.Username = req.Username

	if req.Password != "" || req.ConfirmPassword != "" || req.Password != req.ConfirmPassword {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return "", errors.New("failed to hash password")
		}
		user.PasswordHash = string(hashed)
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return "", errors.New("failed to update user")
	}

	return "User updated successfully", nil
}