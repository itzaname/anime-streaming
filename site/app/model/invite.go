package model

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/itzaname/anime-streaming/site/app/lib/database"
)

// InviteValid checks if an invite is valid
func InviteValid(invite string) (bool, error) {
	return database.DB.Exists(fmt.Sprintf("invite:%s", invite)).Result()
}

// InviteCreate creates an invite
func InviteCreate() (string, error) {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	invite := hex.EncodeToString(b)

	if err := database.DB.Set(fmt.Sprintf("invite:%s", invite), true, -1).Err(); err != nil {
		return "", err
	}

	return invite, database.DB.LPush("invite:id", invite).Err()
}

// InviteDelete deletes an invite from the database
func InviteDelete(invite string) error {
	if err := database.DB.Del(fmt.Sprintf("invite:%s", invite)).Err(); err != nil {
		return err
	}

	return database.DB.LRem("invite:id", -1, invite).Err()
}

// InviteList returns string array of all invites
func InviteList() ([]string, error) {
	return database.DB.LRange("invite:id", 0, -1).Result()
}
