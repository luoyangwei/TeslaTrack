package tesla

import (
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

	OAUTH_TO_TOKEN_URL = "https://auth.tesla.cn/oauth2/v3/token"
)

type Partners struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ToeknType   string `json:"toekn_type"`
}

func GetPartners(clientID, clientSecret string) error {
	values := url.Values{
		"grant_type":    []string{GRANT_TYPE},
		"client_id":     []string{clientID},
		"client_secret": []string{clientSecret},
		"audience":      []string{AUDIENCE},
		"scope":         []string{SCOPE},
	}
	response, err := http.PostForm(OAUTH_TO_TOKEN_URL, values)
	if err != nil {
		return err
	}
	bytes, _ := io.ReadAll(response.Body)

	var partners Partners
	_ = json.Unmarshal(bytes, &partners)

	fmt.Println(partners)
	return nil
}
