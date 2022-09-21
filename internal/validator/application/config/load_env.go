package config

import (
	"github.com/brienze1/crypto-robot-validator/pkg/log"
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

const (
	envKey                    = "VALIDATOR_ENV"
	alternateParentFolderName = "crypto-robot-validator"
	configDirPath             = "/config/"
	envFile                   = ".env"
	envFilePath               = ".env."
	development               = "development"
	test                      = "test"
)

// LoadEnv class is responsible for loading order of .env files. Uses godotenv.Load to inject env variables from
// file into environment;
func LoadEnv() {
	env := os.Getenv(envKey)

	if "" == env {
		env = development
	}
	load(envFilePath + env)
	load(envFile)
}

func load(file string) {
	err := godotenv.Load("." + configDirPath + file)
	if err != nil {
		rootPath := getRootPath(alternateParentFolderName)
		err := godotenv.Load(rootPath + configDirPath + file)
		if err != nil {
			log.Logger().Error(err, "failed loading env file "+file)
			panic("Error loading file: " + file)
		}
	}
}

func getRootPath(dirName string) string {
	projectName := regexp.MustCompile(`^(.*` + dirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	return string(projectName.Find([]byte(currentWorkDirectory)))
}

// LoadTestEnv is used in tests to load .env.test file variables.
func LoadTestEnv() {
	_ = os.Setenv(envKey, test)

	LoadEnv()
}
