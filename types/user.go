package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 8
	emailValid      = 1
	emailRegex      = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	firstName       = "firstName"
	lastName        = "lastName"
	email           = "email"
	password        = "password"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type ValidateMap map[string]string

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encryptPassword, err := generateEncryptedPassword(params.Password)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encryptPassword),
	}, nil
}

func (params CreateUserParams) ValidateCreateUserParams() ValidateMap {
	errors := ValidateMap{}
	validateShortHand(errors, firstName, params.FirstName, minFirstNameLen)
	validateShortHand(errors, lastName, params.LastName, minLastNameLen)
	validateShortHand(errors, password, params.Password, minPasswordLen)
	validateShortHand(errors, email, isEmailValid(params.Email), emailValid)
	return errors
}

func (params UpdateUserParams) ValidateUpdateUserParams() ValidateMap {
	errors := ValidateMap{}
	validateShortHand(errors, firstName, params.FirstName, minFirstNameLen)
	validateShortHand(errors, lastName, params.LastName, minLastNameLen)
	return errors
}

func HandleUserErrors(k string, v int) string {
	if k == "email" {
		return fmt.Sprintf("%v not valid", k)
	}
	return fmt.Sprintf("%v must be at least %d character", k, v)
}

func isEmailValid(e string) string {
	valid := "valid"
	notvalid := ""
	emailRegex := regexp.MustCompile(emailRegex)
	if !emailRegex.MatchString(e) {
		return notvalid
	}
	return valid
}

func generateEncryptedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
}

func (params UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(params.FirstName) > 0 {
		m["firstName"] = params.FirstName
	}
	if len(params.LastName) > 0 {
		m["lastName"] = params.LastName
	}
	return m
}

func validateShortHand(e ValidateMap, f, v string, l int) ValidateMap {
	if len(v) < l {
		e[f] = HandleUserErrors(f, l)
	}
	return e
}
