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

// Vehicle is the complete struct for parsing detailed vehicle information.
type Vehicle struct {
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
func GetVehices(accessToken string) ([]Vehicle, error) {
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
	var data Response[[]Vehicle]
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

// ChargeState contains information about the vehicle's charging status.
type ChargeState struct {
	// BatteryHeaterOn indicates if the battery heater is active.
	BatteryHeaterOn bool `json:"battery_heater_on"`
	// BatteryLevel is the current state of charge of the battery in percent.
	BatteryLevel int `json:"battery_level"`
	// BatteryRange is the estimated remaining driving range in the vehicle's configured units (miles/km).
	BatteryRange float64 `json:"battery_range"`
	// ChargeAmps is the current amperage set for charging.
	ChargeAmps int `json:"charge_amps"`
	// ChargeCurrentRequest is the requested charge current.
	ChargeCurrentRequest int `json:"charge_current_request"`
	// ChargeCurrentRequestMax is the maximum charge current that can be requested.
	ChargeCurrentRequestMax int `json:"charge_current_request_max"`
	// ChargeEnableRequest indicates if a charging request is active.
	ChargeEnableRequest bool `json:"charge_enable_request"`
	// ChargeEnergyAdded is the energy added to the battery during the current charging session in kWh.
	ChargeEnergyAdded float64 `json:"charge_energy_added"`
	// ChargeLimitSoc is the user-defined state of charge limit for charging.
	ChargeLimitSoc int `json:"charge_limit_soc"`
	// ChargeLimitSocMax is the maximum possible state of charge limit.
	ChargeLimitSocMax int `json:"charge_limit_soc_max"`
	// ChargeLimitSocMin is the minimum possible state of charge limit.
	ChargeLimitSocMin int `json:"charge_limit_soc_min"`
	// ChargeLimitSocStd is the standard recommended state of charge limit.
	ChargeLimitSocStd int `json:"charge_limit_soc_std"`
	// ChargeMilesAddedIdeal is the ideal range added during the current charging session.
	ChargeMilesAddedIdeal float64 `json:"charge_miles_added_ideal"`
	// ChargeMilesAddedRated is the rated (EPA) range added during the current charging session.
	ChargeMilesAddedRated float64 `json:"charge_miles_added_rated"`
	// ChargePortColdWeatherMode indicates if the charge port is in a special mode for cold weather.
	ChargePortColdWeatherMode bool `json:"charge_port_cold_weather_mode"`
	// ChargePortColor indicates the status color of the charge port LED.
	ChargePortColor string `json:"charge_port_color"`
	// ChargePortDoorOpen indicates if the charge port door is open.
	ChargePortDoorOpen bool `json:"charge_port_door_open"`
	// ChargePortLatch indicates the state of the charge port latch.
	ChargePortLatch string `json:"charge_port_latch"`
	// ChargeRate is the current charging rate in the vehicle's configured units (e.g., mi/hr).
	ChargeRate float64 `json:"charge_rate"`
	// ChargerActualCurrent is the actual current being drawn by the charger.
	ChargerActualCurrent int `json:"charger_actual_current"`
	// ChargerPilotCurrent is the pilot signal current from the EVSE.
	ChargerPilotCurrent int `json:"charger_pilot_current"`
	// ChargerPower is the current power being delivered by the charger in kW.
	ChargerPower int `json:"charger_power"`
	// ChargerVoltage is the current voltage of the charger.
	ChargerVoltage int `json:"charger_voltage"`
	// ChargingState is the current state of charging (e.g., "Charging", "Disconnected").
	ChargingState string `json:"charging_state"`
	// ConnChargeCable is the type of charge cable connected.
	ConnChargeCable string `json:"conn_charge_cable"`
	// EstBatteryRange is another estimation of the battery range.
	EstBatteryRange float64 `json:"est_battery_range"`
	// FastChargerBrand is the brand of the connected DC fast charger.
	FastChargerBrand string `json:"fast_charger_brand"`
	// FastChargerPresent indicates if a DC fast charger is connected.
	FastChargerPresent bool `json:"fast_charger_present"`
	// FastChargerType is the type of DC fast charger connected (e.g., "Supercharger").
	FastChargerType string `json:"fast_charger_type"`
	// IdealBatteryRange is the ideal estimated remaining driving range.
	IdealBatteryRange float64 `json:"ideal_battery_range"`
	// ManagedChargingActive indicates if a managed charging session is active.
	ManagedChargingActive bool `json:"managed_charging_active"`
	// ManagedChargingUserCanceled indicates if the user has canceled managed charging.
	ManagedChargingUserCanceled bool `json:"managed_charging_user_canceled"`
	// MaxRangeChargeCounter is the number of times the vehicle has been charged to max range.
	MaxRangeChargeCounter int `json:"max_range_charge_counter"`
	// MinutesToFullCharge is the estimated minutes remaining to complete charging to the set limit.
	MinutesToFullCharge int `json:"minutes_to_full_charge"`
	// OffPeakChargingEnabled indicates if off-peak charging is enabled by the user.
	OffPeakChargingEnabled bool `json:"off_peak_charging_enabled"`
	// OffPeakChargingTimes defines the schedule for off-peak charging.
	OffPeakChargingTimes string `json:"off_peak_charging_times"`
	// OffPeakHoursEndTime is the end time for the off-peak charging window, in minutes from midnight.
	OffPeakHoursEndTime int `json:"off_peak_hours_end_time"`
	// PreconditioningEnabled indicates if preconditioning for departure is enabled.
	PreconditioningEnabled bool `json:"preconditioning_enabled"`
	// PreconditioningTimes defines the schedule for preconditioning.
	PreconditioningTimes string `json:"preconditioning_times"`
	// ScheduledChargingMode defines the mode for scheduled charging (e.g., "Off", "Scheduled").
	ScheduledChargingMode string `json:"scheduled_charging_mode"`
	// ScheduledChargingPending indicates if a scheduled charge is pending.
	ScheduledChargingPending bool `json:"scheduled_charging_pending"`
	// ScheduledDepartureTime is the Unix timestamp for the scheduled departure.
	ScheduledDepartureTime int64 `json:"scheduled_departure_time"`
	// ScheduledDepartureTimeMinutes is the departure time in minutes from midnight.
	ScheduledDepartureTimeMinutes int `json:"scheduled_departure_time_minutes"`
	// SuperchargerSessionTripPlanner indicates if the current supercharging session was initiated by the trip planner.
	SuperchargerSessionTripPlanner bool `json:"supercharger_session_trip_planner"`
	// TimeToFullCharge is the estimated time to full charge in hours.
	TimeToFullCharge float64 `json:"time_to_full_charge"`
	// Timestamp is the Unix timestamp (in milliseconds) when this data was recorded.
	Timestamp int64 `json:"timestamp"`
	// TripCharging indicates if trip charging is active.
	TripCharging bool `json:"trip_charging"`
	// UsableBatteryLevel is the usable state of charge of the battery in percent.
	UsableBatteryLevel int `json:"usable_battery_level"`
	// --- Nullable fields ---
	// ChargerPhases is the number of charger phases being used (e.g., 1 or 3). Can be null.
	ChargerPhases *int `json:"charger_phases,omitempty"`
	// ManagedChargingStartTime is the Unix timestamp for when managed charging is scheduled to start. Can be null.
	ManagedChargingStartTime *int64 `json:"managed_charging_start_time,omitempty"`
	// NotEnoughPowerToHeat indicates if there is insufficient power to heat the battery. Can be null.
	NotEnoughPowerToHeat *bool `json:"not_enough_power_to_heat,omitempty"`
	// ScheduledChargingStartTime is the Unix timestamp for when scheduled charging will begin. Can be null.
	ScheduledChargingStartTime *int64 `json:"scheduled_charging_start_time,omitempty"`
	// UserChargeEnableRequest indicates the user's last request to enable/disable charging. Can be null.
	UserChargeEnableRequest *bool `json:"user_charge_enable_request,omitempty"`
}

// ClimateState contains information about the vehicle's climate control system.
type ClimateState struct {
	// AllowCabinOverheatProtection indicates if the user allows Cabin Overheat Protection to be turned on.
	AllowCabinOverheatProtection bool `json:"allow_cabin_overheat_protection"`
	// AutoSeatClimateLeft indicates if the left seat climate is in automatic mode.
	AutoSeatClimateLeft bool `json:"auto_seat_climate_left"`
	// AutoSeatClimateRight indicates if the right seat climate is in automatic mode.
	AutoSeatClimateRight bool `json:"auto_seat_climate_right"`
	// AutoSteeringWheelHeat indicates if the steering wheel heat is in automatic mode.
	AutoSteeringWheelHeat bool `json:"auto_steering_wheel_heat"`
	// BatteryHeater indicates if the battery heater is currently active.
	BatteryHeater bool `json:"battery_heater"`
	// BioweaponMode indicates if Bioweapon Defense Mode is active.
	BioweaponMode bool `json:"bioweapon_mode"`
	// CabinOverheatProtection is the current setting for Cabin Overheat Protection (e.g., "On", "Off").
	CabinOverheatProtection string `json:"cabin_overheat_protection"`
	// CabinOverheatProtectionActivelyCooling indicates if the system is actively cooling to prevent cabin overheat.
	CabinOverheatProtectionActivelyCooling bool `json:"cabin_overheat_protection_actively_cooling"`
	// ClimateKeeperMode is the current state of the Climate Keeper mode (e.g., "off", "on").
	ClimateKeeperMode string `json:"climate_keeper_mode"`
	// CopActivationTemperature is the user-set activation temperature for Cabin Overheat Protection.
	CopActivationTemperature string `json:"cop_activation_temperature"`
	// DefrostMode is the state of the defrost mode. 0 for off.
	DefrostMode int `json:"defrost_mode"`
	// DriverTempSetting is the driver's side temperature setting in Celsius.
	DriverTempSetting float64 `json:"driver_temp_setting"`
	// FanStatus is the current speed of the climate control fan. 0 for off.
	FanStatus int `json:"fan_status"`
	// HvacAutoRequest indicates the automatic request state for the HVAC system.
	HvacAutoRequest string `json:"hvac_auto_request"`
	// InsideTemp is the current temperature inside the cabin in Celsius.
	InsideTemp float64 `json:"inside_temp"`
	// IsAutoConditioningOn indicates if auto conditioning is currently on.
	IsAutoConditioningOn bool `json:"is_auto_conditioning_on"`
	// IsClimateOn indicates if the climate control system is currently active.
	IsClimateOn bool `json:"is_climate_on"`
	// IsFrontDefrosterOn indicates if the front defroster is on.
	IsFrontDefrosterOn bool `json:"is_front_defroster_on"`
	// IsPreconditioning indicates if the vehicle is currently preconditioning.
	IsPreconditioning bool `json:"is_preconditioning"`
	// IsRearDefrosterOn indicates if the rear defroster is on.
	IsRearDefrosterOn bool `json:"is_rear_defroster_on"`
	// LeftTempDirection is an internal value for the left temperature vent direction.
	LeftTempDirection int `json:"left_temp_direction"`
	// MaxAvailTemp is the maximum available temperature setting in Celsius.
	MaxAvailTemp float64 `json:"max_avail_temp"`
	// MinAvailTemp is the minimum available temperature setting in Celsius.
	MinAvailTemp float64 `json:"min_avail_temp"`
	// OutsideTemp is the current temperature outside the vehicle in Celsius.
	OutsideTemp float64 `json:"outside_temp"`
	// PassengerTempSetting is the passenger's side temperature setting in Celsius.
	PassengerTempSetting float64 `json:"passenger_temp_setting"`
	// RemoteHeaterControlEnabled indicates if remote heater control is enabled.
	RemoteHeaterControlEnabled bool `json:"remote_heater_control_enabled"`
	// RightTempDirection is an internal value for the right temperature vent direction.
	RightTempDirection int `json:"right_temp_direction"`
	// SeatHeaterLeft is the heating level for the left seat (0-3).
	SeatHeaterLeft int `json:"seat_heater_left"`
	// SeatHeaterRearCenter is the heating level for the rear center seat (0-3).
	SeatHeaterRearCenter int `json:"seat_heater_rear_center"`
	// SeatHeaterRearLeft is the heating level for the rear left seat (0-3).
	SeatHeaterRearLeft int `json:"seat_heater_rear_left"`
	// SeatHeaterRearRight is the heating level for the rear right seat (0-3).
	SeatHeaterRearRight int `json:"seat_heater_rear_right"`
	// SeatHeaterRight is the heating level for the right seat (0-3).
	SeatHeaterRight int `json:"seat_heater_right"`
	// SideMirrorHeaters indicates if the side mirror heaters are on.
	SideMirrorHeaters bool `json:"side_mirror_heaters"`
	// SteeringWheelHeatLevel is the heating level for the steering wheel.
	SteeringWheelHeatLevel int `json:"steering_wheel_heat_level"`
	// SteeringWheelHeater indicates if the steering wheel heater is on.
	SteeringWheelHeater bool `json:"steering_wheel_heater"`
	// SupportsFanOnlyCabinOverheatProtection indicates if the vehicle supports fan-only mode for Cabin Overheat Protection.
	SupportsFanOnlyCabinOverheatProtection bool `json:"supports_fan_only_cabin_overheat_protection"`
	// Timestamp is the Unix timestamp (in milliseconds) when this data was recorded.
	Timestamp int64 `json:"timestamp"`
	// WiperBladeHeater indicates if the wiper blade heater is on.
	WiperBladeHeater bool `json:"wiper_blade_heater"`
	// --- Nullable fields ---
	// BatteryHeaterNoPower indicates if the battery heater has no power. Can be null.
	BatteryHeaterNoPower *bool `json:"battery_heater_no_power,omitempty"`
}

// DriveState contains information about the vehicle's driving status and location.
type DriveState struct {
	// ActiveRouteLatitude is the latitude of the active navigation route destination.
	ActiveRouteLatitude float64 `json:"active_route_latitude"`
	// ActiveRouteLongitude is the longitude of the active navigation route destination.
	ActiveRouteLongitude float64 `json:"active_route_longitude"`
	// ActiveRouteTrafficMinutesDelay is the current traffic delay in minutes for the active route.
	ActiveRouteTrafficMinutesDelay int `json:"active_route_traffic_minutes_delay"`
	// GpsAsOf is the Unix timestamp of the last GPS fix.
	GpsAsOf int64 `json:"gps_as_of"`
	// Heading is the current direction of the vehicle in degrees (0-359).
	Heading int `json:"heading"`
	// Latitude is the current latitude of the vehicle.
	Latitude float64 `json:"latitude"`
	// Longitude is the current longitude of the vehicle.
	Longitude float64 `json:"longitude"`
	// NativeLatitude is the native latitude provided by the GPS hardware.
	NativeLatitude float64 `json:"native_latitude"`
	// NativeLocationSupported indicates if native location data is supported.
	NativeLocationSupported int `json:"native_location_supported"`
	// NativeLongitude is the native longitude provided by the GPS hardware.
	NativeLongitude float64 `json:"native_longitude"`
	// NativeType is the type of native location data (e.g., "wgs").
	NativeType string `json:"native_type"`
	// Power is the current power draw or regeneration in kW.
	Power int `json:"power"`
	// Timestamp is the Unix timestamp (in milliseconds) when this data was recorded.
	Timestamp int64 `json:"timestamp"`
	// --- Nullable fields ---
	// ShiftState is the current state of the gear shift (e.g., "P", "D", "R"). Can be null.
	ShiftState *string `json:"shift_state,omitempty"`
	// Speed is the current speed of the vehicle. Can be null.
	Speed *float64 `json:"speed,omitempty"`
}

// GUISettings contains user interface settings for the vehicle's display.
type GUISettings struct {
	// Gui24HourTime indicates if the time is displayed in 24-hour format.
	Gui24HourTime bool `json:"gui_24_hour_time"`
	// GuiChargeRateUnits is the user-selected unit for charge rate (e.g., "mi/hr").
	GuiChargeRateUnits string `json:"gui_charge_rate_units"`
	// GuiDistanceUnits is the user-selected unit for distance (e.g., "mi/hr" or "km/hr").
	GuiDistanceUnits string `json:"gui_distance_units"`
	// GuiRangeDisplay is the user-selected display for range (e.g., "Rated").
	GuiRangeDisplay string `json:"gui_range_display"`
	// GuiTemperatureUnits is the user-selected unit for temperature (e.g., "F").
	GuiTemperatureUnits string `json:"gui_temperature_units"`
	// GuiTirepressureUnits is the user-selected unit for tire pressure (e.g., "Psi").
	GuiTirepressureUnits string `json:"gui_tirepressure_units"`
	// ShowRangeUnits indicates if range units are shown.
	ShowRangeUnits bool `json:"show_range_units"`
	// Timestamp is the Unix timestamp (in milliseconds) when this data was recorded.
	Timestamp int64 `json:"timestamp"`
}

// VehicleConfig contains the configuration and specifications of the vehicle.
type VehicleConfig struct {
	// AuxParkLamps is the type of auxiliary park lamps.
	AuxParkLamps string `json:"aux_park_lamps"`
	// BadgeVersion is the version of the vehicle's badge.
	BadgeVersion int `json:"badge_version"`
	// CanAcceptNavigationRequests indicates if the vehicle can accept navigation requests from the API.
	CanAcceptNavigationRequests bool `json:"can_accept_navigation_requests"`
	// CanActuateTrunks indicates if the vehicle's trunks (frunk/trunk) can be opened remotely.
	CanActuateTrunks bool `json:"can_actuate_trunks"`
	// CarSpecialType is the special type of the car (e.g., "base").
	CarSpecialType string `json:"car_special_type"`
	// CarType is the model of the car (e.g., "modely").
	CarType string `json:"car_type"`
	// ChargePortType is the type of charge port (e.g., "US").
	ChargePortType string `json:"charge_port_type"`
	// CopUserSetTempSupported indicates if the user can set a custom temperature for Cabin Overheat Protection.
	CopUserSetTempSupported bool `json:"cop_user_set_temp_supported"`
	// DashcamClipSaveSupported indicates if saving dashcam clips is supported.
	DashcamClipSaveSupported bool `json:"dashcam_clip_save_supported"`
	// DefaultChargeToMax indicates if the default charge limit is set to maximum.
	DefaultChargeToMax bool `json:"default_charge_to_max"`
	// DriverAssist is the type of driver assistance package installed (e.g., "TeslaAP3").
	DriverAssist string `json:"driver_assist"`
	// EceRestrictions indicates if European ECE restrictions are applied.
	EceRestrictions bool `json:"ece_restrictions"`
	// EfficiencyPackage is the efficiency package version (e.g., "MY2021").
	EfficiencyPackage string `json:"efficiency_package"`
	// EuVehicle indicates if the vehicle is a European-spec model.
	EuVehicle bool `json:"eu_vehicle"`
	// ExteriorColor is the exterior color of the vehicle.
	ExteriorColor string `json:"exterior_color"`
	// ExteriorTrim is the exterior trim type (e.g., "Black").
	ExteriorTrim string `json:"exterior_trim"`
	// ExteriorTrimOverride is an override for the exterior trim.
	ExteriorTrimOverride string `json:"exterior_trim_override"`
	// HasAirSuspension indicates if the vehicle has air suspension.
	HasAirSuspension bool `json:"has_air_suspension"`
	// HasLudicrousMode indicates if the vehicle has Ludicrous Mode.
	HasLudicrousMode bool `json:"has_ludicrous_mode"`
	// HasSeatCooling indicates if the vehicle has cooled seats.
	HasSeatCooling bool `json:"has_seat_cooling"`
	// HeadlampType is the type of headlamps installed.
	HeadlampType string `json:"headlamp_type"`
	// InteriorTrimType is the type of interior trim.
	InteriorTrimType string `json:"interior_trim_type"`
	// KeyVersion is the version of the vehicle's key system.
	KeyVersion int `json:"key_version"`
	// MotorizedChargePort indicates if the vehicle has a motorized charge port.
	MotorizedChargePort bool `json:"motorized_charge_port"`
	// PaintColorOverride is an override for the paint color.
	PaintColorOverride string `json:"paint_color_override"`
	// PerformancePackage is the performance package type.
	PerformancePackage string `json:"performance_package"`
	// Plg indicates if power lift gate is present.
	Plg bool `json:"plg"`
	// Pws indicates if a Pedestrian Warning System is present.
	Pws bool `json:"pws"`
	// RearDriveUnit is the type of the rear drive unit.
	RearDriveUnit string `json:"rear_drive_unit"`
	// RearSeatHeaters is the level of rear seat heaters available.
	RearSeatHeaters int `json:"rear_seat_heaters"`
	// RearSeatType is the type of rear seats.
	RearSeatType int `json:"rear_seat_type"`
	// Rhd indicates if the vehicle is right-hand drive.
	Rhd bool `json:"rhd"`
	// RoofColor is the color/type of the roof.
	RoofColor string `json:"roof_color"`
	// SpoilerType is the type of spoiler installed.
	SpoilerType string `json:"spoiler_type"`
	// SupportsQrPairing indicates if the vehicle supports pairing via QR code.
	SupportsQrPairing bool `json:"supports_qr_pairing"`
	// ThirdRowSeats describes the third-row seat configuration.
	ThirdRowSeats string `json:"third_row_seats"`
	// Timestamp is the Unix timestamp (in milliseconds) when this data was recorded.
	Timestamp int64 `json:"timestamp"`
	// TrimBadging is the badging on the trim.
	TrimBadging string `json:"trim_badging"`
	// UseRangeBadging indicates if range badging is used.
	UseRangeBadging bool `json:"use_range_badging"`
	// UtcOffset is the current UTC offset in seconds.
	UtcOffset int `json:"utc_offset"`
	// WebcamSelfieSupported indicates if taking selfies with the webcam is supported.
	WebcamSelfieSupported bool `json:"webcam_selfie_supported"`
	// WebcamSupported indicates if a webcam is present and supported.
	WebcamSupported bool `json:"webcam_supported"`
	// WheelType is the type of wheels installed.
	WheelType string `json:"wheel_type"`
	// --- Nullable fields ---
	// SeatType is the type of front seats. Can be null.
	SeatType *int `json:"seat_type,omitempty"`
	// SunRoofInstalled indicates if a sunroof is installed. Can be null.
	SunRoofInstalled *int `json:"sun_roof_installed,omitempty"`
}

// MediaInfo contains details about the currently playing media.
type MediaInfo struct {
	// A2dpSourceName is the name of the connected A2DP Bluetooth source.
	A2dpSourceName string `json:"a2dp_source_name"`
	// AudioVolume is the current audio volume level.
	AudioVolume float64 `json:"audio_volume"`
	// AudioVolumeIncrement is the step value for volume adjustments.
	AudioVolumeIncrement float64 `json:"audio_volume_increment"`
	// AudioVolumeMax is the maximum possible audio volume.
	AudioVolumeMax float64 `json:"audio_volume_max"`
	// MediaPlaybackStatus is the status of media playback (e.g., "Playing").
	MediaPlaybackStatus string `json:"media_playback_status"`
	// NowPlayingAlbum is the album of the currently playing track.
	NowPlayingAlbum string `json:"now_playing_album"`
	// NowPlayingArtist is the artist of the currently playing track.
	NowPlayingArtist string `json:"now_playing_artist"`
	// NowPlayingDuration is the total duration of the current track in seconds.
	NowPlayingDuration int `json:"now_playing_duration"`
	// NowPlayingElapsed is the elapsed time of the current track in seconds.
	NowPlayingElapsed int `json:"now_playing_elapsed"`
	// NowPlayingSource is the source of the currently playing media.
	NowPlayingSource string `json:"now_playing_source"`
	// NowPlayingStation is the radio station or channel currently playing.
	NowPlayingStation string `json:"now_playing_station"`
	// NowPlayingTitle is the title of the currently playing track or program.
	NowPlayingTitle string `json:"now_playing_title"`
}

// MediaState contains the overall state of the media system.
type MediaState struct {
	// RemoteControlEnabled indicates if remote control of media is enabled.
	RemoteControlEnabled bool `json:"remote_control_enabled"`
}

// SoftwareUpdate contains details about a pending or in-progress software update.
type SoftwareUpdate struct {
	// DownloadPerc is the download percentage of the update.
	DownloadPerc int `json:"download_perc"`
	// ExpectedDurationSec is the expected duration of the update installation in seconds.
	ExpectedDurationSec int `json:"expected_duration_sec"`
	// InstallPerc is the installation percentage of the update.
	InstallPerc int `json:"install_perc"`
	// Status is the current status of the software update.
	Status string `json:"status"`
	// Version is the version number of the software update.
	Version string `json:"version"`
}

// SpeedLimitMode contains settings for the vehicle's speed limit mode.
type SpeedLimitMode struct {
	// Active indicates if speed limit mode is currently active.
	Active bool `json:"active"`
	// CurrentLimitMph is the current speed limit set in miles per hour.
	CurrentLimitMph float64 `json:"current_limit_mph"`
	// MaxLimitMph is the maximum possible speed limit that can be set.
	MaxLimitMph int `json:"max_limit_mph"`
	// MinLimitMph is the minimum possible speed limit that can be set.
	MinLimitMph int `json:"min_limit_mph"`
	// PinCodeSet indicates if a PIN code is set for this mode.
	PinCodeSet bool `json:"pin_code_set"`
}

// VehicleState contains the dynamic state of the vehicle, such as doors locked, odometer, etc.
type VehicleState struct {
	// APIVersion is the API version supported by the vehicle's software.
	APIVersion int `json:"api_version"`
	// AutoparkStateV3 is the current state of the autopark system.
	AutoparkStateV3 string `json:"autopark_state_v3"`
	// AutoparkStyle is the style of autoparking available.
	AutoparkStyle string `json:"autopark_style"`
	// CalendarSupported indicates if the vehicle supports calendar integration.
	CalendarSupported bool `json:"calendar_supported"`
	// CarVersion is the full software version string of the vehicle.
	CarVersion string `json:"car_version"`
	// CenterDisplayState is the state of the center display (e.g., on/off).
	CenterDisplayState int `json:"center_display_state"`
	// DashcamClipSaveAvailable indicates if saving a dashcam clip is currently possible.
	DashcamClipSaveAvailable bool `json:"dashcam_clip_save_available"`
	// DashcamState is the current state of the dashcam (e.g., "Recording", "Unavailable").
	DashcamState string `json:"dashcam_state"`
	// Df is the status of the driver's front door (0=closed).
	Df int `json:"df"`
	// Dr is the status of the driver's rear door (0=closed).
	Dr int `json:"dr"`
	// FdWindow is the status of the front driver's window.
	FdWindow int `json:"fd_window"`
	// FeatureBitmask is a bitmask representing enabled features.
	FeatureBitmask string `json:"feature_bitmask"`
	// FpWindow is the status of the front passenger's window.
	FpWindow int `json:"fp_window"`
	// Ft is the status of the front trunk (frunk) (0=closed).
	Ft int `json:"ft"`
	// HomelinkDeviceCount is the number of programmed Homelink devices.
	HomelinkDeviceCount int `json:"homelink_device_count"`
	// HomelinkNearby indicates if a programmed Homelink device is nearby.
	HomelinkNearby bool `json:"homelink_nearby"`
	// IsUserPresent indicates if a user is detected in the vehicle.
	IsUserPresent bool `json:"is_user_present"`
	// LastAutoparkError is the error message from the last autopark attempt.
	LastAutoparkError string `json:"last_autopark_error"`
	// Locked indicates if the vehicle is locked.
	Locked bool `json:"locked"`
	// MediaInfo contains detailed information about the currently playing media.
	MediaInfo MediaInfo `json:"media_info"`
	// MediaState contains the general state of the media system.
	MediaState MediaState `json:"media_state"`
	// NotificationsSupported indicates if the vehicle supports sending notifications.
	NotificationsSupported bool `json:"notifications_supported"`
	// Odometer is the vehicle's odometer reading in the configured units (miles/km).
	Odometer float64 `json:"odometer"`
	// ParsedCalendarSupported indicates if the vehicle supports parsed calendar events.
	ParsedCalendarSupported bool `json:"parsed_calendar_supported"`
	// Pf is the status of the passenger's front door (0=closed).
	Pf int `json:"pf"`
	// Pr is the status of the passenger's rear door (0=closed).
	Pr int `json:"pr"`
	// RdWindow is the status of the rear driver's side window.
	RdWindow int `json:"rd_window"`
	// RemoteStart indicates if remote start is currently active.
	RemoteStart bool `json:"remote_start"`
	// RemoteStartEnabled indicates if remote start is enabled for the vehicle.
	RemoteStartEnabled bool `json:"remote_start_enabled"`
	// RemoteStartSupported indicates if the vehicle supports remote start.
	RemoteStartSupported bool `json:"remote_start_supported"`
	// RpWindow is the status of the rear passenger's side window.
	RpWindow int `json:"rp_window"`
	// Rt is the status of the rear trunk (0=closed).
	Rt int `json:"rt"`
	// SantaMode indicates if Santa Mode is active.
	SantaMode int `json:"santa_mode"`
	// SentryMode indicates if Sentry Mode is currently active.
	SentryMode bool `json:"sentry_mode"`
	// SentryModeAvailable indicates if Sentry Mode is available on the vehicle.
	SentryModeAvailable bool `json:"sentry_mode_available"`
	// ServiceMode indicates if the vehicle is in service mode.
	ServiceMode bool `json:"service_mode"`
	// ServiceModePlus indicates if the vehicle is in service mode plus.
	ServiceModePlus bool `json:"service_mode_plus"`
	// SmartSummonAvailable indicates if Smart Summon is available.
	SmartSummonAvailable bool `json:"smart_summon_available"`
	// SoftwareUpdate contains information about any available software updates.
	SoftwareUpdate SoftwareUpdate `json:"software_update"`
	// SpeedLimitMode contains information about the speed limit mode settings.
	SpeedLimitMode SpeedLimitMode `json:"speed_limit_mode"`
	// SummonStandbyModeEnabled indicates if Summon standby mode is enabled.
	SummonStandbyModeEnabled bool `json:"summon_standby_mode_enabled"`
	// Timestamp is the Unix timestamp (in milliseconds) when this data was recorded.
	Timestamp int64 `json:"timestamp"`
	// TpmsHardWarningFl indicates a hard tire pressure warning for the front-left tire.
	TpmsHardWarningFl bool `json:"tpms_hard_warning_fl"`
	// TpmsHardWarningFr indicates a hard tire pressure warning for the front-right tire.
	TpmsHardWarningFr bool `json:"tpms_hard_warning_fr"`
	// TpmsHardWarningRl indicates a hard tire pressure warning for the rear-left tire.
	TpmsHardWarningRl bool `json:"tpms_hard_warning_rl"`
	// TpmsHardWarningRr indicates a hard tire pressure warning for the rear-right tire.
	TpmsHardWarningRr bool `json:"tpms_hard_warning_rr"`
	// TpmsLastSeenPressureTimeFl is the Unix timestamp of the last pressure reading for the front-left tire.
	TpmsLastSeenPressureTimeFl int64 `json:"tpms_last_seen_pressure_time_fl"`
	// TpmsLastSeenPressureTimeFr is the Unix timestamp of the last pressure reading for the front-right tire.
	TpmsLastSeenPressureTimeFr int64 `json:"tpms_last_seen_pressure_time_fr"`
	// TpmsLastSeenPressureTimeRl is the Unix timestamp of the last pressure reading for the rear-left tire.
	TpmsLastSeenPressureTimeRl int64 `json:"tpms_last_seen_pressure_time_rl"`
	// TpmsLastSeenPressureTimeRr is the Unix timestamp of the last pressure reading for the rear-right tire.
	TpmsLastSeenPressureTimeRr int64 `json:"tpms_last_seen_pressure_time_rr"`
	// TpmsPressureFl is the pressure of the front-left tire in the configured units (e.g., PSI).
	TpmsPressureFl float64 `json:"tpms_pressure_fl"`
	// TpmsPressureFr is the pressure of the front-right tire.
	TpmsPressureFr float64 `json:"tpms_pressure_fr"`
	// TpmsPressureRl is the pressure of the rear-left tire.
	TpmsPressureRl float64 `json:"tpms_pressure_rl"`
	// TpmsPressureRr is the pressure of the rear-right tire.
	TpmsPressureRr float64 `json:"tpms_pressure_rr"`
	// TpmsRcpFrontValue is the recommended cold pressure for the front tires.
	TpmsRcpFrontValue float64 `json:"tpms_rcp_front_value"`
	// TpmsRcpRearValue is the recommended cold pressure for the rear tires.
	TpmsRcpRearValue float64 `json:"tpms_rcp_rear_value"`
	// TpmsSoftWarningFl indicates a soft tire pressure warning for the front-left tire.
	TpmsSoftWarningFl bool `json:"tpms_soft_warning_fl"`
	// TpmsSoftWarningFr indicates a soft tire pressure warning for the front-right tire.
	TpmsSoftWarningFr bool `json:"tpms_soft_warning_fr"`
	// TpmsSoftWarningRl indicates a soft tire pressure warning for the rear-left tire.
	TpmsSoftWarningRl bool `json:"tpms_soft_warning_rl"`
	// TpmsSoftWarningRr indicates a soft tire pressure warning for the rear-right tire.
	TpmsSoftWarningRr bool `json:"tpms_soft_warning_rr"`
	// ValetMode indicates if valet mode is active.
	ValetMode bool `json:"valet_mode"`
	// ValetPinNeeded indicates if a PIN is needed to disable valet mode.
	ValetPinNeeded bool `json:"valet_pin_needed"`
	// VehicleName is the user-assigned name for the vehicle.
	VehicleName string `json:"vehicle_name"`
	// VehicleSelfTestProgress is the progress of a vehicle self-test.
	VehicleSelfTestProgress int `json:"vehicle_self_test_progress"`
	// VehicleSelfTestRequested indicates if a self-test has been requested.
	VehicleSelfTestRequested bool `json:"vehicle_self_test_requested"`
	// WebcamAvailable indicates if the interior webcam is available.
	WebcamAvailable bool `json:"webcam_available"`
}

// VehicleData represents the entire top-level JSON object for a vehicle's comprehensive data.
type VehicleData struct {
	// ID is the unique identifier for the vehicle record.
	ID int64 `json:"id"`
	// UserID is the ID of the user associated with this vehicle.
	UserID int64 `json:"user_id"`
	// VehicleID is another identifier for the vehicle itself.
	VehicleID int64 `json:"vehicle_id"`
	// Vin is the Vehicle Identification Number.
	Vin string `json:"vin"`
	// AccessType is the type of access the user has to the vehicle (e.g., "OWNER").
	AccessType string `json:"access_type"`
	// GranularAccess contains settings for fine-grained access control.
	GranularAccess GranularAccess `json:"granular_access"`
	// Tokens are the access tokens associated with the vehicle.
	Tokens []string `json:"tokens"`
	// State is the overall state of the vehicle (e.g., "online", "offline").
	State string `json:"state"`
	// InService indicates if the vehicle is currently in a service appointment.
	InService bool `json:"in_service"`
	// IDS is the string representation of the vehicle's ID.
	IDS string `json:"id_s"`
	// CalendarEnabled indicates if calendar integration is enabled.
	CalendarEnabled bool `json:"calendar_enabled"`
	// APIVersion is the general API version for this data payload.
	APIVersion int `json:"api_version"`
	// ChargeState contains all data related to charging.
	ChargeState ChargeState `json:"charge_state"`
	// ClimateState contains all data related to the climate control system.
	ClimateState ClimateState `json:"climate_state"`
	// DriveState contains all data related to the vehicle's driving and location.
	DriveState DriveState `json:"drive_state"`
	// GuiSettings contains all data related to the vehicle's UI settings.
	GuiSettings GUISettings `json:"gui_settings"`
	// VehicleConfig contains all data related to the vehicle's hardware and software configuration.
	VehicleConfig VehicleConfig `json:"vehicle_config"`
	// VehicleState contains all dynamic state information about the vehicle.
	VehicleState VehicleState `json:"vehicle_state"`
	// --- Nullable fields ---
	// Color is the exterior color of the vehicle. Can be null.
	Color *string `json:"color,omitempty"`
	// BackseatToken is a token for the backseat entertainment system. Can be null.
	BackseatToken *string `json:"backseat_token,omitempty"`
	// BackseatTokenUpdatedAt is a timestamp for when the backseat token was last updated. Can be null.
	BackseatTokenUpdatedAt *string `json:"backseat_token_updated_at,omitempty"`
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
func GetVehiceData(accessToken, vin string) (*VehicleData, error) {
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

	// Unmarshal the JSON response into our generic Response struct containing a slice of VehicleData.
	var data Response[VehicleData]
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, errors.Join(err, fmt.Errorf("unmarshal response bytes error"))
	}

	// Check if the API response contains an error.
	if data.Error != "" {
		return nil, fmt.Errorf("tesla response error %s:%s", data.Error, data.ErrorDescription)
	}

	// Return the slice of vehicles from the response.
	return &data.Response, nil
}

// requestAppendAuthorization adds authorization and other necessary headers to an http request.
func requestAppendAuthorization(request *http.Request, accessToken string) {
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Origin", "https://teslatrack.wallora.top")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
}
