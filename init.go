package main

import (
	"git-me/models"
)

func InitModels() error {
	if err := models.PrepareUser(); err != nil {
		return err
	}

	if err := models.PrepateTaskInfo(); err != nil {
		return err
	}

	return models.PrepareHistory()
}
