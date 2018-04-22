package main

import (
	"git-me/models"
)

func Init() error {
	if err :=  models.PrepareUser();err != nil {
		return err
	}

	return models.PrepareHistory()
}
