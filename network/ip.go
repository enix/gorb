package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// AddressInfo contains informations from ip command about an ip address
type AddressInfo struct {
	Family            string `json:"family"`
	Local             string `json:"local"`
	PrefixLen         int    `json:"prefixlen"`
	Scope             string `json:"scope"`
	Label             string `json:"label"`
	ValidLifeTime     int    `json:"valid_life_time"`
	PreferredLifeTime int    `json:"preferred_life_time"`
}

// InterfaceInfo contains informations from ip command about an interface
type InterfaceInfo struct {
	Ifindex     int           `json:"ifindex"`
	Ifname      string        `json:"ifname"`
	Flags       []string      `json:"flags"`
	Mtu         int           `json:"mtu"`
	Qdisc       string        `json:"qdisc"`
	Operstate   string        `json:"operstate"`
	Group       string        `json:"group"`
	LinkType    string        `json:"link_type"`
	Address     string        `json:"address"`
	Broadcast   string        `json:"broadcast"`
	AddressInfo []AddressInfo `json:"addr_info"`
}

// GetAddress get an ip address from an interface
func GetAddress(ip string, ifaceName string) (*AddressInfo, error) {
	ifaceInfo, err := GetInterface(ifaceName)

	if err != nil {
		return nil, err
	}

	for _, addressInfo := range ifaceInfo.AddressInfo {
		if ip == fmt.Sprintf("%s/%d", addressInfo.Local, addressInfo.PrefixLen) {
			return &addressInfo, nil
		}
	}

	return nil, errors.New("Address not found")
}

// AddAddress attaches an ip address on an interface
func AddAddress(ip string, ifaceName string) (*AddressInfo, error) {
	if _, err := ipAddress("add", ip, "dev", ifaceName); err != nil {
		return nil, err
	}

	return GetAddress(ip, ifaceName)
}

// DeleteAddress detaches an ip address from an interface
func DeleteAddress(ip string, ifaceName string) error {
	_, err := ipAddress("delete", ip, "dev", ifaceName)
	return err
}

// GetInterface get an interface
func GetInterface(ifaceName string) (*InterfaceInfo, error) {
	output, err := ipAddress("show", ifaceName)
	infos := []InterfaceInfo{}

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(output, &infos); err != nil {
		return nil, err
	}

	return &infos[0], nil
}

// AddInterface create an interface
func AddInterface(name string, interfaceType string) (*InterfaceInfo, error) {
	if _, err := ipLink("add", name, "type", interfaceType); err != nil {
		return nil, err
	}

	if _, err := ipLink("set", name, "up"); err != nil {
		return nil, err
	}

	return GetInterface(name)
}

// DeleteInterface delete an interface
func DeleteInterface(name string) error {
	_, err := ipLink("delete", name)
	return err
}

func ipAddress(args ...string) ([]byte, error) {
	return ip("address", args...)
}

func ipLink(args ...string) ([]byte, error) {
	return ip("link", args...)
}

func ip(command string, args ...string) ([]byte, error) {
	args = append([]string{"--json", command}, args...)
	cmd := exec.Command("ip", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("command failed (ip %s): %s (%w)", strings.Join(args, " "), output, err)
	}

	return output, nil
}