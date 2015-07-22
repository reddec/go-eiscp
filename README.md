# go-eiscp
eISCP protocol for Onkyo
Provides minimal eISCP protocol

# Example

### Create connection

```go
dev, err := eiscp.NewReceiver(*host)
if err != nil {
  panic(err)
}
defer dev.Close()
// Do something else.....
```

### Set volume level to 50%

```go
err := dev.SetVolume(uint8(50))
// Process error....
```

### Get power state

```go
enabled, err := dev.GetPower()
// Process error and state...
```

### Write Onkyo command

```go
err := dev.WriteCommand("PWR", "01")
// Process error...
// Usually requires wait response
```

### Write raw eISCP message

```go
msg := Message{}
msg.Destination = 0x31  // Destination object
msg.Version = 0x01      // ISCP version
msg.ISCP = []byte("SOME-Command")
err := dev.WriteMessage(msg)
// Process error....
```

### Read raw eISCP message

```go
msg, err := dev.ReadMessage()
// Process error and message...
```
