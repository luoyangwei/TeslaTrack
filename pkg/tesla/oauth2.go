package tesla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	GRANT_TYPE = "client_credentials"
	// audience
	AUDIENCE = "https://fleet-api.prd.cn.vn.cloud.tesla.cn"
	// scope
	SCOPE = "openid user_data vehicle_device_data vehicle_location vehicle_cmds vehicle_charging_cmds energy_device_data energy_cmds"
)

const (
	OAUTH_TO_TOKEN_URL = "https://auth.tesla.cn/oauth2/v3/token"
	// https://fleet-api.prd.cn.vn.cloud.tesla.cn/api/1/partner_accounts
	FLEET_API_PARTNER_ACCOUNTS_URL   = "https://fleet-api.prd.cn.vn.cloud.tesla.cn/api/1/partner_accounts"
	FLEET_API_PARTNER_PUBLIC_KEY_URL = "https://fleet-api.prd.cn.vn.cloud.tesla.cn/api/1/partner_accounts/public_key?domain=%s"
)

type Partner struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func GetPartner(clientID, clientSecret string) (*Partner, error) {
	values := url.Values{
		"grant_type":    []string{GRANT_TYPE},
		"client_id":     []string{clientID},
		"client_secret": []string{clientSecret},
		"audience":      []string{AUDIENCE},
		"scope":         []string{SCOPE},
	}
	response, err := http.PostForm(OAUTH_TO_TOKEN_URL, values)
	if err != nil {
		return nil, err
	}
	bytes, _ := io.ReadAll(response.Body)

	var partner Partner
	_ = json.Unmarshal(bytes, &partner)

	return &partner, nil
}

type (
	RegisterPartnerRequest struct {
		Domain string `json:"domain"`
	}
)

func RegisterPartner(partner *Partner, domain string) {
	registerPartnerRequest := RegisterPartnerRequest{Domain: domain}
	body, _ := json.Marshal(registerPartnerRequest)

	request, _ := http.NewRequest(http.MethodPost, FLEET_API_PARTNER_ACCOUNTS_URL, bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", partner.TokenType+" "+partner.AccessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	fmt.Println(string(responseBody))
}

func GetPartnerPublicKey(partner *Partner, domain string) {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(FLEET_API_PARTNER_PUBLIC_KEY_URL, domain), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", partner.TokenType+" "+partner.AccessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	responseBody, _ := io.ReadAll(response.Body)
	fmt.Println(string(responseBody))
}
