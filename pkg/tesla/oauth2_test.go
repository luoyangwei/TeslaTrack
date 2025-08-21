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

func TestGetPartner(t *testing.T) {
	if _, err := tesla.GetPartner(os.Getenv("TESLA_CLIENT_ID"), os.Getenv("TESLA_CLIENT_SECRET")); err != nil {
		panic(err)
	}
	fmt.Println("success")
}

func TestRegisterPartner(t *testing.T) {
	partner, err := tesla.GetPartner(os.Getenv("TESLA_CLIENT_ID"), os.Getenv("TESLA_CLIENT_SECRET"))
	if err != nil {
		panic(err)
	}

	// {"response":{"account_id":"8f066518-7d6f-41b5-a752-cc28c3368558","domain":"teslatrack.wallora.top","name":"TeslaTrack","description":"“特行记”是一款专为特斯拉用户打造的行驶数据记录与交流平台。用户可以便捷地记录和管理自己的车辆行驶数据，包括行程、能耗、驾驶习惯等多维度信息。同时，应用内设有社区功能，方便车主们分享用车体验、交流驾驶心得、获取最新资讯，打造专属特斯拉车主的互动空间。","client_id":"59748905-9613-419e-8685-fd2267ab5757","ca":null,"created_at":"2025-08-21T07:47:43.793Z","updated_at":"2025-08-21T07:59:37.726Z","enterprise_tier":"pay_as_you_go","issuer":null,"csr":null,"csr_updated_at":null,"public_key":"0449154b994ba82752c5f31eb39376a236e9470031fbbb4e4d06a5db11bb99c1d5ab8d3aa84c74312455205355998aa4187ade901ddf39827d5b16930e3abd18be","public_key_hash":"de6c3b382e0782ebd8ab4dd7323a8528"}}
	tesla.RegisterPartner(partner, "teslatrack.wallora.top")
}

func TestGetPartnerPublicKey(t *testing.T) {
	partner, err := tesla.GetPartner(os.Getenv("TESLA_CLIENT_ID"), os.Getenv("TESLA_CLIENT_SECRET"))
	if err != nil {
		panic(err)
	}
	// This account does not have access to teslatrack.wallora.top
	tesla.GetPartnerPublicKey(partner, "teslatrack.wallora.top")
}
