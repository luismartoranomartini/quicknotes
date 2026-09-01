package handlers

import (
	"net/http"
	render "quicknotes/internal/render"
)

type homeHandler struct {
	render *render.RenderTenplate
}

func NewHomeHandler(render *render.RenderTenplate) *homeHandler {
	return &homeHandler{render: render}
}

func (hh *homeHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	hh.render.RenderPage(w, r, http.StatusOK, "home.html", nil)
}
