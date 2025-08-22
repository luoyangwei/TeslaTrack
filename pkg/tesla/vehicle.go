package tesla

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CN_BASE_URL is the base URL for Tesla's fleet API in the China region.
const CN_BASE_URL = "https://fleet-api.prd.cn.vn.cloud.tesla.cn"

const (
	// VEHICLES_API is the API endpoint for fetching the list of vehicles.
	// It returns a list of vehicles under this account. The default page size is 100.
	VEHICLES_API     = CN_BASE_URL + "/api/1/vehicles"
	VEHICLE_DATA_API = CN_BASE_URL + "/api/1/vehicles/%s/vehicle_data"
)

// client is a global http client.
var client *http.Client

// init function initializes the http client.
func init() {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
}

// Response is a generic struct for handling Tesla API responses.
type Response[T any] struct {
	Response         T                 `json:"response,omitempty"`
	Error            string            `json:"error,omitempty"`
	ErrorDescription string            `json:"error_description,omitempty"`
	Messages         map[string]string `json:"messages,omitempty"`
}

// GranularAccess corresponds to the "granular_access" object in the JSON, representing fine-grained access settings.
type GranularAccess struct {
	// HidePrivate indicates whether to hide private information.
	HidePrivate bool `json:"hide_private"`
}

// VehicleData is the complete struct for parsing detailed vehicle information.
type VehicleData struct {
	// ID is the unique numerical identifier for the vehicle.
	ID int64 `json:"id"`

	// VehicleID is an internal vehicle ID.
	VehicleID int64 `json:"vehicle_id"`

	// VIN is the Vehicle Identification Number of the vehicle.
	VIN string `json:"vin"`

	// DisplayName is the custom display name set by the user for the vehicle.
	DisplayName string `json:"display_name"`

	// AccessType indicates the current user's access permission type, e.g., "OWNER".
	AccessType string `json:"access_type"`

	// State represents the current state of the vehicle, e.g., "offline", "online".
	State string `json:"state"`

	// IDs is the string representation of the vehicle ID.
	IDS string `json:"id_s"`

	// APIVersion is the API version used for vehicle communication.
	APIVersion int `json:"api_version"`

	// InService indicates whether the vehicle is in a maintenance or service state.
	InService bool `json:"in_service"`

	// CalendarEnabled indicates whether the vehicle's calendar function is enabled.
	CalendarEnabled bool `json:"calendar_enabled"`

	// BleAutopairEnrolled indicates whether the vehicle is enrolled in the Bluetooth auto-pair feature.
	BleAutopairEnrolled bool `json:"ble_autopair_enrolled"`

	// GranularAccess contains the detailed settings for fine-grained access permissions.
	GranularAccess GranularAccess `json:"granular_access"`

	// Color is the color of the vehicle. Defined as a pointer type to be nullable, as it can be null in the original JSON.
	Color *string `json:"color,omitempty"`

	// OptionCodes is a list of optional configuration codes for the vehicle. Defined as a pointer type to be nullable.
	OptionCodes *[]string `json:"option_codes,omitempty"`

	// Tokens is a list of tokens related to vehicle access. Defined as a pointer type to be nullable.
	Tokens *[]string `json:"tokens,omitempty"`

	// BackseatToken is the token for the backseat entertainment system. Defined as a pointer type to be nullable.
	BackseatToken *string `json:"backseat_token,omitempty"`

	// BackseatTokenUpdatedAt is the last update time of the backseat token. Defined as a pointer type to be nullable.
	BackseatTokenUpdatedAt *string `json:"backseat_token_updated_at,omitempty"`
}

// GetVehices fetches the list of vehicles using the given access token.
func GetVehices(accessToken string) ([]VehicleData, error) {
	// Create a new HTTP GET request to the VEHICLES_API endpoint.
	request, err := http.NewRequest(http.MethodGet, VEHICLES_API, nil)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("must be to newRequest"))
	}
	// Add authorization and other necessary headers to the request.
	requestAppendAuthorization(request, accessToken)

	// Execute the HTTP request.
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("response error"))
	}
	defer response.Body.Close()

	// Read the entire response body.
	bytes, _ := io.ReadAll(response.Body)
	// For debugging purposes, print the raw response body.
	fmt.Println(string(bytes))

	// Unmarshal the JSON response into our generic Response struct containing a slice of VehicleData.
	var data Response[[]VehicleData]
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, errors.Join(err, fmt.Errorf("unmarshal response bytes error"))
	}

	// Check if the API response contains an error.
	if data.Error != "" {
		return nil, fmt.Errorf("tesla response error %s:%s", data.Error, data.ErrorDescription)
	}

	// Return the slice of vehicles from the response.
	return data.Response, nil
}

// GetVehiceData fetches detailed data for a specific vehicle identified by its VIN.
// The response is a complex JSON object containing real-time vehicle status.
// The structure of the `response` object is detailed below:
//
// --- Core Vehicle Information (核心车辆信息) ---
// Basic identification information for the vehicle and owner. (这部分包含了车辆和车主的基本识别信息。)
//
//	id: (e.g., 100021) The unique ID of the vehicle in the Tesla API. (车辆在特斯拉API中的唯一ID。)
//	user_id: (e.g., 800001) The ID of the vehicle's owner. (车辆所属用户的ID。)
//	vehicle_id: (e.g., 99999) Another unique identifier for the vehicle. (车辆的另一个唯一标识符。)
//	vin: (e.g., "TEST00000000VIN01") The Vehicle Identification Number. (车辆识别码 (VIN)。)
//	color: (e.g., null) The color of the vehicle. (车辆颜色。)
//	access_type: (e.g., "OWNER") The access permission type, here it is the owner. (访问权限类型，此处为车主。)
//	tokens: (e.g., ["4f993c5b9e2b937b", "7a3153b1bbb48a96"]) Tokens used for API authentication. (用于API认证的令牌。)
//	state: (e.g., "online") The current network state of the vehicle. (车辆当前的网络状态。)
//	in_service: (e.g., false) Whether the vehicle is currently in service/maintenance. (车辆是否正在维修。)
//	id_s: (e.g., "100021") The string format of the ID. (ID的字符串格式。)
//	calendar_enabled: (e.g., true) Whether the calendar function is enabled. (是否启用了日历功能。)
//	api_version: (e.g., 54) The API version currently used by the vehicle. (车辆当前使用的API版本。)
//
// --- charge_state (充电状态) ---
// Detailed information about the battery and charging. (这部分详细描述了电池和充电相关的所有信息。)
//
//	battery_heater_on: (e.g., false) Whether the battery heater is on. (电池加热器是否开启。)
//	battery_level: (e.g., 42) The current battery level percentage. (当前电池电量百分比。)
//	battery_range: (e.g., 133.99) The current estimated driving range (in miles). (当前预估续航里程（英里）。)
//	charge_amps: (e.g., 48) The current charging amperage. (当前充电电流（安培）。)
//	charge_limit_soc: (e.g., 90) The configured charge limit percentage. (设定的充电上限百分比。)
//	charge_miles_added_rated: (e.g., 202) The rated range added in this charging session. (本次充电增加的额定续航里程。)
//	charge_port_door_open: (e.g., false) Whether the charge port door is open. (充电口门是否打开。)
//	charge_port_latch: (e.g., "Engaged") Whether the charge port is latched. (充电口是否锁定。)
//	charge_rate: (e.g., 0) The current charging rate (in miles/hour). (当前充电速率（英里/小时）。)
//	charger_power: (e.g., 0) The charger power (in kW). (充电器功率（千瓦）。)
//	charger_voltage: (e.g., 2) The charger voltage. (充电器电压。)
//	charging_state: (e.g., "Disconnected") The state of charging (e.g., Charging, Complete, Disconnected). (充电状态，例如：充电中、已完成、已断开。)
//	est_battery_range: (e.g., 143.88) The estimated battery range (in miles). (预估的电池续航里程（英里）。)
//	fast_charger_present: (e.g., false) Whether a fast charger is connected. (是否连接了快充。)
//	minutes_to_full_charge: (e.g., 0) The minutes remaining to full charge. (距离充满还需的分钟数。)
//	timestamp: (e.g., 1692141038420) The timestamp of the data update. (数据更新的时间戳。)
//	usable_battery_level: (e.g., 42) The usable battery level percentage. (可用的电池电量百分比。)
//
// --- climate_state (空调状态) ---
// Status of the vehicle's internal environmental control system. (这部分包含了车辆内部环境控制系统的状态。)
//
//	allow_cabin_overheat_protection: (e.g., true) Whether cabin overheat protection is allowed. (是否允许座舱过热保护。)
//	cabin_overheat_protection: (e.g., "On") The status of cabin overheat protection. (座舱过热保护状态。)
//	driver_temp_setting: (e.g., 21) The temperature setting for the driver's side (in Celsius). (驾驶员侧设定的温度（摄氏度）。)
//	fan_status: (e.g., 0) The fan status/speed. (风扇状态/风速。)
//	inside_temp: (e.g., 38.4) The interior temperature (in Celsius). (车内温度（摄氏度）。)
//	is_auto_conditioning_on: (e.g., true) Whether auto-conditioning is on. (是否开启自动空调。)
//	is_climate_on: (e.g., false) Whether the climate control is on. (空调是否开启。)
//	is_preconditioning: (e.g., false) Whether preconditioning is in progress. (是否正在进行预处理（提前调节温度）。)
//	outside_temp: (e.g., 36.5) The exterior temperature (in Celsius). (车外温度（摄氏度）。)
//	passenger_temp_setting: (e.g., 21) The temperature setting for the passenger's side (in Celsius). (乘客侧设定的温度（摄氏度）。)
//	seat_heater_left: (e.g., 0) The heating level for the left seat (0-3). (左侧座椅加热等级 (0-3)。)
//	steering_wheel_heater: (e.g., false) Whether the steering wheel heater is on. (方向盘加热器是否开启。)
//	timestamp: (e.g., 1692141038419) The timestamp of the data update. (数据更新的时间戳。)
//
// --- drive_state (驾驶状态) ---
// Geographic location and driving status of the vehicle. (包含了车辆的地理位置和行驶状态。)
//
//	gps_as_of: (e.g., 1692137422) The timestamp of the GPS data update. (GPS数据更新的时间戳。)
//	heading: (e.g., 289) The vehicle's heading (0-359 degrees, 0 is North). (车辆朝向（0-359度，0为正北）。)
//	latitude: (e.g., 37.7765494) The latitude. (纬度。)
//	longitude: (e.g., -122.4195418) The longitude. (经度。)
//	power: (e.g., 1) The current motor power (in kW). (当前电机功率（千瓦）。)
//	shift_state: (e.g., null) The gear state (P, R, N, D). (档位状态 (P, R, N, D)。)
//	speed: (e.g., null) The current speed (in mph). (当前车速（英里/小时）。)
//	timestamp: (e.g., 1692141038420) The timestamp of the data update. (数据更新的时间戳。)
//
// --- gui_settings (界面设置) ---
// User preferences on the in-car display. (用户在车载屏幕上的偏好设置。)
//
//	gui_24_hour_time: (e.g., false) Whether to use 24-hour format for time display. (是否使用24小时制显示时间。)
//	gui_charge_rate_units: (e.g., "mi/hr") The unit for charge rate. (充电速率单位。)
//	gui_distance_units: (e.g., "mi/hr") The unit for distance. (距离单位。)
//	gui_range_display: (e.g., "Rated") The display mode for range (Rated or Typical). (续航显示模式（额定或典型）。)
//	gui_temperature_units: (e.g., "F") The unit for temperature (Fahrenheit). (温度单位（华氏度）。)
//	timestamp: (e.g., 1692141038420) The timestamp of the data update. (数据更新的时间戳。)
//
// --- vehicle_config (车辆配置) ---
// Hardware and software configuration of the vehicle from the factory. (车辆出厂时的硬件和软件配置信息。)
//
//	can_actuate_trunks: (e.g., true) Whether remote opening/closing of trunks is supported. (是否支持远程开关前后备箱。)
//	car_type: (e.g., "modely") The car model. (车型。)
//	charge_port_type: (e.g., "US") The type of charge port. (充电口类型。)
//	driver_assist: (e.g., "TeslaAP3") The version of the driver assistance system. (驾驶辅助系统版本。)
//	exterior_color: (e.g., "MidnightSilver") The exterior color. (外观颜色。)
//	has_ludicrous_mode: (e.g., false) Whether the vehicle has Ludicrous Mode. (是否配备狂暴模式。)
//	interior_trim_type: (e.g., "Black2") The type of interior trim. (内饰类型。)
//	motorized_charge_port: (e.g., true) Whether it has a motorized charge port. (是否为电动充电口。)
//	rhd: (e.g., false) Whether it is a right-hand drive vehicle. (是否为右舵驾驶。)
//	roof_color: (e.g., "RoofColorGlass") The type of roof. (车顶类型。)
//	wheel_type: (e.g., "Apollo19") The type of wheels. (轮毂类型。)
//
// --- vehicle_state (车辆状态) ---
// A comprehensive set of real-time vehicle states. (综合了车辆的各种实时状态。)
//
//	car_version: (e.g., "2023.7.20 7910d26d5c64") The vehicle's software version. (车辆软件版本。)
//	center_display_state: (e.g., 0) The state of the center display (0 is off). (中控屏幕状态（0为关闭）。)
//	df, dr, pf, pr: (e.g., 0) The status of the driver front/rear and passenger front/rear doors (0 is closed). (驾驶员前门、后门，乘客前门、后门状态 (0为关闭)。)
//	ft, rt: (e.g., 0) The status of the front and rear trunks (0 is closed). (前备箱、后备箱状态 (0为关闭)。)
//	homelink_nearby: (e.g., false) Whether the vehicle is within Homelink range. (是否在 Homelink 范围内。)
//	is_user_present: (e.g., false) Whether a user is present in the vehicle. (车内是否有人。)
//	locked: (e.g., true) Whether the vehicle is locked. (车辆是否上锁。)
//	media_info: (e.g., {...}) Information about media playback. (媒体播放信息。)
//	odometer: (e.g., 15720.074889) The total mileage (in miles). (总行驶里程（英里）。)
//	sentry_mode: (e.g., false) Whether Sentry Mode is enabled. (哨兵模式是否开启。)
//	software_update: (e.g., {...}) The status of software updates. (软件更新状态。)
//	speed_limit_mode: (e.g., {...}) Information about speed limit mode. (限速模式信息。)
//	tpms_pressure_fl, fr, rl, rr: (e.g., 3.1, 3.1, 3.15, 3) The tire pressure for each tire. (各轮胎胎压。)
//	valet_mode: (e.g., false) Whether Valet Mode is enabled. (代客模式是否开启。)
//	vehicle_name: (e.g., "grADOFIN") The custom name of the vehicle. (车辆自定义名称。)
func GetVehiceData(accessToken, vin string) (any, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(VEHICLE_DATA_API, vin), nil)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("must be to newRequest"))
	}
	// Add authorization and other necessary headers to the request.
	requestAppendAuthorization(request, accessToken)

	// Execute the HTTP request.
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("response error"))
	}
	defer response.Body.Close()

	// Read the entire response body.
	bytes, _ := io.ReadAll(response.Body)
	// For debugging purposes, print the raw response body.
	fmt.Println(string(bytes))

	return nil, nil
}

// requestAppendAuthorization adds authorization and other necessary headers to an http request.
func requestAppendAuthorization(request *http.Request, accessToken string) {
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Origin", "https://teslatrack.wallora.top")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
}
