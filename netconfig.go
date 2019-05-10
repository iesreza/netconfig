package netconfig

import (
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

type Network struct {
	LocalIP net.IP
	DNS []string
	SubnetMask net.IP
	DefaultGateway net.IP
	InterfaceName string
	HardwareAddress net.HardwareAddr
	Suffix string
	Interface net.Interface
}

func GetNetworkConfig() *Network  {
	network := Network{}
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	network.LocalIP = conn.LocalAddr().(*net.UDPAddr).IP

	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {

		if addrs, err := interf.Addrs(); err == nil {
			for _, addr := range addrs {
				if strings.Contains(addr.String(), network.LocalIP.String()) {
					network.InterfaceName = interf.Name
					network.HardwareAddress = interf.HardwareAddr
					network.Interface = interf
				}
			}
		}
	}



	if runtime.GOOS == "windows" {
		network.getWindows()
	}else{
		network.getLinux()
	}

	return  &network
}

func (network *Network) getLinux(){

}

func (network *Network) String() string{

	res := "InterfaceName:"+network.InterfaceName+"\r\n"
	res += "HardwareAddress:"+network.HardwareAddress.String()+"\r\n"
	res += "LocalIP:"+network.LocalIP.String()+"\r\n"
	res += "DNS:"+strings.Join(network.DNS,",")+"\r\n"
	res += "SubnetMask:"+network.SubnetMask.String()+"\r\n"
	res += "DefaultGateway:"+network.DefaultGateway.String()+"\r\n"
	res += "Suffix:"+network.Suffix+"\r\n"

	return res
}

func (network *Network) getWindows() {
	out, err := exec.Command("ipconfig","/all").Output()
	if err != nil {
		log.Fatal(err)
	}
	items := strings.Split(string(out),"Ethernet adapter ")
	for _,item := range items{
		if strings.HasPrefix(item,network.InterfaceName){
			lines := strings.Split(item,"\r\n")
			network.DefaultGateway = net.ParseIP(extractDotted(lines,"Default Gateway")[0])
			network.DNS = extractDotted(lines,"DNS Servers")
			network.Suffix = extractDotted(lines,"Connection-specific DNS Suffix")[0]
			network.SubnetMask = net.ParseIP(extractDotted(lines,"Subnet Mask")[0])

		}
	}
}

func extractDotted(lines []string,key string) []string  {
	result := ""
	found := false

	for _,line := range lines{
		if !found {
			if strings.HasPrefix(line,"   "+key){
				result = line[39:]+""
				found = true
			}
		}else{

			if len(line) > 39 && strings.TrimSpace(line[0:39]) == ""{
				result += ","+strings.TrimSpace(line[39:])
			}else{
				break
			}
		}

	}

	return  strings.Split( strings.Trim(result,",") , "," )
}