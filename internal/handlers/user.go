package handlers

import (
	"fmt"
	"net/http"
	"quicknotes/internal/mailer"
	render "quicknotes/internal/render"
	"quicknotes/internal/repositories"
	"quicknotes/utils"
	"regexp"
	"strings"

	"github.com/alexedwards/scs/v2"
)

type userHandler struct {
	render  *render.RenderTenplate
	session *scs.SessionManager
	mail    mailer.MailService
	repo    repositories.UserRepository
}

// NewUserHandler -> função de construção
func NewUserHandler(render *render.RenderTenplate, session *scs.SessionManager, mail mailer.MailService, repo repositories.UserRepository) *userHandler {
	return &userHandler{render: render, session: session, mail: mail, repo: repo}
}

func (uh *userHandler) Me(w http.ResponseWriter, r *http.Request) error {

	fmt.Fprintf(w, "Dados do usuário")
	return nil
}

func (uh *userHandler) SigninForm(w http.ResponseWriter, r *http.Request) error {
	userID := uh.session.GetInt64(r.Context(), "userID")
	fmt.Println("USER ID:", userID)
	return uh.render.RenderPage(w, r, http.StatusOK, "user-signin.html", nil)
}

func (uh *userHandler) Signin(w http.ResponseWriter, r *http.Request) error {

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
	if !isEmailValid(data.Email) {
		data.AddFieldError("email", "Email é inválido")
	}

	if !data.Valid() {
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signin.html", data)
	}

	// Consultar o usuário pelo email
	user, err := uh.repo.FindByEmail(r.Context(), data.Email)
	if err != nil {
		data.AddFieldError("validation", "credenciais inválidas")
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signin.html", data)
	}

	//verificar se usuário está ativo
	if !user.Active.Bool {
		data.AddFieldError("validation", "usuário não confirmou o cadastro")
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signin.html", data)
	}

	// Validação da senha
	if !utils.ValidatePassword(data.Password, user.Password.String) {
		data.AddFieldError("validation", "credenciais inválidas")
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signin.html", data)
	}
	//Renew Token
	err = uh.session.RenewToken(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return err
	}
	// objeto de sessão
	// armazena o ID na sessão
	uh.session.Put(r.Context(), "userID", user.ID.Int.Int64())
	uh.session.Put(r.Context(), "userEmail", user.Email.String)

	http.Redirect(w, r, "/note", http.StatusSeeOther)
	return nil
}

func (uh *userHandler) Confirm(w http.ResponseWriter, r *http.Request) error {
	token := r.PathValue("token")
	err := uh.repo.ConfirmUserByToken(r.Context(), token)
	msg := "Seu cadastro foi confirmado. Pode fazer o login"
	if err != nil {
		msg = "Esse cadastro já foi confirmado ou o token é inválido"
	}
	return uh.render.RenderPage(w, r, http.StatusOK, "user-confirm.html", msg)
}

func (uh *userHandler) SignupForm(w http.ResponseWriter, r *http.Request) error {
	return uh.render.RenderPage(w, r, http.StatusOK, "user-signup.html", nil)
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
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signup.html", data)
	}

	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return err
	}

	hasToken := utils.GenerateTokenKey()

	_, token, err := uh.repo.Create(r.Context(), data.Email, hash, hasToken)
	if err == repositories.ErrDuplicateEmail {
		data.AddFieldError("email", "Email já existe")
		return uh.render.RenderPage(w, r, http.StatusUnprocessableEntity, "user-signup.html", data)
	}
	if err != nil {
		return err
	}
	// fmt.Println("Usuário criado:", user.ID)

	body, err := uh.render.RenderMailBody("confirmation.html", token)
	if err != nil {
		return err
	}
	// enviar email de confirmaçãod de cadastro
	uh.mail.Send(mailer.MailMessage{
		To:      []string{data.Email},
		Subject: "Confirmação de Cadastro",
		IsHTML:  true,
		Body:    body,
	})

	return uh.render.RenderPage(w, r, http.StatusOK, "user-signup-success.html", token)
}
func (uh *userHandler) Signout(w http.ResponseWriter, r *http.Request) error {
	//Renew Token
	err := uh.session.RenewToken(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	uh.session.Remove(r.Context(), "userID")
	http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
