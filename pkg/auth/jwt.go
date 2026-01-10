package auth

import (
	"encoding/json"
	"shopify/internal/db/sqlc"
	"shopify/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
}

type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     int32  `json:"role"`
}

func NewJWTService() TokenService {
	return &jwtService{}
}

var (
	jwtSecret     = []byte(utils.GetEnv("JWT_SECRET", "jwt-secret-govn"))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", "12345678901234567890123456789012"))
)

const AccessTokenTTL = 15 * time.Minute

func (js *jwtService) GenerateAccessToken(user sqlc.User) (string, error) {
	payload := &EncryptedPayload{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     user.UserLevel,
	}

	// claims := &Claims{

	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		ID:        uuid.NewString(),
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
	// 		IssuedAt:  jwt.NewNumericDate(time.Now()),
	// 		Issuer:    "govn",
	// 	},
	// }

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encrypted, err := utils.EncryptAES(rawData, jwtEncryptKey)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"data": encrypted,
		"jti":  uuid.NewString(),
		"exp":  time.Now().Add(AccessTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "govn",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (js *jwtService) GenerateRefreshToken() {

}
