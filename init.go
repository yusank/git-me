package main

import (
	"git-me/model"
)

func Init() error {
	return model.PrepareUser()
}
