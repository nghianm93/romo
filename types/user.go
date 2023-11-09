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

func HandleUserErrors(k string, v int) ValidateMap {
	if k == "email" {
		return ValidateMap{k: fmt.Sprintf("%v not valid", k)}
	}
	return ValidateMap{k: fmt.Sprintf("%v must be at least %d character", k, v)}
}

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

func (params CreateUserParams) ValidateCreateUserParams() []ValidateMap {
	var errors []ValidateMap
	if len(params.FirstName) < minFirstNameLen {
		errors = append(errors, HandleUserErrors(firstName, minFirstNameLen))
	}
	if len(params.LastName) < minLastNameLen {
		errors = append(errors, HandleUserErrors(lastName, minLastNameLen))
	}
	if len(params.Password) < minPasswordLen {
		errors = append(errors, HandleUserErrors(password, minPasswordLen))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, HandleUserErrors(email, minFirstNameLen))
	}
	return errors
}

func (p UpdateUserParams) ValidateUpdateUserParams() []ValidateMap {
	var errors []ValidateMap
	if len(p.FirstName) < minFirstNameLen {
		errors = append(errors, HandleUserErrors(firstName, minFirstNameLen))
	}
	if len(p.LastName) < minLastNameLen {
		errors = append(errors, HandleUserErrors(lastName, minLastNameLen))
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(emailRegex)
	return emailRegex.MatchString(e)
}

func generateEncryptedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}
	return m
}
