package service

import (
	"test/restapi/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_ServiceCreate(t *testing.T) {
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
	}

	for i, tc := range tcs {
		svc := New(mockStore)

		err := svc.Create(tc.input)

		assert.Equalf(t, tc.expErr, err, "Test[%v] failed", i)
	}
}

func Test_ServiceGet(t *testing.T){
	ctrl := gomock.NewController(t)
	mockStore := NewMockUserStore(ctrl)

	tests := []struct{
		input int
		expectedOutput models.User
	}{
		{1, models.User{ID: 1, Name: "ko"}},
	}

	for _, test :=  range tests{
		mockStore.EXPECT().Get(test.input).Return(&models.User{ID: 1, Name: "ko"}, nil)

		svc := New(mockStore)
		_, err := svc.Get(test.input)

		if err!=nil{
			t.Fatalf("not expecting error")
		}
	}
}
