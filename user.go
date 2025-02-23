package main

import (
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
)

type User struct {
	SqliteId     int       `json:"-"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"passwordHash"`
	AccessToken  string    `json:"accessToken"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"createdAt"`
	Level        int       `json:"level"`
}

func (me *User) randomize() {
	var qualityWord = QUALITY_WORDS[rand.IntN(len(QUALITY_WORDS))]
	var animalName = ANIMAL_NAMES[rand.IntN(len(ANIMAL_NAMES))]
	var index = rand.IntN(1000)
	me.Name = qualityWord + " " + animalName + " " + strconv.Itoa(index)
	me.PasswordHash = ""
	for i := 0; i < 2; i++ {
		me.PasswordHash += strconv.Itoa(rand.Int())
	}
	for i := 0; i < 4; i++ {
		me.AccessToken += strconv.Itoa(rand.Int())
	}
	me.Email = qualityWord + "@" + strings.ToLower(animalName) + ".com"
	var createdAgo = time.Duration(rand.IntN(1000)) * time.Minute
	me.CreatedAt = time.Now().Add(-createdAgo).Truncate(time.Second)
	me.Level = rand.IntN(100)
}

func generateRandomUsers(count int) []*User {
	var users = make([]*User, count)
	for i := 0; i < count; i++ {
		var user = new(User)
		user.randomize()
		users[i] = user
	}
	return users
}
