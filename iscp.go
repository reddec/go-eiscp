package eiscp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)

// DeviceType - device destination code in ISCP
type DeviceType byte

// Destination code
const (
	Receiver DeviceType = 0x31
)

// Message eISCP
type Message struct {
	Version     byte
	Destination byte
	ISCP        []byte
}

// Parse eISCP message from reader (one-way)
func (msg *Message) Parse(reader *bufio.Reader) error {
	chunk := make([]byte, 4)
	_, err := reader.Read(chunk)
	if err != nil {
		return err
	}
	if string(chunk) != "ISCP" {
		return fmt.Errorf("This is not EISCP message")
	}
	_, err = reader.Read(chunk) //Header size
	if err != nil {
		return err
	}
	if binary.BigEndian.Uint32(chunk) != 16 {
		return fmt.Errorf("Invalid header size")
	}
	reader.Read(chunk) // Data size
	dataSize := binary.BigEndian.Uint32(chunk)
	msg.Version, err = reader.ReadByte()
	if err != nil {
		return err
	}
	reserved := make([]byte, 3)
	_, err = reader.Read(reserved) // Skip reserved
	if err != nil {
		return err
	}
	_, err = reader.ReadByte() // Skip start character
	if err != nil {
		return err
	}
	msg.Destination, err = reader.ReadByte()
	if err != nil {
		return err
	}
	msg.ISCP = make([]byte, dataSize-2) // Trim leading control characters
	_, err = reader.Read(msg.ISCP)
	msg.ISCP = msg.ISCP[:len(msg.ISCP)-3] // Trim trailing EOF and EOL characters
	return err
}

// BuildISCP - Build ISCP message
func (msg *Message) BuildISCP() []byte {
	buffer := bytes.Buffer{}
	buffer.WriteRune('!')             // Start character
	buffer.WriteByte(msg.Destination) // Receiver
	buffer.Write(msg.ISCP)
	buffer.Write([]byte{0x0D})
	return buffer.Bytes()
}

// BuildEISCP - Build ISCP message into ethernet frame
func (msg *Message) BuildEISCP() []byte {
	iscp := msg.BuildISCP()
	sizebuf := make([]byte, 4)
	buffer := bytes.Buffer{}
	buffer.WriteString("ISCP")
	buffer.Write([]byte{0, 0, 0, 0x10}) // Header size

	binary.BigEndian.PutUint32(sizebuf, uint32(len(iscp)))
	buffer.Write(sizebuf)         // Data size
	buffer.WriteByte(msg.Version) // Version
	buffer.Write([]byte{0, 0, 0}) // Reserved

	buffer.Write(iscp) //Data
	return buffer.Bytes()
}
