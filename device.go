package usbpower

import (
	"errors"
	"time"

	"github.com/sstallion/go-hid"
)

// ErrDeviceNotFound is returned when a compatible device is not found on the
// USB bus of the system.
var ErrDeviceNotFound = errors.New("device not found")

// Device represents a USB device that can read power samples.
type Device interface {
	Read() (Sample, error)
	Close() error
}

// Sample represents a power sample read from the device.
type Sample struct {
	Timestamp time.Time `json:"ts"`
	Voltage   float32   `json:"v"`
	Current   float32   `json:"i"`
}

// OpenDevice opens the first available USB device that matches a known USB
// power meter device. It returns a Device interface that can be used to read
// power samples from the device.
//
// If no compatible device is found, it returns ErrDeviceNotFound.
func OpenDevice() (Device, error) {
	var path string
	err := hid.Enumerate(witrnVendorID, hid.ProductIDAny, func(d *hid.DeviceInfo) error {
		path = d.Path
		return nil
	})
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, ErrDeviceNotFound
	}
	device, err := hid.OpenPath(path)
	if err != nil {
		return nil, err
	}
	return &witrnDevice{
		Device: device,
	}, nil
}
