package models

var fakeDB string

func SaveMessage(text string) { fakeDB = text }

func GetMessage() string { return fakeDB }
