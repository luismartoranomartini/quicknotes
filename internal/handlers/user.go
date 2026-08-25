package handlers

import (
	"fmt"
	"net/http"
	"quicknotes/internal/repositories"
	"quicknotes/utils"
	"regexp"
	"strings"
)

type userHandler struct {
	repo repositories.UserRepository
}

func NewUserHandler(repo repositories.UserRepository) *userHandler {
	return &userHandler{repo: repo}
}

func (uh *userHandler) SignupForm(w http.ResponseWriter, r *http.Request) error {
	return render(w, http.StatusOK, "user-signup.html", nil)
}

func (uh *userHandler) Signup(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	data := newUserRequest(email, password)

	if strings.TrimSpace(data.Password) == "" {
		data.AddFieldError("password", "Senha é obrigatória")
	}

	if len(strings.TrimSpace(data.Password)) < 6 {
		data.AddFieldError("password", "Senha precisa ter no mínimo 6 caracteres")
	}

	if !isEmailValid(data.Email) || strings.TrimSpace(data.Password) == "" {
		data.AddFieldError("email", "Email é inválido")
	}

	if !data.Valid() {
		return render(w, http.StatusUnprocessableEntity, "user-signup.html", data)
	}

	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return err
	}

	hasToken := utils.GenerateTokenKey()

	user, token, err := uh.repo.Create(r.Context(), data.Email, hash, hasToken)
	if err == repositories.ErrDuplicateEmail {
		data.AddFieldError("email", "Email já existe")
		return render(w, http.StatusUnprocessableEntity, "user-signup.html", data)
	}
	if err != nil {
		return err
	}
	fmt.Println("Usuário criado:", user.ID)

	return render(w, http.StatusOK, "user-signup-success.html", token)
}

func (uh *userHandler) Confirm(w http.ResponseWriter, r *http.Request) error {
	token := r.PathValue("token")
	err := uh.repo.ConfirmUserByToken(r.Context(), token)
	msg := "Seu cadastro foi confirmado. Pode fazer o login"
	if err != nil {
		msg = "Esse cadastro já foi confirmado ou o token é inválido"
	}
	return render(w, http.StatusOK, "user-confirm.html", msg)
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
