package eiscp_test

import (
	"github.com/reddec/go-eiscp"
)

func ExampleNewReceiver() {
	host := "192.168.1.2:60128"         // 60128 is default port for Onkyo receiver
	dev, err := eiscp.NewReceiver(host) // Create device and connect
	if err != nil {
		panic(err)
	}
	defer dev.Close() // Do not forget close connection
}

func ExampleDevice_WriteCommand() {
	err = dev.WriteCommand("PWR", "01") // Send command and argument to receiver
}

func ExampleDevice_SetSource() {
	source := eiscp.SourceByName["pc"] // Find source by name 'PC'
	err = dev.SetSource(source)        // Setup source on receiver
}

func ExampleDevice_GetSource() {
	code, err := dev.GetSource()
	if err != nil {
		panic(err)
	}
	source := SourceToName[code] // Get readable name of source channel
}
