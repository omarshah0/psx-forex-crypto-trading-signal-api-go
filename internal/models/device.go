package models

// DeviceType represents the type of device making the request
type DeviceType string

const (
	DeviceTypeWeb    DeviceType = "web"
	DeviceTypeMobile DeviceType = "mobile"
)

// IsValid checks if the device type is valid
func (d DeviceType) IsValid() bool {
	return d == DeviceTypeWeb || d == DeviceTypeMobile
}

// String returns the string representation of the device type
func (d DeviceType) String() string {
	return string(d)
}

