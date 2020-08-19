package services

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func SaveAvatarFromUrl(url string) (string, error) {
	var (
		storagePath = os.Getenv("STORAGE_PATH")
		avatarName  = RandomNumber(20)
	)
	avatarFile, err := os.Create(fmt.Sprintf("%s/uploads/avatars/%s.png", storagePath, avatarName))
	if err != nil {
		return avatarName, err
	}
	defer avatarFile.Close()
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return avatarName, err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(avatarFile, resp.Body); err != nil {
		return avatarName, err
	}
	return avatarName, nil
}

func RandomNumber(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("0123456789")
	b := make([]rune, length)
	for i := range b {
		if letter := letters[rand.Intn(len(letters))]; letter != 0 {
			b[i] = letters[rand.Intn(len(letters))]
		}
		continue
	}
	return string(b)
}

func RandomUserName() string {
	return "u" + RandomNumber(10)
}

func GenerateHash() string {
	return RandomString(40)
}
