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
