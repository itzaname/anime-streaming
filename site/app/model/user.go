package model

import (
	"encoding/json"
	"fmt"

	"github.com/itzaname/anime-streaming/site/app/lib/database"
)

type User struct {
	ID        string
	BytesUsed int64
	Admin     bool
	Status    uint8
}

// UserExists returns if a user exists in the database
func UserExists(id string) (bool, error) {
	return database.DB.Exists(fmt.Sprintf("user:id:%s", id)).Result()
}

// UserById returns a user object
func UserById(id string) (User, error) {
	var user User
	resp, err := database.DB.Get(fmt.Sprintf("user:id:%s", id)).Result()
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(resp), &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UserAdd will add a user
func UserAdd(name string, admin bool) error {
	var user User
	user.ID = name
	user.Admin = admin

	if err := user.Update(); err != nil {
		return err
	}

	return database.DB.LPush("user:id", name).Err()
}

// Key returns the database key for the object
func (u *User) Key() string {
	return fmt.Sprintf("user:id:%s", u.ID)
}

// Update will update the current user object in the redis store
func (u *User) Update() error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return database.DB.Set(u.Key(), data, -1).Err()
}
