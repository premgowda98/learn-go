package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"project/restapi-2/types"
	"testing"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if user payload is invalid", func(t *testing.T){
		payload := types.RegisterUserPayload{
			FirstName: "",
			LastName: "",
			Email: "",
			Password: "",
		}

		marshelled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshelled))

		if err !=nil {
			t.Fatal(err)
		}
	})
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(u types.User) (error) {
	return nil
}
