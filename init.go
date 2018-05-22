package main

import (
	"git-me/models"
)

func InitModels() error {
	if err := models.PrepareUser(); err != nil {
		return err
	}

	if err := models.PrepareTaskInfo(); err != nil {
		return err
	}

	return models.PrepareHistory()
}
