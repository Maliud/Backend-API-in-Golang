package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Maliud/Backend-API-in-Golang/config"
	"github.com/Maliud/Backend-API-in-Golang/types"
	"github.com/Maliud/Backend-API-in-Golang/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// kullanıcı isteğinden belirteci al

		tokenString := getTokenFromRequest(r)
		// JWT'yi doğrulayın
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("token doğrulanamadı:", err)
			permissionDenied(w)
			return
		}
		if !token.Valid {
			log.Println("Token Geçersiz")
			permissionDenied(w)
			return
		}
		// eğer kullanıcı kimliğini DB'den getirmemiz gerekiyorsa (token'dan ise)
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, _ := strconv.Atoi(str)
		u, err := store.GetUserByID(userID)

		if err != nil {
			log.Printf("fkullanıcıyı kimliğine göre alamadı:", err)
			permissionDenied(w)
			return
		}
		// userID 'deki bağlamı ayarla “userID"

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)

	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" {
		return tokenAuth
	}
	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("beklenmedik imzalama yöntemi:", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("izin reddedildi."))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userID
}