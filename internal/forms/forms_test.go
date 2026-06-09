package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

// NOTE: Has and MinLenght function doesn't require r *http.Request, we need to change that in forms, but for learning purposes I had kept it as it is

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	// r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	form := New(url.Values{})
	r := httptest.NewRequest("GET", "/whatever", nil)

	has := form.Has("whatever", r)
	if has {
		t.Error("Form shows has field when it doesn't")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	r.PostForm = postData
	_ = r.ParseForm()

	form = New(r.PostForm)
	has = form.Has("a", r)
	if !has {
		t.Error("The form should have the field a")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("GET", "/whatever", nil)

	form := New(r.PostForm)
	_ = r.ParseForm()

	form.MinLength("first_name", 3, r)
	if form.Valid() {
		t.Error("form shows min length for non existent field")
	}

	isError := form.Errors.Get("first_name")
	if isError == "" {
		t.Error("should have got an error but didn't find one")
	}

	r = httptest.NewRequest("GET", "/whatever", nil)
	postData := url.Values{}
	postData.Add("some_field", "some_value")
	r.PostForm = postData
	form = New(r.PostForm)
	form.Errors = errors{}

	_ = r.ParseForm()

	_ = form.MinLength("some_field", 3, r)

	if !form.Valid() {
		t.Error("shows min length if 3 is not met, what it is")
	}

	isError = form.Errors.Get("some_field")
	if isError != "" {
		t.Error("should have got no error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	// TODO: test IsEmail works as desired
	postData := url.Values{}

	form := New(postData)

	form.IsEmail("x")

	if form.Valid() {
		t.Error("There should have been an error but found none")
	}

	postData.Set("email", "a@amail.com")
	form = New(postData)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("There should not have been an error but found one")
	}

	postData.Set("email", "x")
	form = New(postData)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("There should have been an error but found none")
	}
}
