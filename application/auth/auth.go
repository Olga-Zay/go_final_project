package auth

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go_final_project/config"
	"go_final_project/service/model"
	"net/http"
)

type Auth struct {
	config      *config.Config
	addressAuth map[string]bool
}

func NewAuth(config *config.Config) Auth {
	return Auth{
		config: config,
		addressAuth: map[string]bool{
			"/api/task":      true,
			"/api/tasks":     true,
			"/api/task/done": true,
		},
	}
}

func (a Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.addressAuth[r.RequestURI] {
			next.ServeHTTP(w, r)
			return
		}
		// смотрим наличие пароля
		if len(a.config.Pass) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			if !isJWTValid(a.config.Pass, jwt) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (a Auth) SingIn(w http.ResponseWriter, r *http.Request) {
	request, err := prepareSingInRequest(r)
	if err != nil {
		prepareResponse(w, model.SingInResponseWithError{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	if request.Password == a.config.Pass {
		secret := []byte(a.config.Pass)

		jwtToken := jwt.New(jwt.SigningMethodHS256)

		signedToken, err := jwtToken.SignedString(secret)
		if err != nil {
			prepareResponse(w, model.SingInResponseWithError{Error: err.Error()}, http.StatusInternalServerError)
			return
		}

		prepareResponse(w, model.SingInResponse{Token: signedToken}, http.StatusOK)
		return
	}

	prepareResponse(w, model.SingInResponseWithError{Error: "Неправильный пароль"}, http.StatusUnauthorized)
}

func prepareSingInRequest(r *http.Request) (model.SingInRequest, error) {
	var singInRequest model.SingInRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&singInRequest); err != nil {
		return model.SingInRequest{}, fmt.Errorf("ошибка десериализации JSON: %s", err.Error())
	}

	return singInRequest, nil
}

func prepareResponse(w http.ResponseWriter, response any, httpStatus int) {
	encoderErr := json.NewEncoder(w).Encode(&response)
	if encoderErr != nil {
		http.Error(w, encoderErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatus)
}

func isJWTValid(pass string, token string) bool {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(pass), nil
	})
	if err != nil {
		fmt.Printf("Failed to parse token: %s\n", err)
		return false
	}
	if jwtToken.Valid {
		return true
	}

	return false
}
