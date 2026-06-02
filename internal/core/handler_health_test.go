package core

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	"google.golang.org/genai"
)

func TestHealth(t *testing.T) {
	var tests = []struct {
		name     string
		expected healthResponse
		initFn   func() *App
	}{
		{
			name: "Healthy response",
			expected: healthResponse{
				Status: "OK",
				Uptime: "0s",
				Services: services{
					Database: "UP",
					Model:    "UP",
				},
			},
			initFn: func() *App {
				mock := New(newTestDependencies(&mockAI{
					ListFn: func(
						ctx context.Context,
						config *genai.ListModelsConfig,
					) (genai.Page[genai.Model], error) {
						return genai.Page[genai.Model]{}, nil
					},
				}))
				return mock
			},
		},
		{
			name: "Unhealthy response",
			expected: healthResponse{
				Status: "UNHEALTHY",
				Uptime: "0s",
				Services: services{
					Database: "DOWN",
					Model:    "DOWN",
				},
			},
			initFn: func() *App {
				mock := New(newTestDependencies(&mockAI{
					ListFn: func(
						ctx context.Context,
						config *genai.ListModelsConfig,
					) (genai.Page[genai.Model], error) {
						return genai.Page[genai.Model]{},
							errors.New("could not list models")
					},
				}))
				mock.DB.Close()
				return mock
			},
		},
		{
			name: "Database up, Model down",
			expected: healthResponse{
				Status: "UNHEALTHY",
				Uptime: "0s",
				Services: services{
					Database: "UP",
					Model:    "DOWN",
				},
			},
			initFn: func() *App {
				mock := New(newTestDependencies(&mockAI{
					ListFn: func(
						ctx context.Context,
						config *genai.ListModelsConfig,
					) (genai.Page[genai.Model], error) {
						return genai.Page[genai.Model]{},
							errors.New("could not list models")
					},
				}))
				return mock
			},
		},
		{
			name: "Model up, Database down",
			expected: healthResponse{
				Status: "UNHEALTHY",
				Uptime: "0s",
				Services: services{
					Database: "DOWN",
					Model:    "UP",
				},
			},
			initFn: func() *App {
				mock := New(newTestDependencies(&mockAI{}))
				mock.DB.Close()
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.initFn()
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", baseURL+"/health", nil)

			mock.Health(recorder, request)

			var actual healthResponse
			if err := json.NewDecoder(recorder.Body).Decode(&actual); err != nil {
				t.Fatal("response: failed to decode response body")
			}

			if !reflect.DeepEqual(tt.expected, actual) {
				t.Fatalf("response: want %v, got %v", tt.expected, actual)
			}
		})
	}
}
