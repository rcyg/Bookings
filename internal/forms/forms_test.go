package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	//valid function only checks whether the form received has errors
	//we can use it for further test function
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm) //got a empty form which has no error

	isValid := form.Valid() //if empty should be no errors
	if !isValid {           //if not empty
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	//required function adding errors if required field is empty
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c") //check fields which is empty
	if form.Valid() {            //it should got errors and return false
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c") //it has these fields thus have no errors actually
	if !form.Valid() {           //it should return true
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(url.Values{})

	has := form.Has("whatever")
	if has { // it don't has so it should return false
		t.Error("form shows has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	has = form.Has("a")
	if !has { //it actually has so it should return true
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("x", 10)
	if form.Valid() { //it has errors so it should be false
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" { //it actually got errors so the isError should not be blank
		t.Error("should have error but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("some_field", "some_value")

	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() { //not satify the length requirement it got errors
		t.Error("shows min length of 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("another_field", "abc123")
	form = New(postedData)

	form.MinLength("another_field", 1)
	if !form.Valid() { //satisfy the length requirement so it will pass
		t.Error("shows min length if 1 is not met when it is")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" { //it pass the length requirement so it will pass
		t.Error("should not have error but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() { //field not exist, shall not pass
		t.Error("form shows valid email for non-existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "me@here.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() { //actually a email, should pass
		t.Error("got an invalid email when we should not have")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "x")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() { //invalid email, shall not pass
		t.Error("got valid for invalid email address")
	}
}
