package security

import (
	"bacancy/go-boiler-plate/app/config"
	"bacancy/go-boiler-plate/app/models/user"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"golang.org/x/crypto/scrypt"
)

const PW_SALT_BYTES = 32
const PW_HASH_BYTES = 64

func Login(c *gin.Context) {

	email := c.PostForm("email")
	password := c.PostForm("password")

	u, found, err := user.GetUserByEmail(email)
	if found == false {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something went wrong", "detail": err})
		return
	}

	hash, err := Hash(password, u.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something went wrong", "detail": err})
		return
	}

	if u.Password != hash {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Something went wrong", "detail": "Invalid password"})
		return
	}

	token, err := CreateToken(u.ID, u.Name, u.Email, u.Admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something went wrong", "detail": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "email": u.Email, "admin": u.Admin, "Name": u.Name, "Lastname": u.LastName, "id": u.ID})
}

func LoginAdmin(c *gin.Context) {

	email := c.PostForm("email")
	password := c.PostForm("password")

	u, found, err := user.GetUserByEmail(email)
	if found == false {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something went wrong", "detail": err})
		return
	}

	hash, err := Hash(password, u.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something went wrong", "detail": err})
		return
	}

	if u.Password != hash {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	if u.Admin == false {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	token, err := CreateToken(u.ID, u.Name, u.Email, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Something went wrong", "detail": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "email": u.Email, "admin": u.Admin, "Name": u.Name, "Lastname": u.LastName, "id": u.ID})
}

func GetSalt() (string, error) {
	salt := make([]byte, PW_SALT_BYTES)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", salt), nil
}

func Hash(password string, salt string) (string, error) {

	hash, err := scrypt.Key([]byte(password), []byte(salt), 1<<14, 8, 1, PW_HASH_BYTES)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash), nil
}

type JWTCustomClaims struct {
	Id    uint
	Email string
	Admin bool
	jwt.StandardClaims
}

type JWTToken struct {
	Id    uint
	Name  string
	Email string
	Admin bool
}

type JWTStripeCustomClaims struct {
	UserID  uint
	EventID uint
	jwt.StandardClaims
}

type JWTStripeToken struct {
	UserID  uint
	EventID uint
}

func CreateToken(id uint, name string, email string, admin bool) (string, error) {
	claims := JWTCustomClaims{
		id,
		email,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			Issuer:    name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.GetConfig().JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetTokenData(tokenString string) (JWTToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWT_SECRET), nil
	})
	if err != nil {
		return JWTToken{}, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return JWTToken{
			Id:    claims.Id,
			Name:  claims.StandardClaims.Issuer,
			Email: claims.Email,
			Admin: claims.Admin,
		}, nil
	} else {
		return JWTToken{}, errors.New("Invalid token")
	}
}
