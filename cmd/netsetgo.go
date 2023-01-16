package cmd

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/P1-Team/netsetgo"
	"github.com/P1-Team/netsetgo/configurer"
	"github.com/P1-Team/netsetgo/device"
	"github.com/P1-Team/netsetgo/netns"
)

func main() {
	var bridgeName, bridgeAddress, containerAddress, vethNamePrefix string
	var pid int

	flag.StringVar(&bridgeName, "bridgeName", "brg0", "Name to assign to bridge device")
	flag.StringVar(&bridgeAddress, "bridgeAddress", "10.10.10.1/24", "Address to assign to bridge device (CIDR notation)")
	flag.StringVar(&vethNamePrefix, "vethNamePrefix", "veth", "Name prefix for veth devices")
	flag.StringVar(&containerAddress, "containerAddress", "10.10.10.2/24", "Address to assign to the container (CIDR notation)")
	flag.IntVar(&pid, "pid", 0, "pid of a process in the container's network namespace")
	flag.Parse()

	if pid == 0 {
		fmt.Println("ERROR - netsetgo needs a pid")
		os.Exit(1)
	}

	bridgeCreator := device.NewBridge()
	vethCreator := device.NewVeth()
	netnsExecer := &netns.Execer{}

	hostConfigurer := configurer.NewHostConfigurer(bridgeCreator, vethCreator)
	containerConfigurer := configurer.NewContainerConfigurer(netnsExecer)
	netset := netsetgo.New(hostConfigurer, containerConfigurer)

	bridgeIP, bridgeSubnet, err := net.ParseCIDR(bridgeAddress)
	Check(err)

	containerIP, _, err := net.ParseCIDR(containerAddress)
	Check(err)

	netConfig := netsetgo.NetworkConfig{
		BridgeName:     bridgeName,
		BridgeIP:       bridgeIP,
		ContainerIP:    containerIP,
		Subnet:         bridgeSubnet,
		VethNamePrefix: vethNamePrefix,
	}

	Check(netset.ConfigureHost(netConfig, pid))
	Check(netset.ConfigureContainer(netConfig, pid))
}

func Check(err error) error {
	if err != nil {
		return err
	}
	return nil
}
