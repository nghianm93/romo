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
	UpdateLastName  = "Update LastName"
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
	getReq.Header.Add("Content-Type", "application/json")
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

func TestUserHandler_HandleGetUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.dropDown(t)
	apptest := fiber.New()
	UserHandler := NewUserHandler(tdb.User)
	apptest.Post("/", UserHandler.HandlePostUser)
	createdUser := createUserAndReturn(apptest, t, TestFirstName, TestLastName, TestEmail, TestPassWord)
	id := createdUser.ID.Hex()
	apptest.Get("/:id", UserHandler.HandleGetUser)
	getUrl := fmt.Sprintf("/%s", id)
	getReq := httptest.NewRequest("GET", getUrl, nil)
	getReq.Header.Add("Content-Type", "application/json")
	getRes, err := apptest.Test(getReq)
	if err != nil {
		t.Error(err)
	}
	var User types.User
	json.NewDecoder(getRes.Body).Decode(&User)

	if User.FirstName != TestFirstName {
		t.Errorf("Expected %s but got %s", TestFirstName, User.FirstName)
	}
	if User.LastName != TestLastName {
		t.Errorf("Expected %s but got %s", TestLastName, User.LastName)
	}
	if User.Email != TestEmail {
		t.Errorf("Expected %s but got %s", TestEmail, User.Email)
	}
	if len(User.EncryptedPassword) > 0 {
		t.Error("Dont return password")
	}
}

func TestUserHandler_HandleGetUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.dropDown(t)
	UserHandler := NewUserHandler(tdb.User)
	apptest := fiber.New()
	apptest.Post("/", UserHandler.HandlePostUser)
	numUsersToCreate := 2
	for i := 0; i < numUsersToCreate; i++ {
		createUserAndReturn(apptest, t, fmt.Sprint(TestFirstName, i+1), fmt.Sprint(TestLastName, i+1), TestEmail, TestPassWord)
	}
	apptest.Get("/users", UserHandler.HandleGetUsers)
	getUsersReq := httptest.NewRequest("GET", "/users", nil)
	getUsersReq.Header.Add("Content-Type", "application/json")
	getUsersRes, err := apptest.Test(getUsersReq)
	if err != nil {
		t.Error(err)
	}
	var Users []types.User
	json.NewDecoder(getUsersRes.Body).Decode(&Users)
	if len(Users) < numUsersToCreate {
		t.Errorf("Something wrong, expected %d users but got %d", numUsersToCreate, len(Users))
	}
	if Users[0].FirstName != fmt.Sprint(TestFirstName, 1) {
		t.Errorf("Something wrong, expected %s but got %s", Users[0].FirstName, fmt.Sprint(TestFirstName, 1))
	}
}

func TestUserHandler_HandleDeleteUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.dropDown(t)
	UserHandler := NewUserHandler(tdb.User)
	apptest := fiber.New()
	apptest.Post("/", UserHandler.HandlePostUser)
	createdUser := createUserAndReturn(apptest, t, TestFirstName, TestLastName, TestEmail, TestPassWord)
	deleteUrl := fmt.Sprintf("/%s", createdUser.ID.Hex())
	apptest.Delete("/:id", UserHandler.HandleDeleteUser)
	deleteUserReq := httptest.NewRequest("DELETE", deleteUrl, nil)
	deleteUserReq.Header.Add("Content-Type", "application/json")
	deleteRes, err := apptest.Test(deleteUserReq)
	if err != nil {
		t.Error(err)
	}
	if deleteRes.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status code %d, got %d", fiber.StatusOK, deleteRes.StatusCode)
	}
}
