package main

import (
	"math/rand/v2"
	"strconv"
	"time"
)

type User struct {
	name         string
	passwordHash string
	email        string
	createdAt    time.Time
	level        int
}

func (this *User) randomize() {
	var qualityWord = QUALITY_WORDS[rand.IntN(len(QUALITY_WORDS))]
	var animalName = ANIMAL_NAMES[rand.IntN(len(ANIMAL_NAMES))]
	var index = rand.IntN(1000)
	this.name = qualityWord + " " + animalName + " " + strconv.Itoa(index)
	this.passwordHash = ""
	for i := 0; i < 4; i++ {
		this.passwordHash += strconv.Itoa(rand.Int())
	}
	this.email = strconv.Itoa(rand.Int()) + "@" + strconv.Itoa(rand.Int()) + ".com"
	this.createdAt = time.Now()
	this.level = rand.IntN(100)
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
