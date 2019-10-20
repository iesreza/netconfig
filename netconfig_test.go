package netconfig_test

import (
	"fmt"
	"github.com/iesreza/netconfig"
	"testing"
)

func TestGetNetworkConfig(t *testing.T) {
	fmt.Println("start test")
	fmt.Println(netconfig.GetNetworkConfig().String())
}
