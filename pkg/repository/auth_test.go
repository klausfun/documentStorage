package repository

import (
	"documentStorage/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewAuthPostgres(db, nil)

	type mockBehavior func(args models.User)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		user         models.User
		wantErr      bool
		wantLogin    string
	}{
		{
			name: "OK",
			user: models.User{
				Login:    "newperson1",
				Password: "new&Person112",
			},
			mockBehavior: func(args models.User) {
				rows := sqlmock.NewRows([]string{"login"}).AddRow(args.Login)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(args.Login, args.Password).
					WillReturnRows(rows)
			},
			wantLogin: "newperson1",
		},
		{
			name: "Error on Insert",
			user: models.User{
				Login:    "newperson1",
				Password: "new&Person112",
			},
			mockBehavior: func(args models.User) {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(args.Login, args.Password).
					WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.user)

			gotLogin, err := r.CreateUser(testCase.user)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.wantLogin, gotLogin)
			}
		})
	}
}
