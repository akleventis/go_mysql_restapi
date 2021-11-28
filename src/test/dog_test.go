package test

import (
	"bytes"
	handlers "go_mysql/src/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func TestGetDogs(t *testing.T) {
	// create new mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// use mock database
	app := &handlers.App{
		DB: db,
	}

	// create new request to /dogs endpoint, nil => body io.Reader
	req, err := http.NewRequest("GET", "/dogs", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Will store response received from /dogs endpoint => pointer to ResonseRecorder struct
	recorder := httptest.NewRecorder()

	// insert mock data
	rows := sqlmock.NewRows([]string{"id", "name", "age", "color", "gender", "breed", "weight"}).
		AddRow(Dog1Mock.ID, Dog1Mock.Name, Dog1Mock.Age, Dog1Mock.Color, Dog1Mock.Gender, Dog1Mock.Breed, Dog1Mock.Weight).
		AddRow(Dog2Mock.ID, Dog2Mock.Name, Dog2Mock.Age, Dog2Mock.Color, Dog2Mock.Gender, Dog2Mock.Breed, Dog2Mock.Weight)
	mock.ExpectQuery("^SELECT (.+) FROM dogs$").WillReturnRows(rows)

	// execute request with mock database filled with....mock data
	handler := http.HandlerFunc(app.GetDogs)
	// hit endpoint with recorder and request
	handler.ServeHTTP(recorder, req)

	// test recorder status code
	if recorder.Code != http.StatusOK {
		t.Errorf("getDogs return wrong status code: got %v but want %v", recorder.Code, http.StatusOK)
	}

	// test return body. for some
	//"Encode writes the JSON encoding of v to the stream, followed by a newline character." => add new line to expected
	expected := `[{"id":1,"name":"Dog1","age":"5","color":"Brown","gender":"Female","breed":"Husky","weight":"34"},{"id":2,"name":"Dog2","age":"3","color":"Grey","gender":"Male","breed":"Boxer","weight":"4"}]
`
	if recorder.Body.String() != expected {
		t.Errorf("Wrong return body: got %s but want %s", recorder.Body.String(), expected)
	}
}
func TestGetDogById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	app := &handlers.App{
		DB: db,
	}

	req, err := http.NewRequest("GET", "/dogs", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "2"})

	recorder := httptest.NewRecorder()

	rows := sqlmock.NewRows([]string{"id", "name", "age", "color", "gender", "breed", "weight"}).
		AddRow(Dog2Mock.ID, Dog2Mock.Name, Dog2Mock.Age, Dog2Mock.Color, Dog2Mock.Gender, Dog2Mock.Breed, Dog2Mock.Weight)
	mock.ExpectQuery("^SELECT (.+) FROM dogs WHERE id=?").
		WithArgs("2").
		WillReturnRows(rows)

	handler := http.HandlerFunc(app.GetDogById)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("getDogs return wrong status code: got %v but want %v", recorder.Code, http.StatusOK)
	}
	expected := `{"id":2,"name":"Dog2","age":"3","color":"Grey","gender":"Male","breed":"Boxer","weight":"4"}
`
	if recorder.Body.String() != expected {
		t.Errorf("Wrong return body: got %s but want %s", recorder.Body.String(), expected)
	}
}

func TestPostDog(t *testing.T) {
	var jsonString = []byte(`{"name":"Dog1","age":"5","color":"Brown","gender":"Female","breed":"Husky","weight":"34"}`)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Fatalf("an error %s was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	app := &handlers.App{
		DB: db,
	}
	req, err := http.NewRequest("POST", "/dogs", bytes.NewBuffer(jsonString))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	query := "INSERT INTO dogs(name, age, color, gender, breed, weight) VALUES(?,?,?,?,?,?)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(Dog1Mock.Name, Dog1Mock.Age, Dog1Mock.Color, Dog1Mock.Gender, Dog1Mock.Breed, Dog1Mock.Weight).WillReturnResult(sqlmock.NewResult(0, 1))

	handler := http.HandlerFunc(app.PostDog)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("getDogs return wrong status code: got %v but want %v", recorder.Code, http.StatusCreated)
	}
	expected := `{"age":"5","breed":"Husky","color":"Brown","gender":"Female","name":"Dog1","weight":"34"}
`
	if recorder.Body.String() != expected {
		t.Errorf("Wrong return body: got %s but want %s", recorder.Body.String(), expected)
	}
}

func TestUpdateDog(t *testing.T) {
	var jsonString = []byte(`{"color":"Dark Brown"}`)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error %s was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	app := &handlers.App{
		DB: db,
	}

	req, err := http.NewRequest("PATCH", "/dogs", bytes.NewBuffer(jsonString))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	recorder := httptest.NewRecorder()

	query := "UPDATE dogs SET color = 'Dark Brown' WHERE id = 1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	handler := http.HandlerFunc(app.UpdateDog)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("getDogs return wrong status code: got %v but want %v", recorder.Code, http.StatusOK)
	}
	expected := `Dog with ID 1 was updated`
	if recorder.Body.String() != expected {
		t.Errorf("Wrong return body: got %s but want %s", recorder.Body.String(), expected)
	}
}

func TestDeleteDog(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an err %s was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	app := &handlers.App{
		DB: db,
	}

	req, err := http.NewRequest("DELETE", "/dogs", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{"id": "2"})
	recorder := httptest.NewRecorder()

	query := "DELETE FROM dogs WHERE id=?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WillReturnResult(sqlmock.NewResult(0, 1))

	handler := http.HandlerFunc(app.DeleteDog)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("getDogs return wrong status code: got %v but want %v", recorder.Code, http.StatusOK)
	}
	expected := `Dog with ID 2 was deleted`
	if recorder.Body.String() != expected {
		t.Errorf("Wrong return body: got %s but want %s", recorder.Body.String(), expected)
	}
}
