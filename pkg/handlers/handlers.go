package handlers

import (
	"github.com/fouched/go-bookings/pkg/config"
	"github.com/fouched/go-bookings/pkg/models"
	"github.com/fouched/go-bookings/pkg/render"
	"net/http"
)

// Repo the repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.DisplayTemplate(w, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.DisplayTemplate(w, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.DisplayTemplate(w, "make-reservation.page.gohtml", &models.TemplateData{})
}

// Generals renders the generals room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.DisplayTemplate(w, "generals.page.gohtml", &models.TemplateData{})
}

// Majors renders the majors room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.DisplayTemplate(w, "majors.page.gohtml", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.DisplayTemplate(w, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Posted to search availability"))
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.DisplayTemplate(w, "contact.page.gohtml", &models.TemplateData{})
}
