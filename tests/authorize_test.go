package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const DEFAULT_SERVER_HOST = "http://127.0.0.1:8100"

var (
	// createAuthorizeUrl
	createAuthorizeUrl   = DEFAULT_SERVER_HOST + "/api/v1/authorize"
	callbaclAuthorizeUrl = DEFAULT_SERVER_HOST + "/api/v1/authorize/callback"
	redirectAuthorizeUrl = DEFAULT_SERVER_HOST + "/api/v1/authorize/redirect"
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

type (
	redirectRequest struct {
		ClientId string `json:"clientId"`
	}
	redirectReply struct {
		Scope                  string `json:"scope"`
		State                  string `json:"state"`
		Nonce                  string `json:"nonce"`
		PromptMissingScopes    bool   `json:"promptMissingScopes"`
		RequireRequestedScopes bool   `json:"requireRequestedScopes"`
		RedirectUri            string `json:"redirectUri"`
	}
)

func TestRedirect(t *testing.T) {
	request := redirectRequest{
		ClientId: os.Getenv("TESLA_CLIENT_ID"),
	}
	body, _ := json.Marshal(request)
	response, err := http.Post(redirectAuthorizeUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	fmt.Printf("statusCode: %d, status:%s \n", response.StatusCode, response.Status)

	responseBody, _ := io.ReadAll(response.Body)

	var reply redirectReply
	_ = json.Unmarshal(responseBody, &reply)
	fmt.Printf("responseBody: %+v \n\n", reply)

	teslaAuthorizeUrl := "https://auth.tesla.cn/oauth2/v3/authorize?&client_id=" + os.Getenv("TESLA_CLIENT_ID") + "&locale=en-US&prompt=login&redirect_uri=" + reply.RedirectUri + "&response_type=code&scope=" + url.QueryEscape(reply.Scope) + "&state=" + reply.State
	fmt.Println(teslaAuthorizeUrl)
}

func TestCallback(t *testing.T) {
	id, _ := uuid.NewV7()
	http.Get(callbaclAuthorizeUrl + "?code=" + id.String())
}
