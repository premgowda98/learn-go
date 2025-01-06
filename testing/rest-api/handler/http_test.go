package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"test/restapi/models"
	"testing"

	"go.uber.org/mock/gomock"
)

func Test_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := NewMockUserService(ctrl)

	validData := []byte(`{"name": "ko"}`)
	r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(validData))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	mockService.EXPECT().Create(&models.UserRequest{ID: 0, Name: "ko"}).Return(errors.New("something went wrong"))
	
	handler := New(mockService)

	handler.CreateUser(w, r)

	if w.Code != http.StatusOK{
		t.Fatalf("failed to send http request")
	}

}
