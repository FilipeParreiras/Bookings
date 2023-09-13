package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FilipeParreiras/Bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"make-res", "/make-reservation", "GET", http.StatusOK},

	//{"post-search-availability", "/search-availability", "Post", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-02"},
	//}, http.StatusOK},
	//{"post-search-availability-json", "/search-availability-json", "Post", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-02"},
	//}, http.StatusOK},
	{"make-reservation", "/make-reservation", "Post", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	request, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getConstext(request)
	request = request.WithContext(ctx)

	// responseRecorder simulates what we get from request response life cycle when someone passes a request in our
	//website and the responseWriter gives the response
	responseRecorder := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	// Cast reservation handler into handlerFunc
	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Reservation handler retured wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	request, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getConstext(request)
	request = request.WithContext(ctx)
	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler retured wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existing room
	request, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getConstext(request)
	request = request.WithContext(ctx)
	responseRecorder = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler retured wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	request, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getConstext(request)
	request = request.WithContext(ctx)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler retured wrong response code: got %d, wanted %d", responseRecorder.Code,
			http.StatusSeeOther)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// First case -> rooms are not available
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get context with session
	ctx := getConstext(req)
	req = req.WithContext(ctx)

	// send request header
	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	// make handler handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	//get response recorder
	responseRecorder := httptest.NewRecorder()

	// make request to handler
	handler.ServeHTTP(responseRecorder, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(responseRecorder.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}

func getConstext(request *http.Request) context.Context {
	ctx, err := session.Load(request.Context(), request.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
