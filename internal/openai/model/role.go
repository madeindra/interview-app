package model

type Role string

const (
	ROLE_SYSTEM    Role = "system"
	ROLE_ASSISTANT Role = "assistant"
	ROLE_USER      Role = "user"
)
