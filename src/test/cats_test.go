package test

import (
	handlers "go_mysql/src/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetCats(t *testing.T) {
	// create new mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when openning a stub database connection", err)
	}
	defer db.Close()

	// use mock database
	app := &handlers.App{
		DB: db,
	}

	// create new request to /cats endpoint, nil => body io.Reader
	req, err := http.NewRequest("GET", "/cats", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Will store response received from /cats endpoint => pointer to ResonseRecorder struct
	recorder := httptest.NewRecorder()

	// insert mock data
	rows := sqlmock.NewRows([]string{"id", "name", "age", "color", "gender", "breed", "weight"}).
		AddRow(Cat1Mock.ID, Cat1Mock.Name, Cat1Mock.Age, Cat1Mock.Color, Cat1Mock.Gender, Cat1Mock.Breed, Cat1Mock.Weight).
		AddRow(Cat2Mock.ID, Cat2Mock.Name, Cat2Mock.Age, Cat2Mock.Color, Cat2Mock.Gender, Cat2Mock.Breed, Cat2Mock.Weight)
	mock.ExpectQuery("^SELECT (.+) FROM cats$").WillReturnRows(rows)

	// execute request with mock database filled with....mock data
	handler := http.HandlerFunc(app.GetCats)
	// hit endpoint with recorder and request
	handler.ServeHTTP(recorder, req)

	// test recorder status code
	if recorder.Code != http.StatusOK {
		t.Errorf("getCats return wrong status code: got %v but want %v", recorder.Code, http.StatusOK)
	}

	// test return body. for some
	//"Encode writes the JSON encoding of v to the stream, followed by a newline character." => add new line to expected
	expected := `[{"id":1,"name":"Cat1","age":"1","color":"White","gender":"female","breed":"Munchkin","weight":"8"},{"id":2,"name":"Cat2","age":"3","color":"Orange","gender":"Male","breed":"Bengal","weight":"4"}]
`
	if recorder.Body.String() != expected {
		t.Errorf("Wrong return body: got %s but want %s", recorder.Body.String(), expected)
	}
}
func TestGetCatById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when openning a stub database connection", err)
	}
	defer db.Close()

	app := &handlers.App{
		DB: db,
	}

	req, err := http.NewRequest("GET", "/cats/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	rows := sqlmock.NewRows([]string{"id", "name", "age", "color", "gender", "breed", "weight"}).
		AddRow(Cat1Mock.ID, Cat1Mock.Name, Cat1Mock.Age, Cat1Mock.Color, Cat1Mock.Gender, Cat1Mock.Breed, Cat1Mock.Weight).
		AddRow(Cat2Mock.ID, Cat2Mock.Name, Cat2Mock.Age, Cat2Mock.Color, Cat2Mock.Gender, Cat2Mock.Breed, Cat2Mock.Weight)
	mock.ExpectQuery("^SELECT (.+) FROM cats WHERE id=?").WithArgs("").WillReturnRows(rows)

	handler := http.HandlerFunc(app.GetCatById)
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("getCats return wrong status code: got %v but want %v", recorder.Code, http.StatusOK)
	}

	expected := `{"id":2,"name":"Cat2","age":"3","color":"Orange","gender":"Male","breed":"Bengal","weight":"4"}
`
	if recorder.Body.String() != expected {
		t.Errorf("Wrong return body: got %s but want %s", recorder.Body.String(), expected)
	}
}
