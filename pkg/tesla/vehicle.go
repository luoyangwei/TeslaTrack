package tesla

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const CN_BASE_URL = "https://fleet-api.prd.cn.vn.cloud.tesla.cn"

const (
	VEHICLES_API = CN_BASE_URL + "/api/v1/vehicles" // VEHICLES 返回此账户下车辆的列表。默认页面大小为100。
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
}

func GetVehices(accessToken string) {
	request, err := http.NewRequest(http.MethodGet, VEHICLES_API, nil)
	if err != nil {
		panic(err)
	}
	requestAppendAuthorization(request, accessToken)

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	bytes, _ := io.ReadAll(response.Body)
	fmt.Println(string(bytes))
}

func requestAppendAuthorization(request *http.Request, accessToken string) {
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Origin", "https://teslatrack.wallora.top")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
}
