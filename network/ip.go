package network

import (
	"encoding/json"
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
	AddressInfo []AddressInfo
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
