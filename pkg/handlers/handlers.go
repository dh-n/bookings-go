// Package handlers
package handlers

import (
	"net/http"

	"github.com/dh-n/bookings/pkg/config"
	"github.com/dh-n/bookings/pkg/models"
	"github.com/dh-n/bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		a,
	}
}

func SetNewRepoForHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remoteIp", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["data"] = "Hello there this is test data. Awesome"

	remoteIP := m.App.Session.GetString(r.Context(), "remoteIp")
	stringMap["remoteIp"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
