package models

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // used interface because  we dont now the data type that is
	CSRFToken string
	Flash     string // any type of rapid message
	Warning   string
	Error     string
}
