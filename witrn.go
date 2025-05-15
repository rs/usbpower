package usbpower

import (
	"encoding/binary"
	"errors"
	"math"
	"time"

	"github.com/sstallion/go-hid"
)

// Offsets for packet fields
const (
	witrnVendorID = 0x716

	witrnPacketLength = 64

	witrnOffsetHeader              = 0
	witrnOffsetTimestampSec        = 2
	witrnOffsetTimestampMSec       = 3
	witrnOffsetUnknown1            = 4
	witrnOffsetTimestampMSecMod100 = 6
	witrnOffsetUnknown2            = 7
	witrnOffsetAh                  = 14
	witrnOffsetWh                  = 18
	witrnOffsetRecordTime          = 22
	witrnOffsetTime                = 26
	witrnOffsetDPlus               = 30
	witrnOffsetDMinus              = 34
	witrnOffsetUnknown3            = 38
	witrnOffsetUnknown4            = 42
	witrnOffsetVoltage             = 46
	witrnOffsetCurrent             = 50
	witrnOffsetUnknown5            = 54
	witrnOffsetUnknown6            = 58
	witrnDevicePayloadChecksum     = 62
	witrnDeviceHeaderChecksum      = 63
)

type witrnDevice struct {
	Device *hid.Device

	timestampEpoch       time.Time
	lastTimeS            int
	timestampWrapArounds int
	buf                  []byte
}

func (w *witrnDevice) Read() (Sample, error) {
	if w.buf == nil {
		w.buf = make([]byte, witrnPacketLength)
	}

	n, err := w.Device.ReadWithTimeout(w.buf, 500*time.Millisecond)
	if err != nil {
		return Sample{}, err
	}
	packet := w.buf[0:n]
	if !w.isValidPacket(packet) {
		return Sample{}, errors.New("invalid packet")
	}
	return Sample{
		Timestamp: w.computeTimestamp(packet),
		Voltage:   math.Float32frombits(binary.LittleEndian.Uint32(packet[witrnOffsetVoltage:])),
		Current:   math.Float32frombits(binary.LittleEndian.Uint32(packet[witrnOffsetCurrent:])),
	}, nil
}

func (w *witrnDevice) isValidPacket(packet []byte) bool {
	if len(packet) != witrnPacketLength {
		return false
	}
	if packet[witrnOffsetHeader] != 0xff || packet[witrnOffsetHeader+1] != 0x55 {
		return false
	}
	var packetChecksum byte
	for i := 0; i < 8; i++ {
		packetChecksum += packet[i]
	}
	var payloadChecksum byte
	for i := 8; i < witrnDevicePayloadChecksum; i++ {
		payloadChecksum += packet[i]
	}
	packetChecksum += payloadChecksum
	return packetChecksum == packet[witrnDeviceHeaderChecksum] &&
		payloadChecksum == packet[witrnDevicePayloadChecksum]

}

func (w *witrnDevice) computeTimestamp(data []byte) time.Time {
	timeS := int(data[witrnOffsetTimestampMSec])
	timeMs := int(data[witrnOffsetTimestampMSec])
	timeMsMod100 := int(data[witrnOffsetTimestampMSecMod100])

	if timeS < 10 && w.lastTimeS > 250 {
		w.timestampWrapArounds++
	}
	w.lastTimeS = timeS

	var baseMs int
	found := false
	for j := 0; j <= 3; j++ {
		ms := (256*j + timeMs) % 1000
		if ms%100 == timeMsMod100 {
			baseMs = ms
			found = true
			break
		}
	}
	if !found {
		baseMs = 0
	}

	msSinceEpoch := (timeS+w.timestampWrapArounds*256)*1000 + baseMs
	durationSinceEpoch := time.Duration(msSinceEpoch) * time.Millisecond

	if w.timestampEpoch.IsZero() {
		w.timestampEpoch = time.Now().Add(-durationSinceEpoch)
	}

	return w.timestampEpoch.Add(durationSinceEpoch)
}

func (w *witrnDevice) Close() error {
	return w.Device.Close()
}
