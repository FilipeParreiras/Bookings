package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/FilipeParreiras/Bookings/internal/config"
	"github.com/FilipeParreiras/Bookings/internal/forms"
	"github.com/FilipeParreiras/Bookings/internal/models"
	"github.com/FilipeParreiras/Bookings/internal/render"
	"log"
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

// Handler function
// An handler function needs to have these two arguments
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// gets remote IP Address and stores it in the session
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	// Render the template
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
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
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	log.Println("Get")
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	// Render the template
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostReservation handles the posting
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	log.Println("Post")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	// Render the template
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	// Render the template
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	// Render the template
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// Get data from user input
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

// jsonResponse is only used in AvailabilityJSON
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles requests for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	// Parameters to send
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	// Creates the JSON
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))
	// What kind of response am i sending
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	// Render the template
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	// Render the template
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
