package service

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	
    "Auth/internal/config"
    "Auth/internal/database"
    "Auth/internal/entity"
	"Auth/internal/DTO"
)

func Login(req DTO.LoginRequest) (string, error) {

    var user entity.User
    if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
        return "", errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return "", errors.New("invalid credentials")
    }
    return generateJWT(user.ID)
}

// generateJWT signs a new token for a given user ID
func generateJWT(userID uuid.UUID) (string, error) {
    cfg := config.Load()
    claims := jwt.MapClaims{
        "sub": userID.String(),
        "exp": time.Now().Add(72 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.JWTSecret))
}




func ExtractUserIDFromToken(tokenStr string) (uuid.UUID, error) {
    cfg := config.Load()
    token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(cfg.JWTSecret), nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userIDStr := claims["sub"].(string)
        return uuid.Parse(userIDStr)
    }

    return uuid.Nil, errors.New("Invalid Token, Please Log In Again")
}
