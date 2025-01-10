MAKEFLAGS += --always-make

aider:
	aider --multiline go.mod main.go

actionlint:
	actionlint .github/workflows/go.yml
