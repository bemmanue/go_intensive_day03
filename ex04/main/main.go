package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"go_intensive/go_intensive_day03/ex04/db"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthToken struct {
	TokenType string `json:"token_type"`
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

type AuthTokenClaim struct {
	*jwt.StandardClaims
	User
}

type Data struct {
	Name   string
	Places []db.Place
}

type ErrorMsg struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/api/get_token", handlerGetToken)
	http.HandleFunc("/api/recommend", handler)
	log.Fatalln(http.ListenAndServe("localhost:8888", nil))
}

func handlerGetToken(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	expiresAt := time.Now().Add(time.Minute * 1).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &AuthTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		User{user.Username, user.Password},
	}

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthToken{
		Token:     tokenString,
		TokenType: "Bearer",
		ExpiresIn: expiresAt,
	})
}

func handler(w http.ResponseWriter, r *http.Request) {

	bearer := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(bearer, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMsg{Message: err.Error()})
		return
	}
	if !token.Valid {
		json.NewEncoder(w).Encode(ErrorMsg{Message: "Invalid authorization token"})
	}

	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	lat, err := strconv.ParseFloat(values["lat"][0], 64)
	if err != nil {
		log.Fatalln(err)
	}

	lon, err := strconv.ParseFloat(values["lon"][0], 64)
	if err != nil {
		log.Fatalln(err)
	}

	places, err := db.GetPlaces(lat, lon)
	if err != nil {
		log.Println(err)
	}

	data := Data{
		Name:   "Recommendation",
		Places: places,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalln(err)
	}
}
