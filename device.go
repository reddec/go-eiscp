// Package eiscp provides basic support for eISCP/ISCP protocol
package eiscp

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
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

// Device of Onkyo receiver
type Device struct {
	conn            net.Conn
	channel         *bufio.Reader
	destinationType DeviceType
	version         byte
	Verbose         bool
}

// NewDevice - create and connect to eISCP device (Onkyo)
func NewDevice(host string, deviceType DeviceType, iscpVersion byte) (*Device, error) {
	dev := new(Device)
	dev.destinationType = deviceType
	dev.version = iscpVersion
	err := dev.connect(host)
	if err != nil {
		dev.Close()
		return nil, err
	}
	return dev, nil
}

// NewReceiver - sugar for NewDevice with Receiver as device type and version 1
func NewReceiver(host string) (*Device, error) {
	return NewDevice(host, Receiver, 0x01)
}

// Close connection
func (d *Device) Close() error {
	return d.conn.Close()
}

// ReadMessage - read raw eISCP message
func (d *Device) ReadMessage() (Message, error) {
	msg := Message{}
	err := msg.Parse(d.channel)
	if err != nil {
		return msg, err
	}
	return msg, err
}

// WriteMessage - write raw eISCP message
func (d *Device) WriteMessage(msg Message) error {
	req := msg.BuildEISCP()
	if d.Verbose {
		fmt.Printf("Req: % x\n", req)
	}
	_, err := d.conn.Write(req)
	return err
}

// WriteCommand - write command with arg to remote connection
func (d *Device) WriteCommand(command, arg string) error {
	if d.conn == nil {
		return fmt.Errorf("Not connected")
	}
	msg := Message{}
	msg.Destination = byte(d.destinationType)
	msg.Version = d.version
	msg.ISCP = []byte(command + arg)
	req := msg.BuildEISCP()
	if d.Verbose {
		fmt.Printf("Req: % x\n", req)
	}
	_, err := d.conn.Write(req)
	return err
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

// Utility functions
//  |           |
//  V           V

// Connect device to Onkyo host
func (d *Device) connect(host string) error {
	conn, err := net.Dial("tcp", host)
	d.conn = conn
	if err == nil {
		d.channel = bufio.NewReader(d.conn)
		_, err := d.ReadMessage()
		if err != nil {
			conn.Close()
			return err
		}
	}
	return err
}

func (d *Device) readResponse(command string) error {
	msg, err := d.ReadMessage()
	if err != nil {
		return err
	}
	if string(msg.ISCP[3:]) == "N/A" {
		return fmt.Errorf("Not available")
	}
	return nil
}

func (d *Device) requestSet(command, arg string) error {
	err := d.WriteCommand(command, arg)
	if err != nil {
		return err
	}
	return d.readResponse(command)
}
