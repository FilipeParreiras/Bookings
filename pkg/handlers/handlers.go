package handlers

import (
	"net/http"

	"github.com/FilipeParreiras/Bookings/pkg/config"
	"github.com/FilipeParreiras/Bookings/pkg/models"
	"github.com/FilipeParreiras/Bookings/pkg/render"
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

// Handler function
// An handler function needs to have these two arguments
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// gets remote IP Address and stores it in the session
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	// Render the template
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
// gives acess to anything inside repository
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// Creates stringMap that contains key-values to use in HTML code
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	// Gets the remote IP from session and stores it in stringmap
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// Render the template
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
