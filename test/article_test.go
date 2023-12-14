package test

import (
	"net/http"
	"os"
	"testing"
)

// USER
func TestUserGetAllArticles(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		expected    int
	}{
		{"ValidParams", "?limit=10&offset=0", http.StatusOK},
		{"MissingLimit", "?offset=0", http.StatusBadRequest},
		{"MissingOffset", "?limit=10", http.StatusBadRequest},
		{"MissingBoth", "", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := BaseURL + "/users/articles" + tt.queryParams

			response, err := http.Get(url)
			if err != nil {
				t.Fatalf("Error making GET request: %s", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tt.expected {
				t.Errorf("Expected status code %d, got %d", tt.expected, response.StatusCode)
			}
		})
	}
}

func TestUserGetArticleByID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected int
	}{
		{"ValidID", "33", http.StatusOK},
		{"NotFoundID", "10", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := BaseURL + "/users/articles/" + tt.id

			response, err := http.Get(url)
			if err != nil {
				t.Fatalf("Error making GET request: %s", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tt.expected {
				t.Errorf("Expected status code %d, got %d", tt.expected, response.StatusCode)
			}
		})
	}
}

func TestUserGetAllArticlesByTitle(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		expected    int
	}{
		{"ValidParams", "?title=coba&limit=5&offset=0", http.StatusOK},
		{"InvalidValueTitleParams", "?title=xyz&limit=5&offset=0", http.StatusNotFound},
		{"MissingTitle", "?limit=5&offset=0", http.StatusBadRequest},
		{"MissingLimit", "?title=coba&offset=0", http.StatusBadRequest},
		{"MissingOffset", "?title=coba&limit=5", http.StatusBadRequest},
		{"MissingAll", "", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := BaseURL + "/users/article" + tt.queryParams

			response, err := http.Get(url)
			if err != nil {
				t.Fatalf("Error making GET request: %s", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tt.expected {
				t.Errorf("Expected status code %d, got %d", tt.expected, response.StatusCode)
			}
		})
	}
}

// DOCTOR
var DoctorValidBearerToken = os.Getenv("DOCTOR_TOKEN")

func TestDoctorDeleteArticle(t *testing.T) {}
