package auth

import (
	"reflect"
	"testing"

	"github.com/sschwartz96/minimongo/db"
	"github.com/sschwartz96/minimongo/mock"
	"github.com/sschwartz96/syncapod/internal/database"
	"github.com/sschwartz96/syncapod/internal/models"
	"github.com/sschwartz96/syncapod/internal/protos"
)

func TestCreateAuthorizationCode(t *testing.T) {
	mockDB := mock.CreateDB()

	type args struct {
		dbClient db.Database
		userID   *protos.ObjectID
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				dbClient: mockDB,
				clientID: "testClient",
				userID:   protos.ObjectIDFromHex("testUserID"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateAuthorizationCode(tt.args.dbClient, tt.args.userID, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAuthorizationCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var found models.AuthCode
				err = tt.args.dbClient.FindOne(database.ColAuthCode, &found, &db.Filter{"code": got}, db.CreateOptions())
				if err != nil {
					t.Errorf("CreateAuthorizationCode() error looking for auth code = %v", err)
				}
			}
		})
	}
}

func TestCreateAccessToken(t *testing.T) {
	type args struct {
		dbClient db.Database
		authCode *models.AuthCode
	}
	tests := []struct {
		name    string
		args    args
		want    *models.AccessToken
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateAccessToken(tt.args.dbClient, tt.args.authCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestValidateAuthCode(t *testing.T) {
	type args struct {
		dbClient db.Database
		code     string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.AuthCode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAuthCode(tt.args.dbClient, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAuthCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateAuthCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	type args struct {
		dbClient db.Database
		token    string
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAccessToken(tt.args.dbClient, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insertAuthCode(t *testing.T) {
	type args struct {
		dbClient db.Database
		code     *models.AuthCode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := insertAuthCode(tt.args.dbClient, tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("insertAuthCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findAuthCode(t *testing.T) {
	type args struct {
		dbClient db.Database
		code     string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.AuthCode
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findAuthCode(tt.args.dbClient, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("findAuthCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAuthCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insertAccessToken(t *testing.T) {
	type args struct {
		dbClient db.Database
		token    *models.AccessToken
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := insertAccessToken(tt.args.dbClient, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("insertAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindOauthAccessToken(t *testing.T) {
	type args struct {
		dbClient db.Database
		token    string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.AccessToken
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindOauthAccessToken(tt.args.dbClient, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOauthAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOauthAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteOauthAccessToken(t *testing.T) {
	type args struct {
		dbClient db.Database
		token    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteOauthAccessToken(tt.args.dbClient, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("DeleteOauthAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
