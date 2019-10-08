package netconfig

import (
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

// Network is the interface which store network configuration data
type Network struct {
	LocalIP         net.IP
	DNS             []string
	SubnetMask      net.IP
	DefaultGateway  net.IP
	InterfaceName   string
	HardwareAddress net.HardwareAddr
	Suffix          string
	Interface       *net.Interface
}

var instance *Network

func Refresh() *Network {
	if instance != nil {
		instance = nil
	}
	return GetNetworkConfig()
}

// GetNetworkConfig create instance of network configuration.
func GetNetworkConfig() *Network {
	if instance != nil {
		return instance
	}
	network := Network{}

	if runtime.GOOS == "windows" {
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
						//network.Interface = &interf
					}
				}
			}
		}
		network.getWindows()
	} else {
		network.getLinux()
	}
	instance = &network
	return &network
}

// getLinux read network data for linux
func (network *Network) getLinux() {

	out, err := exec.Command("/bin/ip", "route", "get", "8.8.8.8").Output()

	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(string(out), " ")
	network.DefaultGateway = net.ParseIP(parts[2])
	network.InterfaceName = parts[4]
	network.LocalIP = net.ParseIP(parts[6])

	interf, err := net.InterfaceByName(network.InterfaceName)
	if err == nil {
		network.HardwareAddress = interf.HardwareAddr
		network.Interface = interf
		log.Println(interf)
	}

	out, err = exec.Command("/sbin/ifconfig", network.InterfaceName).Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")

		if len(lines) > 1 {
			network.SubnetMask = net.ParseIP(strings.Split(strings.TrimSpace(lines[1]), " ")[4])
		}
	}

	out, err = exec.Command("grep", "domain-name", "/var/lib/dhcp/dhclient."+network.InterfaceName+".leases").Output()

	if err == nil {
		dnslist := ""
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			if strings.Contains(line, "domain-name-servers") {
				if len(line) > 26 {
					line = strings.TrimRight(strings.TrimSpace(line)[26:], ";")
					list := strings.Split(line, ",")
					for _, dnsitem := range list {
						if !strings.Contains(dnslist, dnsitem) {
							dnslist += dnsitem + ","
						}
					}
				}

			} else {
				if len(line) > 18 {
					network.Suffix = strings.TrimSpace(strings.TrimRight(strings.TrimSpace(line)[18:], ";"))
				}
			}
			dnslist = strings.TrimRight(dnslist, ",")
		}

		network.DNS = strings.Split(dnslist, ",")
	}

}

// String return network information as string
func (network *Network) String() string {

	res := "InterfaceName:" + network.InterfaceName + "\r\n"
	res += "HardwareAddress:" + network.HardwareAddress.String() + "\r\n"
	res += "LocalIP:" + network.LocalIP.String() + "\r\n"
	res += "DNS:" + strings.Join(network.DNS, ",") + "\r\n"
	res += "SubnetMask:" + network.SubnetMask.String() + "\r\n"
	res += "DefaultGateway:" + network.DefaultGateway.String() + "\r\n"
	res += "Suffix:" + network.Suffix + "\r\n"

	return res
}

// getWindows read network data in windows
func (network *Network) getWindows() {
	out, err := exec.Command("ipconfig", "/all").Output()
	if err != nil {
		log.Fatal(err)
	}
	items := strings.Split(string(out), "Ethernet adapter ")
	for _, item := range items {
		if strings.HasPrefix(item, network.InterfaceName) {
			lines := strings.Split(item, "\r\n")
			network.DefaultGateway = net.ParseIP(extractDotted(lines, "Default Gateway")[0])
			network.DNS = extractDotted(lines, "DNS Servers")
			network.Suffix = extractDotted(lines, "Connection-specific DNS Suffix")[0]
			network.SubnetMask = net.ParseIP(extractDotted(lines, "Subnet Mask")[0])

		}
	}
}

// extractDotted extract data of ipconfig
func extractDotted(lines []string, key string) []string {
	result := ""
	found := false

	for _, line := range lines {
		if !found {
			if strings.HasPrefix(line, "   "+key) {
				result = line[39:] + ""
				found = true
			}
		} else {

			if len(line) > 39 && strings.TrimSpace(line[0:39]) == "" {
				result += "," + strings.TrimSpace(line[39:])
			} else {
				break
			}
		}

	}

	return strings.Split(strings.Trim(result, ","), ",")
}
