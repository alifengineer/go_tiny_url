package postgres

import (
	"context"
	pb "go_auth_api_gateway/genproto/auth_service"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func createUser(t *testing.T) (resp *pb.UserPrimaryKey) {

	repo := NewUserRepo(db)

	resp, err := repo.Create(context.Background(), &pb.CreateUserRequest{
		FirstName: fakeData.FirstName(),
		LastName:  fakeData.LastName(),
		Phone:     fakeData.PhoneNumber(),
		Username:  fakeData.UserName(),
		Password:  "123456",
	})
	assert.NoError(t, err)
	return resp
}

func TestCreateUser(t *testing.T) {

	tests := []struct {
		name    string
		give    *pb.CreateUserRequest
		wantErr error
	}{
		{
			name: "SUCCESS: Create user",
			give: &pb.CreateUserRequest{
				FirstName: "Said",
				LastName:  "Amir",
				Phone:     fakeData.PhoneNumber(),
				Username:  fakeData.UserName(),
			},
			wantErr: nil,
		},
		{
			name: "ERROR: Create user, Dublicate phone",
			give: &pb.CreateUserRequest{
				FirstName: "Said",
				LastName:  "Amir",
				Phone:     fakeData.PhoneNumber(),
				Username:  fakeData.UserName(),
			},
			wantErr: nil,
		},
	}

	userRepo := NewUserRepo(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := userRepo.Create(context.Background(), tt.give)
			assert.NoError(t, err)
		})
	}
}

func TestGetUserByPK(t *testing.T) {

	tests := []struct {
		name string
		give *pb.UserPrimaryKey
		want error
	}{
		{
			name: "SUCCESS: Get user",
			give: &pb.UserPrimaryKey{
				Id: createUser(t).GetId(),
			},
			want: nil,
		},
		{
			name: "ERROR: Get user, User not found",
			give: &pb.UserPrimaryKey{
				Id: createRandomId(t),
			},
			want: pgx.ErrNoRows,
		},
	}

	userRepo := NewUserRepo(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := userRepo.GetByPK(context.Background(), tt.give)
			assert.Equal(t, tt.want, err)
		})
	}
}
