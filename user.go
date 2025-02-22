package main

import (
	"math/rand/v2"
	"strconv"
	"time"
)

type User struct {
	Name         string    `json:"name"`
	PasswordHash string    `json:"passwordHash"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"createdAt"`
	Level        int       `json:"level"`
}

func (this *User) randomize() {
	var qualityWord = QUALITY_WORDS[rand.IntN(len(QUALITY_WORDS))]
	var animalName = ANIMAL_NAMES[rand.IntN(len(ANIMAL_NAMES))]
	var index = rand.IntN(1000)
	this.Name = qualityWord + " " + animalName + " " + strconv.Itoa(index)
	this.PasswordHash = ""
	for i := 0; i < 4; i++ {
		this.PasswordHash += strconv.Itoa(rand.Int())
	}
	this.Email = strconv.Itoa(rand.Int()) + "@" + strconv.Itoa(rand.Int()) + ".com"
	this.CreatedAt = time.Now()
	this.Level = rand.IntN(100)
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
