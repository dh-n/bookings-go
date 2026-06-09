// Package handlers
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dh-n/bookings/internal/config"
	"github.com/dh-n/bookings/internal/forms"
	"github.com/dh-n/bookings/internal/helpers"
	"github.com/dh-n/bookings/internal/models"
	"github.com/dh-n/bookings/internal/render"
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

	_ = render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["data"] = "Hello there this is test data. Awesome"

	m.App.Session.Put(r.Context(), "device", r.Host)

	remoteIP := m.App.Session.GetString(r.Context(), "remoteIp")
	stringMap["remoteIp"] = remoteIP

	device := m.App.Session.GetString(r.Context(), "device")
	stringMap["device"] = device

	_ = render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]any)
	data["reservation"] = emptyReservation
	_ = render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)

	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "phone", "email")

	form.MinLength("first_name", 3, r)

	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation
		_ = render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Request received"))
}

type jsonResponseData struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// PostAvailabilityJSON renders the json data
func (m *Repository) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	data, err := json.MarshalIndent(jsonResponseData{
		Ok:      true,
		Message: "Data received",
	}, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary renders the summary page of the reservation
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// data := make(map[string]any)
	// data["reservation"] = m.App.Session.Get(r.Context(), "reservation")

	// alternative way to do this using type assertion
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) // asserting the type
	if !ok {
		m.App.ErrorLog.Println("cannot get reservation from session")
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]any)
	data["reservation"] = reservation
	_ = render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
