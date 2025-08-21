package tesla

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// CN_BASE_URL是中国区的特斯拉车队API的基础URL
const CN_BASE_URL = "https://fleet-api.prd.cn.vn.cloud.tesla.cn"

const (
	// VEHICLES_API 是获取车辆列表的API端点
	VEHICLES_API = CN_BASE_URL + "/api/v1/vehicles" // VEHICLES 返回此账户下车辆的列表。默认页面大小为100。
)

// client 是一个全局的http客户端
var client *http.Client

// init 函数用于初始化http客户端
func init() {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
}

// GetVehices 使用给定的访问令牌获取车辆列表
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

// requestAppendAuthorization 向http请求添加认证头和其他必要的头信息
func requestAppendAuthorization(request *http.Request, accessToken string) {
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Origin", "https://teslatrack.wallora.top")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
}
