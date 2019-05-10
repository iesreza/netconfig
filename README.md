# netconfig
Golang Network Configuration Reader
Read Cross Platform Network Configuration
Currently Can read
- InterfaceName
- HardwareAddress
- LocalIP
- DNS
- SubnetMask
- DefaultGateway
- Suffix

## Example
```
package main

import (
	"fmt"
	"netconfig"
	"strings"
)

func main() {
	
	network := netconfig.GetNetworkConfig()

	res := "InterfaceName:"+network.InterfaceName+"\r\n"
	res += "HardwareAddress:"+network.HardwareAddress.String()+"\r\n"
	res += "LocalIP:"+network.LocalIP.String()+"\r\n"
	res += "DNS:"+strings.Join(network.DNS,",")+"\r\n"
	res += "SubnetMask:"+network.SubnetMask.String()+"\r\n"
	res += "DefaultGateway:"+network.DefaultGateway.String()+"\r\n"
	res += "Suffix:"+network.Suffix+"\r\n"

	fmt.Println(res)

}
```
