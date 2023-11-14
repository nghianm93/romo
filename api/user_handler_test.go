package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nghianm93/romo/types"
	"net/http/httptest"
	"testing"
)

const (
	TestFirstName   = "Tester"
	TestLastName    = "Nguyen Van"
	TestEmail       = "nvtester@gmail.com"
	TestPassWord    = "abc123xyz"
	UpdateFirstName = "Update FirstName"
	UpdateLastName  = "Update LasttName"
)

func createUserAndReturn(apptest *fiber.App, t *testing.T, firstName, lastName, email, password string) types.User {
	createUserParams := types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}

	createUserRequest, _ := json.Marshal(createUserParams)
	createReq := httptest.NewRequest("POST", "/", bytes.NewReader(createUserRequest))
	createReq.Header.Add("Content-Type", "application/json")

	createRes, err := apptest.Test(createReq, -1)
	if err != nil {
		t.Error(err)
	}

	var createdUser types.User
	json.NewDecoder(createRes.Body).Decode(&createdUser)

	return createdUser
}

func TestUserHandler_HandlePostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.dropDown(t)
	apptest := fiber.New()
	UserHandler := NewUserHandler(tdb.User)
	apptest.Post("/", UserHandler.HandlePostUser)
	User := createUserAndReturn(apptest, t, TestFirstName, TestLastName, TestEmail, TestPassWord)
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

func TestUserHandler_HandlePutUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.dropDown(t)
	apptest := fiber.New()
	UserHandler := NewUserHandler(tdb.User)
	apptest.Post("/", UserHandler.HandlePostUser)
	createdUser := createUserAndReturn(apptest, t, TestFirstName, TestLastName, TestEmail, TestPassWord)
	putUrl := fmt.Sprintf("/%s", createdUser.ID.Hex())
	apptest.Put("/:id", UserHandler.HandlePutUser)
	updateUserParams := types.UpdateUserParams{
		FirstName: UpdateFirstName,
		LastName:  UpdateLastName,
	}
	updateUserRequest, _ := json.Marshal(updateUserParams)
	putReq := httptest.NewRequest("PUT", putUrl, bytes.NewReader(updateUserRequest))
	putReq.Header.Add("Content-Type", "application/json")
	putRes, err := apptest.Test(putReq)
	if err != nil {
		t.Error(err)
	}
	if putRes.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status code %d, got %d", fiber.StatusOK, putRes.StatusCode)
	}
	apptest.Get("/:id", UserHandler.HandleGetUser)
	getReq := httptest.NewRequest("GET", putUrl, nil)
	putReq.Header.Add("Content-Type", "application/json")
	getRes, err := apptest.Test(getReq)
	if err != nil {
		t.Error(err)
	}
	var updatedUser types.User
	json.NewDecoder(getRes.Body).Decode(&updatedUser)

	if updatedUser.FirstName != UpdateFirstName {
		t.Errorf("Expected %s but got %s", UpdateFirstName, updatedUser.FirstName)
	}
	if updatedUser.LastName != UpdateLastName {
		t.Errorf("Expected %s but got %s", UpdateLastName, updatedUser.LastName)
	}

}
