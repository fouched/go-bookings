package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fouched/go-bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	{"about", "/", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-existent", "/green/eggs/and/ham", "GET", http.StatusNotFound},

	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new reservations", "/admin/reservations-new", "GET", http.StatusOK},
	{"all reservations", "/admin/reservations-all", "GET", http.StatusOK},
	{"show reservation", "/admin/reservations/new/1/show", "GET", http.StatusOK},

	//{"post-search-availability", "/search-availability", "POST", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-02"},
	//}, http.StatusOK},
	//{"post-search-availability-json", "/search-availability-json", "POST", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-02"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	// defer - only gets executed after the function completes
	defer testServer.Close()

	// this is a table test
	for _, e := range theTests {
		//log.Println("Running:" + e.name + " on " + e.method + " request, with url " + testServer.URL + e.url)
		resp, err := testServer.Client().Get(testServer.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
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

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation not in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	// session data
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	// form data
	reqBody := "first_name=Fouche"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=du%20Preez")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=fouche@dupreez.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0828567442")
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	// cleaner way of building form data
	//postedData := url.Values{}
	//postedData.Add("first_name", "Fouche")
	//postedData.Add("last_name", "du Preez")
	//postedData.Add("email", "fouche@dupreez.com")
	//postedData.Add("phone", "0828567442")
	//req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing form body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code for missing post body: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case where reservation not in session
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid form data
	reqBody = "first_name=F"
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code for invalid form data: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for reservation database insert failure
	reqBody = "first_name=Fouche"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=du%20Preez")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=fouche@dupreez.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0828567442")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	reservation.RoomID = 2
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code for failed reservation database insert: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for room restriction database insert failure
	reqBody = "first_name=Fouche"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=du%20Preez")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=fouche@dupreez.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=0828567442")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	reservation.RoomID = 1000
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code for failed restriction database insert: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// rooms are not available
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
} {
	{
		"valid-credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
