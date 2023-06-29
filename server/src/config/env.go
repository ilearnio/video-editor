package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"videoeditor/src/helpers"
)

var IsDevelopmentEnv bool = false
var IsTestEnv bool = false
var IsProductionEnv bool = false

type envVars struct {
	ENVIRONMENT string

	API_URL string

	PLAYHT_SECRET  string
	PLAYHT_USER_ID string
}

var EnvConfig envVars

//
// Order of loading for .env files.
// The values taken from here https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use
// Note: The order is reversed in slices below, because of the weird godotenv loading approach.
//

var envFiles = []string{
	".env.%s.local",
	".env.local",
	".env.%s",
	".env",
}

var testEnvFiles = []string{
	".env.test.local",
	".env.test",
	".env",
}

func LoadEnv() error {
	detectEnvironmentName()
	return LoadEnvConfig()
}

func GetEnvName() string {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "development" || environment == "" {
		return "development"
	}
	return environment
}

func detectEnvironmentName() {
	if flag.Lookup("test.v") != nil {
		IsTestEnv = true
		return
	}

	envName := GetEnvName()
	IsDevelopmentEnv = envName == "development"
	IsTestEnv = envName == "test"
	IsProductionEnv = envName == "production"
}

func findExistingEnvPaths() []string {
	var envFileNames = make([]string, len(envFiles))
	if IsTestEnv {
		envFileNames = testEnvFiles
	} else {
		envName := GetEnvName()
		for i, envFile := range envFiles {
			if strings.Contains(envFile, "%s") {
				envFile = fmt.Sprintf(envFile, envName)
			}
			envFileNames[i] = envFile
		}
	}

	var envFilePaths []string
	for _, envFile := range envFileNames {
		fullPath := path.Join(CWD_PATH, envFile)
		if helpers.FileExists(fullPath) {
			envFilePaths = append(envFilePaths, fullPath)
		}
	}
	return envFilePaths
}

func LoadEnvConfig() error {
	envFilePaths := findExistingEnvPaths()

	err := godotenv.Load(envFilePaths...)
	if err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}
	envconfig.Process("", &EnvConfig)

	return nil
}
