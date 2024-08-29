package register_test

import (
	"fmt"
	"providerHub/internal/handler/auth"
)

type mockUserRegister struct{}

func (m *mockUserRegister) Register(req auth.RegisterRequest) (string, error) {
	if req.Login == "existing_user" {
		return "", fmt.Errorf("user already exists")
	}
	return "new_user_id", nil
}

/*func TestRegister(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	handler := register.New(logger, &mockUserRegister{})

	tests := []struct {
		name           string
		requestBody    dto.RegisterRequest
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful registration",
			requestBody: dto.RegisterRequest{
				Login:    "new_user",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "new_user_id",
		},
		{
			name: "user already exists",
			requestBody: dto.RegisterRequest{
				Login:    "existing_user",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"user already exists"}`,
		},
		{
			name: "invalid request",
			requestBody: dto.RegisterRequest{
				Login:    "",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}*/
