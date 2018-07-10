package main

import(
    "net/http"
    "net/http/httptest"
    "io/ioutil"
    "testing"
)

func TestHandler(t *testing.T){
    req, err := http.NewRequest("GET", "", nil)
    
    if err != nil{
        t.Fatal(err)
    }


    recorder := httptest.NewRecorder()

    hf := http.HandlerFunc(handler)
    hf.ServeHTTP(recorder, req)

    if status := recorder.Code; status != http.StatusOK{
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := `Hello World!`
    actual := recorder.Body.String()
    if actual != expected{
        t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
    }
}

func TestRouter(t *testing.T){
    r := newRouter()
    mockServer := httptest.NewServer(r)
    
    resp, err := http.Get(mockServer.URL + "/hello")
    
    if err != nil{
        t.Fatal(err)
    }
    
    if resp.StatusCode != http.StatusOK{
        t.Errorf("Status should be ok, got %d", resp.StatusCode)
    }
    
    defer resp.Body.Close()
    
    b, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
        t.Fatal(err)
    }
    
    respString := string(b)
    expected := "Hello World!"
    
    if respString != expected{
        t.Errorf("Response should be %s, got %s", expected, respString)
    }
}

func TestRouterForNonExistentRoute(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	// Most of the code is similar. The only difference is that now we make a 
	//request to a route we know we didn't define, like the `POST /hello` route.
	resp, err := http.Post(mockServer.URL+"/hello", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 405 (method not allowed)
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status should be 405, got %d", resp.StatusCode)
	}

	// The code to test the body is also mostly the same, except this time, we 
	// expect an empty body
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	respString := string(b)
	expected := ""

	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}
}


func TestStaticFileServer(t *testing.T){
    r := newRouter()
    mockServer := httptest.NewServer(r)
    
    resp, err := http.Get(mockServer.URL + "/assets/")
    if err != nil{
        t.Fatal(err)
    }
    
    contentType := resp.Header.Get("Content-Type")
    expectedContentType := "text/html; charset=utf-8"
    
    if expectedContentType != contentType{
        t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
    }
}
    
    
    
    
    
    
    
    
    
    
    
    