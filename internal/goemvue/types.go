package goemvue

import "time"

type CustomerInfo struct {
	CustomerGID int      `json:"customerGid"`
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	CreatedAt   string   `json"createdAt"`
}

type CustomerInfoWithDevices struct {
	CustomerInfo
	Devices     []Device `json:"devices"`
}

type Device struct {
	DeviceGid            int                `json:"deviceGid"`
	ManufacturerDeviceID string             `json:"manufacturerDeviceId"`
	Model                string             `json:"model"`
	Firmware             string             `json:"firmware"`
	ParentDeviceGid      *int               `json:"parentDeviceGid"`
	ParentChannelNum     interface{}        `json:"parentChannelNum"`
	LocationProperties   LocationProperties `json:"locationInformation"`
	Outlet               interface{}        `json:"outlet"`
	EvCharger            interface{}        `json:"evCharger"`

	Devices  []Device  `json:"devices"`
	Channels []Channel `json:"channels"`
}

type LocationProperties struct {
	DeviceGid             int                 `json:"deviceGid"`
	Name                  int                 `json:"deviceName"`
	ZipCode               string              `json:"zipCode"`
	TimeZone              string              `json:"timeZone"`
	BillingCycleStartDay  int                 `json:"billingCycleStartDay"`
	UsageCentPerKwHour    float64             `json:"usageCentPerKwHour"`
	PeakDemandDollarPerKw float64             `json:"peakDemandDollarPerKw"`
	LocationInformation   LocationInformation `json:"locationInformation"`
	LatitudeLongitude     LatitudeLongitude   `json:"latitudeLongitude"`
	UtilityRateGid        *int                `json:"utilityRateGid"`
}
type LocationInformation struct {
	HeatSource   string   `json:"heatSource"`
	LocationSqFt *float64 `json:"locationSqFt,string"`
	LocationType string   `json:"locationType"`
	NumPeople    *int     `json:"numPeople,string"`
	HotTub       *bool    `json:"hotTub,string"`
}

type LatitudeLongitude struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Channel struct {
	ChannelTypeGid    int     `json:"channelTypeGid"`
	Name              string  `json:"name"`
	ChannelNum        string  `json:"channelNum"`
	DeviceGid         int     `json:"deviceGid"`
	ChannelMultiplier float64 `json:"channelMultiplier"`
}


/*
 * History response allows for null values
 */
type History struct {
	Start *time.Time `json:"firstUsageInstant"`
	Data  []*float64  `json:"usageList"`
}

type ChannelList []Channel

