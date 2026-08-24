package handlers

import (
	"fmt"
	"net/http"
	"quicknotes/internal/repositories"
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

	user, err := uh.repo.Create(r.Context(), data.Email, data.Password)
	if err == repositories.ErrDuplicateEmail {
		data.AddFieldError("email", "Email já existe")
		return render(w, http.StatusUnprocessableEntity, "user-signup.html", data)
	}
	if err != nil {
		return err
	}
	fmt.Println("Usuário criado:", user.ID)

	return render(w, http.StatusOK, "user-signup-success.html", nil)
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
