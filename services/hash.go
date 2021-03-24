package services

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/CastyLab/grpc.server/storage"
	"github.com/minio/minio-go"
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
	avatarName := RandomNumber(20)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return avatarName, err
	}
	defer resp.Body.Close()
	if _, err = storage.Client.PutObject("avatars", fmt.Sprintf("%s.png", avatarName), resp.Body, -1, minio.PutObjectOptions{}); err != nil {
		return "", err
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
	return fmt.Sprintf("u%s", RandomNumber(10))
}

func GenerateHash() string {
	return RandomString(40)
}
