package main

import (
	"git-me/models"
)

func Init() error {
	return models.PrepareUser()
}
