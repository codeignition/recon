// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// +build linux

// Package network provides different network related data.
package network

import (
	"net"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hariharan-uno/recon/internal/fileutil"
)

// Network addresses
var (
	IPV4Addr string
	IPV6Addr string
	MacAddr  string
)

// Data represents the network data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	d["interfaces"] = make(map[string]interface{})
	ifaces := d["interfaces"].(map[string]interface{})
	out, err := exec.Command("ip", "addr").Output()
	if err != nil {
		return d, err
	}
	lines := strings.Split(string(out), "\n")
	var (
		name  string
		iface map[string]interface{}
		addrs map[string]interface{}
	)
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.ContainsAny(line, "<>") {
			// 1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default
			a := strings.Fields(line)
			name = strings.Trim(a[1], " :")
			ifaces[name] = make(map[string]interface{})
			iface = ifaces[name].(map[string]interface{})
			flags := strings.Split(strings.Trim(a[2], "<> "), ",")
			iface["flags"] = flags
			for i := range a {
				if a[i] == "mtu" {
					iface["mtu"] = a[i+1]
				}

				if a[i] == "state" {
					iface["state"] = strings.ToLower(a[i+1])
				}
			}
			iface["addresses"] = make(map[string]interface{})
			addrs = iface["addresses"].(map[string]interface{})

		} else {
			line = strings.TrimSpace(line)
			// link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
			// link/ether 4c:eb:42:0f:8e:99 brd ff:ff:ff:ff:ff:ff
			if strings.HasPrefix(line, "link/") {
				a := strings.Fields(line)
				iface["encapsulation"] = strings.TrimPrefix(a[0], "link/")
				if len(a) >= 4 {
					if a[1] != "00:00:00:00:00:00" {
						addrs[a[1]] = map[string]string{
							"family": "lladdr",
						}
					}
				}
			}

			if strings.HasPrefix(line, "inet") {
				a := strings.Fields(line)
				// inet 192.168.1.119/24 brd 192.168.1.255 scope global wlan0
				// inet 127.0.0.1/8 scope host lo
				if len(a) >= 5 {
					ip, ipnet, err := net.ParseCIDR(a[1])
					if err != nil {
						return d, err
					}
					ones, _ := ipnet.Mask.Size() // ignore the number of bits
					var tempScope string
					for i := range a {
						if a[i] == "scope" {
							tempScope = scope(a[i+1])
						}
					}
					addrs[ip.String()] = make(map[string]string)
					t := addrs[ip.String()].(map[string]string)
					t["family"] = "inet"
					t["prefixlen"] = strconv.Itoa(ones)
					t["scope"] = tempScope

					// by converting into IP type, we get the string in the form a.b.c.d
					t["netmask"] = net.IP(ipnet.Mask).String()

					for i := range a {
						if a[i] == "brd" {
							t["broadcast"] = a[i+1]
						}

						if a[i] == "peer" {
							t["peer"] = a[i+1]
						}
					}
				}
			}

			if strings.HasPrefix(line, "inet6") {
				a := strings.Fields(line)
				// inet6 ::1/128 scope host
				if len(a) >= 4 {
					ip, ipnet, err := net.ParseCIDR(a[1])
					if err != nil {
						return d, err
					}
					ones, _ := ipnet.Mask.Size() // ignore the number of bits
					addrs[ip.String()] = map[string]string{
						"family":    "inet6",
						"prefixlen": strconv.Itoa(ones),
						"scope":     scope(a[3]),
					}
				}
			}
		}
	}

	if err := defaultGateway(d); err != nil {
		return d, err
	}

	if err := routes(d); err != nil {
		return d, err
	}

	if err := neighbours(d); err != nil {
		return d, err
	}

	// arp data is added in the neighbours function itself.
	// If that fails, this one adds the arp data.
	if err := arp(d); err != nil {
		return d, err
	}

	// last function to call as the data needs to be filled by
	// the above function calls.
	populateAddrs(d)
	return d, nil
}

func scope(s string) string {
	if s == "host" {
		return "Node"
	}
	return strings.Title(s)
}

// defaultGateway adds the default gateway and default interface
// data to the given map.
func defaultGateway(d Data) error {
	out, err := exec.Command("route", "-n").Output()
	if err != nil {
		return err
	}
	s := strings.Split(string(out), "\n")
	// s[0] is the title, s[1] is the column headings. Also, we only
	// consider s[2] for the default interface and gateway.
	a := strings.Fields(s[2])
	d["default_gateway"] = a[1]
	d["default_interface"] = a[7]

	out, err = exec.Command("route", "-6", "-n").Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(out), "\n")
	// lines[0] is the title, lines[1] is the column headings.
	for _, line := range lines[2:] {
		a := strings.Fields(line)
		if len(a) >= 7 {
			if a[1] != "::" {
				d["default_inet6_gateway"] = a[1]
				d["default_inet6_interface"] = a[6]
				break
			}
		}
	}
	return nil
}

// populateAddrs populates the address variables
// exported by the package.
func populateAddrs(d Data) {
	iface := d["default_interface"].(string)
	ifaces := d["interfaces"].(map[string]interface{})
	ifaceMap := ifaces[iface].(map[string]interface{})
	addresses := ifaceMap["addresses"].(map[string]interface{})
	for k, val := range addresses {
		v := val.(map[string]string)
		if v["family"] == "inet" {
			IPV4Addr = k
		}
		if v["family"] == "inet6" {
			IPV6Addr = k
		}
		if v["family"] == "lladdr" {
			MacAddr = k
		}
	}
}

// arp adds the arp data to the given map.
func arp(d Data) error {
	out, err := exec.Command("arp", "-an").Output()
	if err != nil {
		return err
	}
	s := strings.Split(string(out), "\n")
	for i := range s {
		line := s[i]
		// ? (192.164.1.1) at 48:f8:c3:46:03:44 [ether] on wlan
		a := strings.Fields(line)
		if len(a) >= 4 {
			d["arp"] = map[string]string{
				strings.Trim(a[1], "()"): a[3],
			}
		}
	}
	return nil
}

// ipv6Enabled returns true if IPv6 is enabled
// on the machine.
func ipv6Enabled() bool {
	return fileutil.Exists("/proc/net/if_inet6")
}

// families to get default routes from.
var families = []map[string]string{
	{
		"name":               "inet",
		"defaultRoute":       "0.0.0.0/0",
		"defaultPrefix":      "default",
		"neighbourAttribute": "arp",
	},
}

func init() {
	if ipv6Enabled() {
		families = append(families, map[string]string{
			"name":               "inet6",
			"defaultRoute":       "::/0",
			"defaultPrefix":      "default_inet6",
			"neighbourAttribute": "neighbour_inet6",
		})
	}
}

// neighbours adds the neighbours data of families to the given map.
func neighbours(d Data) error {
	for _, family := range families {
		out, err := exec.Command("ip", "-f", family["name"], "neigh", "show").Output()
		if err != nil {
			return err
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			a := strings.Fields(line)
			// 192.167.1.1 dev wlan0 lladdr 48:f8:b3:36:03:44 REACHABLE
			// fe80::4af7:b3ff:fe36:344 dev wlan0 lladdr 48:f8:b3:36:06:44 router STALE
			if len(a) >= 5 {
				d[family["neighbourAttribute"]] = map[string]string{
					a[0]: a[4],
				}
			}
		}
	}
	return nil
}

func routes(d Data) error {
	var routes []map[string]string
	for _, family := range families {
		out, err := exec.Command("ip", "-o", "-f", family["name"], "route", "show").Output()
		if err != nil {
			return err
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			// default via 192.168.1.1 dev wlan0  proto static
			// fd01:37b7:3570::/64 dev wlan0  proto kernel  metric 256  expires 6837sec mtu 1280
			a := strings.Fields(line)
			if len(a) >= 1 {
				m := map[string]string{
					"destination": a[0],
					"family":      family["name"],
				}
				for i := range a {
					switch a[i] {
					case "via":
						m["via"] = a[i+1]
					case "src":
						m["src"] = a[i+1]
					case "proto":
						m["proto"] = a[i+1]
					case "metric":
						m["metric"] = a[i+1]
					case "scope":
						m["scope"] = a[i+1]

					}

				}

				routes = append(routes, m)
			}

		}
	}
	d["routes"] = routes
	return nil
}
