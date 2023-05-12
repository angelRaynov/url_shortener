package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorReader struct{}

func (r *errorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("error reading")
}

// One million iterations
// go test -bench=BenchmarkShortenURLHandler -benchmem -cpu=8 -benchtime=1000000x
// 1000000	     22225 ns/op	    2510 B/op	      22 allocs/op
func BenchmarkShortenURLHandler(b *testing.B) {

	uniqueURLs := make(chan string)

	go generateUniqueURLs(uniqueURLs)

	rr := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request) {
		shortenURLHandler(w, r, uniqueURLs)

	}

	for n := 0; n < b.N; n++ {
		url := fmt.Sprintf("https://www.someverylogngurl%d.com", n)
		req, err := http.NewRequest("GET", "/", strings.NewReader(url))
		if err != nil {
			b.Fatal(err)
		}
		handler(rr, req)

		body, err := io.ReadAll(rr.Body)
		if err != nil {
			b.Fatal(err)
		}

		//request the long url
		resolveURL := string(body)
		req2, err := http.NewRequest("GET", resolveURL, nil)
		if err != nil {
			b.Fatal(err)
		}
		handler(rr, req2)
	}
}

// 100k iterations
// go test -bench=BenchmarkShortenURLHandlerOneMap -benchmem -cpu=8 -benchtime=100000x
// 100000	    417372 ns/op	    2596 B/op	      29 allocs/op
func BenchmarkSlow(b *testing.B) {

	uniqueURLs := make(chan string)

	go generateUniqueURLs(uniqueURLs)

	rr := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request) {
		shortenURLHandlerOneMap(w, r, uniqueURLs)

	}

	for n := 0; n < b.N; n++ {
		url := fmt.Sprintf("https://www.someverylogngurl%d.com", n)
		req, err := http.NewRequest("GET", "/", strings.NewReader(url))
		if err != nil {
			b.Fatal(err)
		}
		handler(rr, req)

		body, err := io.ReadAll(rr.Body)
		if err != nil {
			b.Fatal(err)
		}

		//request the long url
		resolveURL := string(body)
		req2, err := http.NewRequest("GET", resolveURL, nil)
		if err != nil {
			b.Fatal(err)
		}
		handler(rr, req2)
	}
}

func TestShortenURLHandler(t *testing.T) {
	uniqueURLs, expectUniqueURLs := startTestChannels()

	expectedShortURL := "http://localhost:1234/" + <-expectUniqueURLs

	// Create a test request to shorten a URL
	shortenBody := []byte("https://www.example.com/test1")
	shortenRequest, err := http.NewRequest("POST", "/", bytes.NewReader(shortenBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function
	shortenURLHandler(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, shortenRecorder.Code)
	}

	// Check the response body
	if shortenRecorder.Body.String() != expectedShortURL {
		t.Errorf("expected body %q; got %q", expectedShortURL, shortenRecorder.Body.String())
	}

}

func TestShortenURLHandlerNilBody(t *testing.T) {
	uniqueURLs, _ := startTestChannels()

	// Create a test request to shorten a URL
	shortenRequest, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function
	shortenURLHandler(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, shortenRecorder.Code)
	}

}

func TestShortenURLHandlerInvalidBody(t *testing.T) {
	uniqueURLs, _ := startTestChannels()

	// Create a test request to shorten a URL
	shortenRequest, err := http.NewRequest("POST", "/", &errorReader{})
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function
	shortenURLHandler(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, shortenRecorder.Code)
	}

}

func TestShortenURLHandlerSameURL(t *testing.T) {
	uniqueURLs, expectUniqueURLs := startTestChannels()

	expectedShortURL := "http://localhost:1234/" + <-expectUniqueURLs

	shortenBody := []byte("https://www.example.com/test1")

	// Create a test request to shorten the same URL again
	shortenRequest, err := http.NewRequest("POST", "/", bytes.NewReader(shortenBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function again
	shortenURLHandler(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, shortenRecorder.Code)
	}

	// Check the response body
	if shortenRecorder.Body.String() != expectedShortURL {
		t.Errorf("expected body %q; got %q", expectedShortURL, shortenRecorder.Body.String())
	}

}

func TestShortenURLHandlerExpandURL(t *testing.T) {
	uniqueURLs, expectUniqueURLs := startTestChannels()

	expectedShortURL := "http://localhost:1234/" + <-expectUniqueURLs

	// Create a test request to shorten a URL
	shortenBody := []byte("https://www.example.com/test1")
	shortenRequest, err := http.NewRequest("POST", "/", bytes.NewReader(shortenBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function
	shortenURLHandler(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, shortenRecorder.Code)
	}

	// Check the response body
	if shortenRecorder.Body.String() != expectedShortURL {
		t.Errorf("expected body %q; got %q", expectedShortURL, shortenRecorder.Body.String())
	}

	// Create a test request to expand the short URL
	expandRequest, err := http.NewRequest("GET", expectedShortURL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	expandRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function to expand the URL
	shortenURLHandler(expandRecorder, expandRequest, uniqueURLs)

	// Check the response status code
	if expandRecorder.Code != http.StatusMovedPermanently {
		t.Errorf("expected status %d; got %d", http.StatusMovedPermanently, expandRecorder.Code)
	}

	// Check the response header Location
	if expandRecorder.Header().Get("Location") != string(shortenBody) {
		t.Errorf("expected header Location %q; got %q", string(shortenBody), expandRecorder.Header().Get("Location"))
	}

}

func TestFindLongUrl(t *testing.T) {
	haystack := map[string]string{
		"http://example.com": "abc123",
		"http://google.com":  "def456",
		"http://amazon.com":  "ghi789",
	}

	testCases := []struct {
		name            string
		needle          string
		expectedLongUrl string
		expectedFound   bool
	}{
		{"value exist", "abc123", "http://example.com", true},
		{"value does not exist", "xyz789", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			longUrl, found := findLongUrl(haystack, tc.needle)
			if longUrl != tc.expectedLongUrl || found != tc.expectedFound {
				t.Errorf("findLongUrl(%v, %v) = (%v, %v), expected (%v, %v)", haystack, tc.needle, longUrl, found, tc.expectedLongUrl, tc.expectedFound)
			}
		})

	}
}

func TestShortenURLOneMapHandler(t *testing.T) {
	uniqueURLs, expectUniqueURLs := startTestChannels()

	expectedShortURL := "http://localhost:1234/" + <-expectUniqueURLs

	// Create a test request to shorten a URL
	shortenBody := []byte("https://www.example.com/test1")
	shortenRequest, err := http.NewRequest("POST", "/", bytes.NewReader(shortenBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function
	shortenURLHandlerOneMap(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, shortenRecorder.Code)
	}

	// Check the response body
	if shortenRecorder.Body.String() != expectedShortURL {
		t.Errorf("expected body %q; got %q", expectedShortURL, shortenRecorder.Body.String())
	}

}

func TestShortenURLOneMapHandlerNilBody(t *testing.T) {
	uniqueURLs, _ := startTestChannels()

	// Create a test request to shorten a URL
	shortenRequest, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function
	shortenURLHandlerOneMap(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, shortenRecorder.Code)
	}

}

func TestShortenURLOneMapHandlerInvalidBody(t *testing.T) {
	uniqueURLs, _ := startTestChannels()

	// Create a test request to shorten a URL
	shortenRequest, err := http.NewRequest("POST", "/", &errorReader{})
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function
	shortenURLHandlerOneMap(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, shortenRecorder.Code)
	}

}

func TestShortenURLOneMapHandlerSameURL(t *testing.T) {
	uniqueURLs, expectUniqueURLs := startTestChannels()

	expectedShortURL := "http://localhost:1234/" + <-expectUniqueURLs

	shortenBody := []byte("https://www.example.com/test1")

	// Create a test request to shorten the same URL again
	shortenRequest, err := http.NewRequest("POST", "/", bytes.NewReader(shortenBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function again
	shortenURLHandlerOneMap(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, shortenRecorder.Code)
	}

	// Check the response body
	if shortenRecorder.Body.String() != expectedShortURL {
		t.Errorf("expected body %q; got %q", expectedShortURL, shortenRecorder.Body.String())
	}

}

func TestShortenURLHandlerOneMapExpandURL(t *testing.T) {
	uniqueURLs, expectUniqueURLs := startTestChannels()

	expectedShortURL := "http://localhost:1234/" + <-expectUniqueURLs

	// Create a test request to shorten a URL
	shortenBody := []byte("https://www.example.com/test1")
	shortenRequest, err := http.NewRequest("POST", "/", bytes.NewReader(shortenBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	shortenRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function
	shortenURLHandlerOneMap(shortenRecorder, shortenRequest, uniqueURLs)

	// Check the response status code
	if shortenRecorder.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, shortenRecorder.Code)
	}

	// Check the response body
	if shortenRecorder.Body.String() != expectedShortURL {
		t.Errorf("expected body %q; got %q", expectedShortURL, shortenRecorder.Body.String())
	}

	// Create a test request to expand the short URL
	expandRequest, err := http.NewRequest("GET", expectedShortURL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a test response recorder
	expandRecorder := httptest.NewRecorder()

	// Call the shortenURLHandler function to expand the URL
	shortenURLHandlerOneMap(expandRecorder, expandRequest, uniqueURLs)

	// Check the response status code
	if expandRecorder.Code != http.StatusMovedPermanently {
		t.Errorf("expected status %d; got %d", http.StatusMovedPermanently, expandRecorder.Code)
	}

	// Check the response header Location
	if expandRecorder.Header().Get("Location") != string(shortenBody) {
		t.Errorf("expected header Location %q; got %q", string(shortenBody), expandRecorder.Header().Get("Location"))
	}

}

func startTestChannels() (chan string, chan string) {
	uniqueURLs := make(chan string)
	go generateUniqueURLs(uniqueURLs)

	//we will consume this channel for the expected test cases
	expectUniqueURLs := make(chan string)
	go generateUniqueURLs(expectUniqueURLs)
	return uniqueURLs, expectUniqueURLs
}
