package service

import (
	"errors"
	"test/restapi/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := NewMockUserStore(ctrl)

	mockData := &models.User{ID: 1, Name: "example"}

	tcs := []struct {
		desc      string
		input     *models.UserRequest
		mockCalls []*gomock.Call
		expErr    error
	}{
		{
			desc:  "success case: create",
			input: &models.UserRequest{ID: 1, Name: "example"},
			mockCalls: []*gomock.Call{
				mockStore.EXPECT().Create(mockData).Return(nil),
			},
			expErr: nil,
		},
		{
			desc:  "success case: fail",
			input: &models.UserRequest{ID: 1, Name: "fail"},
			mockCalls: []*gomock.Call{
				mockStore.EXPECT().Create(mockData).Return(nil),
			},
			expErr: errors.New("not allowed"),
		},
	}

	for i, tc := range tcs {
		svc := New(mockStore)

		err := svc.Create(tc.input)

		assert.Equalf(t, tc.expErr, err, "Test[%v] failed", i)
	}
}
