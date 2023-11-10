package api

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/nghianm93/romo/types"
	"net/http/httptest"
	"testing"
)

const (
	TestFirstName = "Tester"
	TestLastName  = "Nguyen Van"
	TestEmail     = "nvtester@gmail.com"
	TestPassWord  = "abc123xyz"
)

func TestUserHandler_HandlePostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.dropDown(t)
	apptest := fiber.New()
	UserHandler := NewUserHandler(tdb.User)
	apptest.Post("/", UserHandler.HandlePostUser)
	params := types.CreateUserParams{
		FirstName: TestFirstName,
		LastName:  TestLastName,
		Email:     TestEmail,
		Password:  TestPassWord,
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := apptest.Test(req)
	if err != nil {
		t.Error(err)
	}
	var User types.User
	json.NewDecoder(res.Body).Decode(&User)
	if len(User.ID) == 0 {
		t.Errorf("User should have Id")
	}
	if len(User.EncryptedPassword) > 0 {
		t.Errorf("Don't return pwd on json")
	}
	if User.FirstName != TestFirstName {
		t.Errorf("Expexted %s but got %s", TestFirstName, User.FirstName)
	}
	if User.LastName != TestLastName {
		t.Errorf("Expected %s but got %s", TestLastName, User.LastName)
	}
	if User.Email != TestEmail {
		t.Errorf("Expected %s but got %s", TestEmail, User.Email)
	}

}
