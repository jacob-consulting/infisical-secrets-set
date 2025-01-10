MAKEFLAGS += --always-make

aider:
	aider --multiline go.mod main.go
