package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"gin_project/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

var JWT_SECRET = "eUbP9shywUygMx7u"

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type JWTOutput struct {
	Token string `json:"token"`
	Expires time.Time `json:"expires"`
}
type AuthHandler struct {
	collection *mongo.Collection
	ctx context.Context
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":err.Error(),
		})
		return
	}

	h := sha256.New()
	h.Write([]byte(user.Password))
	b := h.Sum(nil)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"username":user.Username,
		"password":base64.StdEncoding.EncodeToString(b),
	})
	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":"Invalid username or password",
		})
		return
	}

	if user.Username != "admin" || user.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
	}
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	jwtOutput := JWTOutput{
		Token: tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)
}

func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":err.Error(),
		})
	}
	if tkn == nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":"Invalid token",
		})
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Token is not expired yet",
		})
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	jwtOutput := JWTOutput{
		Token: tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx: ctx,
	}
}