package main

import (
	"flag"
	"fmt"
	"github.com/reddec/go-eiscp"
	"strconv"
)

func main() {
	host := flag.String("host", "", "Onkyo host")
	param := flag.String("param", "", "Param name")
	value := flag.String("value", "", "Param value. Empty means only get")
	listSources := flag.Bool("list-source", false, "List source")
	flag.Parse()

	if *listSources {
		for k := range eiscp.SourceByName {
			fmt.Println(k)
		}
		return
	}
	dev, err := eiscp.NewReceiver(*host)
	if err != nil {
		panic(err)
	}
	defer dev.Close()
	if *value == "" {
		switch *param {
		case "power":
			fmt.Println(dev.GetPower())
		case "volume":
			fmt.Println(dev.GetVolume())
		case "source":
			src, err := dev.GetSource()
			if err == nil {
				fmt.Println(eiscp.SourceToName[src])
			} else {
				fmt.Println(err)
			}
		default:
			panic("Unknow param")
		}
	} else {
		switch *param {
		case "power":
			v, err := strconv.ParseBool(*value)
			if err != nil {
				panic(err)
			}
			fmt.Println(dev.SetPower(v))
		case "volume":
			v, err := strconv.ParseInt(*value, 10, 8)
			if err != nil {
				panic(err)
			}
			fmt.Println(dev.SetVolume(uint8(v)))
		case "source":
			src, ok := eiscp.SourceByName[*value]
			if !ok {
				panic("Unknown source")
			}
			fmt.Println(dev.SetSource(src))
		default:
			panic("Unknow param")
		}
	}
}
