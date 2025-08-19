package tesla_test

import (
	"log"
	"os"
	"teslatrack/pkg/tesla"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("警告：无法加载 .env 文件!", err)
	}
}

func getAccessToken() string {
	return os.Getenv("TESLA_ACCESS_TOKEN")
}

func TestGetVehices(t *testing.T) {
	tesla.GetVehices(getAccessToken())
}
