// Package eiscp provides basic support for eISCP/ISCP protocol
package eiscp

import (
	"bufio"
	"fmt"
	"net"
)

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
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
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
