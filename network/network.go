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

	if err := arp(d); err != nil {
		return d, err
	}

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
	return nil
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
