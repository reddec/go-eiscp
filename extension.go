package eiscp

import (
	"encoding/hex"
	"strconv"
	"strings"
)

// Source name of input channel
type Source string

// Sample sources
const (
	SrcVCR           Source = "00"
	SrcCBL                  = "01"
	SrcGame                 = "02"
	SrcAux1                 = "03"
	SrcAux2                 = "04"
	SrcPC                   = "05"
	SrcDVD                  = "10"
	SrcPhono                = "22"
	SrcCD                   = "23"
	SrcFM                   = "24"
	SrcAM                   = "25"
	SrcTuner                = "26"
	SrcDLNA2                = "27"
	SrcInternetRadio        = "28"
	SrcUsbFront             = "29"
	SrcUsbRear              = "2A"
	SrcNetwork              = "2B"
)

// SourceByName - map channel name to source enum const
var SourceByName = map[string]Source{
	"vcr":            SrcVCR,
	"cbl":            SrcCBL,
	"game":           SrcGame,
	"aux1":           SrcAux1,
	"aux2":           SrcAux2,
	"pc":             SrcPC,
	"dvd":            SrcDVD,
	"phono":          SrcPhono,
	"cd":             SrcCD,
	"fm":             SrcFM,
	"am":             SrcAM,
	"tuner":          SrcTuner,
	"dlna2":          SrcDLNA2,
	"internet-radio": SrcInternetRadio,
	"usb-front":      SrcUsbFront,
	"usb-rear":       SrcUsbRear,
	"network":        SrcNetwork,
}

// SourceToName - map source enum to channel name
var SourceToName = map[Source]string{
	SrcVCR:           "vcr",
	SrcCBL:           "cbl",
	SrcGame:          "game",
	SrcAux1:          "aux1",
	SrcAux2:          "aux2",
	SrcPC:            "pc",
	SrcDVD:           "dvd",
	SrcPhono:         "phono",
	SrcCD:            "cd",
	SrcFM:            "fm",
	SrcAM:            "am",
	SrcTuner:         "tuner",
	SrcDLNA2:         "dlna2",
	SrcInternetRadio: "internet-radio",
	SrcUsbFront:      "usb-front",
	SrcUsbRear:       "usb-rear",
	SrcNetwork:       "network",
}

// SetSource - Set Onkyo source channel
func (d *Device) SetSource(source Source) error {
	return d.requestSet("SLI", string(source))
}

// GetSource - Get Onkyo source channel. Use SourceToName to get readable name
func (d *Device) GetSource() (Source, error) {
	err := d.WriteCommand("SLI", "QSTN")
	if err != nil {
		return Source("err"), err
	}
	msg, err := d.ReadMessage()
	return Source(msg.ISCP[3:]), err
}

// SetPower - turn on/off Onkyo device
func (d *Device) SetPower(on bool) error {
	if on {
		return d.requestSet("PWR", "01")
	}
	return d.requestSet("PWR", "00")
}

// GetPower - get Onkyo power state
func (d *Device) GetPower() (bool, error) {
	err := d.WriteCommand("PWR", "QSTN")
	if err != nil {
		return false, err
	}
	msg, err := d.ReadMessage()
	return string(msg.ISCP[3:]) == "01", err
}

// SetVolume - set master volume in Onkyo receiver
func (d *Device) SetVolume(level uint8) error {
	return d.requestSet("MVL", strings.ToUpper(hex.EncodeToString([]byte{level})))
}

// GetVolume - get master volume in Onkyo receiver
func (d *Device) GetVolume() (uint8, error) {
	err := d.WriteCommand("MVL", "QSTN")
	if err != nil {
		return 0, err
	}
	msg, err := d.ReadMessage()
	vol, err := strconv.ParseUint(string(msg.ISCP[3:]), 16, 8)
	return uint8(vol), err
}
