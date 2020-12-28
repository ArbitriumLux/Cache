package server

import (
	"D/Avito/cache"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer_Delete(t *testing.T) {
	testCase := []struct {
		name string
		key  string
	}{
		{
			name: "Actualkey",
			key:  "45",
		},
		{
			name: "Notakey",
			key:  "fefefefefef",
		},
	}
	s := New(NewConfig())
	handler := http.HandlerFunc(s.Delete)
	s.cache.Items["45"] = cache.Item{
		Value:      "pppp",
		Created:    time.Now(),
		Expiration: 1,
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{
				"Key1": []string{tc.key},
			}
			req, err := http.NewRequest("POST", "/delete", strings.NewReader(form.Encode()))
			assert.Nil(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			_, ok := s.cache.Items[tc.key]
			assert.Equal(t, false, ok)
		})
	}
}
func TestServer_Set(t *testing.T) {
	testCase := []struct {
		name     string
		key      string
		value    string
		duration string
	}{
		{
			name:     "Correct",
			key:      "Len",
			value:    "Some words",
			duration: "4m",
		},
		{
			name:     "Notcorrect",
			key:      "23",
			value:    "",
			duration: "0",
		},
	}
	s := New(NewConfig())
	handler := http.HandlerFunc(s.Set)
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{
				"Key":      []string{tc.key},
				"Value":    []string{tc.value},
				"Duration": []string{tc.duration},
			}
			req, err := http.NewRequest("POST", "/set", strings.NewReader(form.Encode()))
			assert.Nil(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			_, ok := s.cache.Items[tc.key]
			assert.Equal(t, true, ok)
		})
	}
}
