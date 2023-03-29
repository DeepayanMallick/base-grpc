package postgres

import (
	"context"
	"testing"

	"github.com/DeepayanMallick/base-grpc/usermgm/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func insertTestUser(t *testing.T, s *Storage) {
	// insert test user into the "users" table
	got, err := s.CreateUser(context.TODO(), storage.User{FirstName: "testfirstname", LastName: "testlastname", Username: "testuser", Email: "testuser@example.com", Password: "password"})
	if err != nil {
		t.Fatalf("Unable to create user: %v", err)
	}

	if got == "" {
		t.Fatalf("Unable to create user: %v", err)
	}
}

func deleteTestUser(t *testing.T, s *Storage, id string) {
	// delete test user from the "users" table
	if err := s.deleteUserPermanently(context.TODO(), id); err != nil {
		t.Fatalf("Unable to delete user: %v", err)
	}
}

func deleteAllTestUser(t *testing.T, s *Storage) {
	// delete test user from the "users" table
	if err := s.deleteUsersPermanently(context.TODO()); err != nil {
		t.Fatalf("Unable to delete users: %v", err)
	}
}

var getUserByIDtestCases = []struct {
	name     string
	id       string
	want     *storage.User
	wantErr  bool
	teardown func(*testing.T, *Storage, string)
	setup    func(*testing.T, *Storage, string)
}{
	{
		name:     "validUserID",
		want:     &storage.User{FirstName: "testfirstname", LastName: "testlastname", Username: "testuser", Email: "testuser@example.com", Password: "password"},
		teardown: deleteTestUser,
	},
	{
		name:    "nonExistentUserID",
		id:      "nonexistentuserid",
		wantErr: true,
	},
}

func TestGetUserByID(t *testing.T) {
	s := newTestStorage(t)

	// create new user
	uid, err := s.CreateUser(context.TODO(), storage.User{FirstName: "testfirstname", LastName: "testlastname", Username: "testuser", Email: "testuser@example.com", Password: "password"})
	if err != nil {
		t.Fatalf("Unable to create user: %v", err)
	}

	opts := cmp.Options{
		cmpopts.IgnoreFields(
			storage.User{},
			"ID",
			"Password",
			"CreatedAt",
			"UpdatedAt",
			"DeletedAt",
			"DeletedBy",
		),
	}

	for _, tc := range getUserByIDtestCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.id != "" {
				uid = tc.id
			}

			got, err := s.GetUserByID(context.Background(), uid)
			if (err != nil) != tc.wantErr {
				t.Errorf("Storage.GetUserByID() gotErr = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !cmp.Equal(got, tc.want, opts...) {
				t.Errorf("Storage.GetUserByID() diff = %v", cmp.Diff(got, tc.want, opts...))
			}

			if tc.teardown != nil {
				tc.teardown(t, s, uid)
			}
		})
	}
}

var testCasesGetUserByEmail = []struct {
	name     string
	email    string
	want     *storage.User
	wantErr  bool
	setup    func(*testing.T, *Storage)
	teardown func(*testing.T, *Storage, string)
}{
	{
		name:     "valid username",
		email:    "testuser@example.com",
		want:     &storage.User{FirstName: "testfirstname", LastName: "testlastname", Username: "testuser", Email: "testuser@example.com", Password: "password"},
		setup:    insertTestUser,
		teardown: deleteTestUser,
	},
	{
		name:    "non existent user email",
		email:   "nonexistenttestuser@example.com",
		wantErr: true,
	},
}

func TestGetUserByEmail(t *testing.T) {
	s := newTestStorage(t)

	opts := cmp.Options{
		cmpopts.IgnoreFields(
			storage.User{},
			"ID",
			"Password",
			"CreatedAt",
			"UpdatedAt",
			"DeletedAt",
			"DeletedBy",
		),
	}

	for _, tc := range testCasesGetUserByEmail {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(t, s)
			}

			got, err := s.GetUserByEmail(context.Background(), tc.email)
			if (err != nil) != tc.wantErr {
				t.Errorf("Storage.GetUserByEmail() gotErr = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !cmp.Equal(got, tc.want, opts...) {
				t.Errorf("Storage.GetUserByEmail() diff = %v", cmp.Diff(got, tc.want, opts...))
			}

			if tc.teardown != nil {
				tc.teardown(t, s, got.ID)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	s := newTestStorage(t)
	testCases := []struct {
		name     string
		in       storage.User
		want     string
		wantErr  bool
		setup    func(*testing.T, *Storage)
		teardown func(*testing.T, *Storage)
	}{
		{
			name: "user creation success",
			in: storage.User{
				FirstName: "testfirstname",
				LastName:  "testlastname",
				Username:  "testuser",
				Email:     "testuser@example.com",
				Password:  "password",
				Status:    1,
			},
			wantErr:  false,
			teardown: deleteAllTestUser,
		}, {
			name: "user creation failed",
			in: storage.User{
				FirstName: "testfirstname",
				LastName:  "testlastname",
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "password",
				Status:   1,
			},
			wantErr:  true,
			setup:    insertTestUser,
			teardown: deleteAllTestUser,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(t, s)
			}

			got, err := s.CreateUser(context.TODO(), tc.in)
			if (err != nil) != tc.wantErr {
				t.Errorf("Storage.CreateUser() gotErr = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && got == "" {
				t.Errorf("Storage.CreateUser() want uuid, got %v", got)
				return
			}

			if tc.teardown != nil {
				tc.teardown(t, s)
			}
		})
	}
}
