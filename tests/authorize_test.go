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
	callbackAuthorizeUrl = DEFAULT_SERVER_HOST + "/api/v1/authorize/callback"
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
		RedirectURI:  "https://teslatrack.wallora.top/api/v1/authorize/callback",
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

// TestExchangeCode 交换特斯拉回调带过来的Code
func TestExchangeCode(t *testing.T) {
	exchangeUrl := "https://auth.tesla.cn/oauth2/v3/token"

	values := url.Values{
		"grant_type":    []string{"authorization_code"},
		"client_id":     []string{os.Getenv("TESLA_CLIENT_ID")},
		"client_secret": []string{os.Getenv("TESLA_CLIENT_SECRET")},
		"code":          []string{"CN_2749410895cdfd6c743aaa86aabf40b23d814176b33438e859bb7e77afe8"},
		"audience":      []string{"https://fleet-api.prd.cn.vn.cloud.tesla.cn"},
		"redirect_uri":  []string{"https://teslatrack.wallora.top/api/v1/authorize/callback"},
	}
	response, err := http.PostForm(exchangeUrl, values)
	if err != nil {
		panic(err)
	}
	fmt.Printf("statusCode: %d, status:%s \n", response.StatusCode, response.Status)

	responseBody, _ := io.ReadAll(response.Body)
	fmt.Println(string(responseBody))
}

func TestCallback(t *testing.T) {
	id, _ := uuid.NewV7()
	http.Get(callbackAuthorizeUrl + "?code=" + id.String())
}
