package settings

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DbConnection string
var DbUsername string
var DbPassword string
var DbHost string
var DbPort string
var DbDatabase string

var RankUtilApi string
var DcWrapperApi string

var Debug bool

var RankTaskApiBaseUrl string

func init() {
	checkEnv()
	LoadSetting()
}

func checkEnv() {
	_ = godotenv.Load()
	needChecks := []string{
		"DB_CONNECTION", "DB_HOST", "DB_PORT", "DB_DATABASE", "DB_USERNAME", "DB_PASSWORD", "RANK_TASK_API_BASE_URL",
		"RANK_UTIL_API", "DC_WRAPPER_API", "RANK_TASK_API_BASE_URL",
	}

	for _, envKey := range needChecks {
		if os.Getenv(envKey) == "" {
			log.Fatalf("env %s missed", envKey)
		}
	}
}

func LoadSetting() {
	DbConnection = os.Getenv("DB_CONNECTION")
	DbUsername = os.Getenv("DB_USERNAME")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbDatabase = os.Getenv("DB_DATABASE")
	RankUtilApi = os.Getenv("RANK_UTIL_API")
	DcWrapperApi = os.Getenv("DC_WRAPPER_API")

	debug := os.Getenv("DEBUG")
	if debug != "" && debug != "false" && debug != "0" {
		Debug = true
	}
	RankTaskApiBaseUrl = os.Getenv("RANK_TASK_API_BASE_URL")
}
