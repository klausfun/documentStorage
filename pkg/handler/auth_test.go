package handler

import (
	"bytes"
	"documentStorage/models"
	"documentStorage/pkg/service"
	mock_service "documentStorage/pkg/service/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           models.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"login":"newperson1","password":"new&Person112", "token":"rfjn4289jnd892vjdsi92uvhnjd8f"}`,
			inputUser: models.User{
				Login:    "newperson1",
				Password: "new&Person112",
				Token:    "rfjn4289jnd892vjdsi92uvhnjd8f",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return("newperson1", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"response":{"login":"newperson1"}}`,
		},
		{
			name:                "Invalid JSON",
			inputBody:           `{"login":"newperson1", "password":}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, user models.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":{"code":400,"text":"invalid input body"}}`,
		},
		{
			name:                "Invalid Login or Password",
			inputBody:           `{"login":"", "password":"new&Person112", "token":"rfjn4289jnd892vjdsi92uvhnjd8f"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, user models.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":{"code":400,"text":"invalid input body"}}`,
		},
		{
			name:                "No Access Rights",
			inputBody:           `{"login":"newperson1","password":"new&Person112", "token":"wvwrvwvwvrva"}`,
			inputUser:           models.User{Login: "newperson1", Password: "new&Person112", Token: "invalid_token"},
			mockBehavior:        func(s *mock_service.MockAuthorization, user models.User) {},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":{"code":403,"text":"no access rights"}}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"login":"newperson1","password":"new&Person112", "token":"rfjn4289jnd892vjdsi92uvhnjd8f"}`,
			inputUser: models.User{
				Login:    "newperson1",
				Password: "new&Person112",
				Token:    "rfjn4289jnd892vjdsi92uvhnjd8f",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return("", errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":{"code":500,"text":"service failure"}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/register", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(testCase.inputBody))

			if testCase.inputUser.Token != "" {
				os.Setenv("REGISTRATION_TOKEN", testCase.inputUser.Token)
			} else {
				os.Setenv("REGISTRATION_TOKEN", "rfjn4289jnd892vjdsi92uvhnjd8f")
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
