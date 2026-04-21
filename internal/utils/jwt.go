package utils

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
    secretKey []byte
    expiry    time.Duration
}

type Claims struct {
    UserID   int64  `json:"userId"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func NewJWTUtil(secretKey string, expiryHours int) *JWTUtil {
    return &JWTUtil{
        secretKey: []byte(secretKey),
        expiry:    time.Duration(expiryHours) * time.Hour,
    }
}

func (j *JWTUtil) GenerateToken(userID int64, username, role string) (string, error) {
    claims := &Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return j.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}