package models

import "github.com/FilipeParreiras/Bookings/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // used interface because  we don't now the data type that is
	CSRFToken string
	Flash     string // any type of rapid message
	Warning   string
	Error     string
	Form      *forms.Form
}
