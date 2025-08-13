package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

const DEFAULT_SERVER_HOST = "http://127.0.0.1:8000"

var (
	// createAuthorizeUrl
	createAuthorizeUrl = DEFAULT_SERVER_HOST + "/api/v1/authorize"
)

type createAuthorize struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	GrantType    string `json:"grantType"`
	RedirectURI  string `json:"redirectURI"`
}

func init() {
	if err := godotenv.Load(".env.test"); err != nil {
		log.Println("警告：无法加载 .env 文件!", err)
	}
}

func TestCreateAuthorize(t *testing.T) {
	body := createAuthorize{
		ClientId:     os.Getenv("TESLA_CLIENT_ID"),
		ClientSecret: os.Getenv("TESLA_CLIENT_SECRET"),
		GrantType:    "authorization_code",
		RedirectURI:  "https://teslatrack.luoyangwei.cn/api/v1/authorize/callback",
	}

	request, _ := json.Marshal(body)
	response, err := http.Post(createAuthorizeUrl, "application/json", bytes.NewReader(request))
	if err != nil {
		panic(err)
	}

	fmt.Printf("statusCode: %d, status:%s \n", response.StatusCode, response.Status)
}
