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

	// {
	// 	"access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjRWZjY5Z1lhZzBWMnZNVDdTNEhOMjN6VURudyJ9.eyJpc3MiOiJodHRwczovL2F1dGgudGVzbGEuY24vb2F1dGgyL3YzL250cyIsImF6cCI6IjU5NzQ4OTA1LTk2MTMtNDE5ZS04Njg1LWZkMjI2N2FiNTc1NyIsInN1YiI6IjY2OTYzZjI0LTNhNzctNDJhMi04ZGFjLTk0NjVlNzJkOTA3MSIsImF1ZCI6WyJodHRwczovL2ZsZWV0LWFwaS5wcmQuY24udm4uY2xvdWQudGVzbGEuY24iLCJodHRwczovL2F1dGgudGVzbGEuY24vb2F1dGgyL3YzL3VzZXJpbmZvIl0sInNjcCI6WyJvcGVuaWQiLCJvZmZsaW5lX2FjY2VzcyIsInVzZXJfZGF0YSIsInZlaGljbGVfZGV2aWNlX2RhdGEiLCJ2ZWhpY2xlX2xvY2F0aW9uIiwidmVoaWNsZV9jbWRzIiwidmVoaWNsZV9jaGFyZ2luZ19jbWRzIiwiZW5lcmd5X2RldmljZV9kYXRhIiwiZW5lcmd5X2NtZHMiXSwiYW1yIjpbInNtc290cCJdLCJleHAiOjE3NTU0NTAxNjUsImlhdCI6MTc1NTQyMTM2NSwib3VfY29kZSI6IkNOIiwibG9jYWxlIjoiemgtQ04iLCJhY2NvdW50X3R5cGUiOiJwZXJzb24iLCJvcGVuX3NvdXJjZSI6ZmFsc2UsImFjY291bnRfaWQiOiI4ZjA2NjUxOC03ZDZmLTQxYjUtYTc1Mi1jYzI4YzMzNjg1NTgiLCJhdXRoX3RpbWUiOjE3NTU0MjEzNjR9.rE7sRDYszQcXJHz4yLo6AmTMPqTb6U7j3hh16VozEXfWr9bEPKlhxzwU7tQL04HtF-RIoRSFpQPPk1k9QOWcF98KfmPWwuEEUyMKcCe77mu6jDhf24S53IAZuCEJKpC8jA8_1dcRqnfg1PYZS_ngM3NVgxAA2-bX5qOI9mJj9kY4RH99iY-2HHgzpKfAT_FdyHsR5J_2R01wFl7urfebGB2zmlgHf2Rjo1xzRWbMiuZoYRVJ2QQUcac069lWuf7You-HwBHn_oG3Li_FwLX-J8I9_89YggjU59zGxhNS4uJpEzmKbGkU4mTCRwa8pb_INx05JWNFNIRf1vBT_T7eVg",
	// 	"refresh_token": "CN_87cc2b76dd4406b78666b2f79c2d19f9eabe30811aa32d015d6cb139c69805f8",
	// 	"id_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjRWZjY5Z1lhZzBWMnZNVDdTNEhOMjN6VURudyJ9.eyJpc3MiOiJodHRwczovL2F1dGgudGVzbGEuY24vb2F1dGgyL3YzL250cyIsImF1ZCI6IjU5NzQ4OTA1LTk2MTMtNDE5ZS04Njg1LWZkMjI2N2FiNTc1NyIsInN1YiI6IjY2OTYzZjI0LTNhNzctNDJhMi04ZGFjLTk0NjVlNzJkOTA3MSIsImV4cCI6MTc1NTQ1MDE2NSwiaWF0IjoxNzU1NDIxMzY1LCJhdXRoX3RpbWUiOjE3NTU0MjEzNjQsImFtciI6WyJzbXNvdHAiXSwidXBkYXRlZF9hdCI6MTc1NDAxODQ3OX0.Yj6iIZZlDwBO1LLuOQ2zUMFe0xw10-js0H80J55Cu_nkwazLd0o32IdwDQGaaGuB7azoNDO2S1c5cUZULtxKbKa9mGilGEeUoWKEU2zHvmU90_d6Ci-0wRh5CicYhoxTn_YLzm9vSo9PM7nawNgkKROFHOkU4Wi5qATJUot-yjod-MMlxwxk_DwSQH94HiGaqgI9mCm3RIco_1aHUwIvyqbPLfGBCPf8YJsKMS6SEox82sGEaB5_Ads4jK4ZT1bTvL48NDyqi6Ux22nlqlq22MrRbpPR-2C-0loOs7nIXyzneRP7-HqR7-ZOZt7Q1L1NxMdq2ZBrFBhTHP3S8wjCtg",
	// 	"expires_in": 28800,
	// 	"state": "0198b743-50d5-79ee-ae9e-9ff8e639ba72",
	// 	"token_type": "Bearer"
	// }
	responseBody, _ := io.ReadAll(response.Body)
	fmt.Println(string(responseBody))
}

func TestCallback(t *testing.T) {
	id, _ := uuid.NewV7()
	http.Get(callbackAuthorizeUrl + "?code=" + id.String())
}
