package controllers

import (
	"alkitab/entitys"
	"alkitab/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/jwt"
	"github.com/olebedev/emitter"
)

type UserController struct {
	s           *services.UserService
	mailService *services.MailService
	EmailJobs   chan services.EmailJob
}

func NewUserController(s *services.UserService,
	mailService *services.MailService) *UserController {
	return &UserController{s: s,
		mailService: mailService}
}

var SharedKey = []byte("sercrethatmaycontainch@r$32chars")

type TokenClaims struct {
	TokenClaims string `json:"tokenClaims"`
	UserID      int    `json:"userID"`
}

func (c *UserController) SignUpAddUser(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		Message string `json:"message"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	e := &emitter.Emitter{}

	e.On("Signup", func(event *emitter.Event) {
		data := event.Args

		username := data[0].(string)
		email := data[1].(string)
		password := data[2].(string)

		hashPw, _ := services.HashPassword(password)
		uuid := uuid.New()
		meUuid := uuid.String()

		user := entitys.UsersLetstalk{
			Username: username,
			Email:    email,
			Password: hashPw,
		}
		c.s.SignUpAddUser(user)

		c.s.StoreCodeVerif(username, meUuid)

		c.EmailJobs <- services.EmailJob{
			To:      email,
			Subject: "Confirmation Verification",
			Body:    "Klik link untuk verifikasi: http://localhost:8081/verif/" + meUuid,
		}

		resp := Response{Message: "success create account, please check your email"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	e.Emit("Signup", req.Username, req.Email, req.Password)
}
func (c *UserController) SigninUser(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Message string `json:"message"`
	}

	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user := entitys.UsersLetstalk{
		Username: req.Username,
		Password: req.Password,
	}

	userID := c.s.FindUserID(req.Username)

	fmt.Println("user id is", userID)
	myClaims := TokenClaims{
		TokenClaims: req.Username,
		UserID:      userID,
	}

	token, err := jwt.Sign(jwt.HS256, SharedKey, myClaims, jwt.MaxAge(24*10*time.Hour))
	if err != nil {
		panic(err)
	}

	fmt.Println(req.Username, userID)
	login := c.s.SigninUser(user)
	if login == true {
		w.Header().Set("Authorization", "Bearer "+string(token))
		w.Header().Set("Content-Type", "application/json")
		fmt.Print(token)
		json.NewEncoder(w).Encode("success login")
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Failed to login")
	}
}

func (c *UserController) Verification(w http.ResponseWriter, r *http.Request) {
	const prefix = "/verif/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}

	pathVerif := strings.TrimPrefix(r.URL.Path, prefix)
	pathVerif = strings.TrimSuffix(pathVerif, "/")

	c.s.VerifyCode(pathVerif)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("success to create account")
}
