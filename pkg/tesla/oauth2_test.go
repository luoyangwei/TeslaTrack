package tesla_test

import (
	"fmt"
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

func TestGetPartners(t *testing.T) {
	if err := tesla.GetPartners(os.Getenv("TESLA_CLIENT_ID"), os.Getenv("TESLA_CLIENT_SECRET")); err != nil {
		panic(err)
	}
	fmt.Println("success")
}
