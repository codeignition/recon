// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// +build linux

// Package top provides selective data provided by the `top` command.
// It collects system summary data, current running tasks, etc.
package top

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Data represents the `top` data
type Data struct {
	Uptime  string
	CPU     CPUData
	Memory  MemoryData
	Swap    SwapData
	LoadAvg LoadAverageData
}

// CPU represents the percentage CPU usages, averaged over the cores
type CPUData struct {
	Userspace,
	Idle,
	System,
	IOWait,
	Stolen float64
}

// LoadAverageData for the last 1min, 5mins, 15mins
type LoadAverageData struct {
	LastOneMin,
	LastFiveMin,
	LastFifteenMin float64
}

// MemoryData represents the memory data in kilobytes
type MemoryData struct {
	Total,
	Used,
	Free,
	Buffers,
	Cached int
}

// SwapData represents the swap memory data in kilobytes
type SwapData struct {
	Total,
	Used,
	Free int
}

// CollectData collects the data and returns
// an error if any.
func CollectData() (*Data, error) {
	d := new(Data)
	count := 0 // count of the top output iteration
	iters := 2 // number of iterations that top command should collect
	out, err := exec.Command("top", "-bn", strconv.Itoa(iters)).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n") // use a bufio.Scanner if memory problems arise
	if len(lines) < 7 {
		return nil, errors.New("top: unexpected output")
	}
	base := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "top") {
			count++
			base = i
		}
		if count == iters {
			switch i - base {
			case 0:
				err := d.parseUptimeLoadAvgData(line)
				if err != nil {
					return nil, err
				}
			case 1:
				// we are not collecting the tasks data for now
				continue
			case 2:
				err := d.parseCPUData(line)
				if err != nil {
					return nil, err
				}
			case 3:
				err := d.parseMemoryData(line)
				if err != nil {
					return nil, err
				}
			case 4:
				err := d.parseSwapData(line)
				if err != nil {
					return nil, err
				}
			default:
				break
			}
		}
	}
	return d, nil
}

func (d *Data) parseUptimeLoadAvgData(s string) error {
	a := strings.SplitN(s, "load average: ", 2)
	b := strings.SplitN(a[0], "up", 2)
	t := strings.SplitN(b[1], "users", 2)
	c := strings.TrimSpace(t[0])
	i := strings.LastIndex(c, ",")
	if i == -1 {
		return errors.New("top: unable to parse uptime data")
	}
	d.Uptime = c[:i]
	l := strings.Split(a[1], ",")
	if len(l) != 3 {
		return errors.New("top: unexpected number of load averages found")
	}
	var f [3]float64
	for i := range l {
		x, err := strconv.ParseFloat(strings.TrimSpace(l[i]), 64)
		f[i] = x
		if err != nil {
			return err
		}
	}
	d.LoadAvg = LoadAverageData{
		LastOneMin:     f[0],
		LastFiveMin:    f[1],
		LastFifteenMin: f[2],
	}
	return nil
}

func (d *Data) parseCPUData(s string) error {
	a := strings.TrimPrefix(s, "%Cpu(s): ")
	a = strings.TrimPrefix(a, "Cpu(s): ") // some systems have the prefix Cpu(s): without the % sign
	b := strings.Split(a, ",")
	if len(b) != 8 {
		return errors.New("top: unknown number of CPU data")
	}
	var c [8]float64
	for i := range b {
		t := strings.TrimSpace(b[i])
		x, err := strconv.ParseFloat(t[:len(t)-3], 64)
		c[i] = x
		if err != nil {
			return err
		}
	}
	d.CPU = CPUData{
		Userspace: c[0],
		Idle:      c[3],
		System:    c[1],
		IOWait:    c[4],
		Stolen:    c[7],
	}
	return nil
}

func (d *Data) parseMemoryData(s string) error {
	if strings.HasPrefix(s, "KiB Mem:") {
		var m MemoryData
		_, err := fmt.Sscanf(s, "KiB Mem:\t%d total,\t%d used,\t%d free,\t%d buffers", &m.Total, &m.Used, &m.Free, &m.Buffers)
		if err != nil {
			return fmt.Errorf("top: unable to parse memory data; %s", err)
		}
		d.Memory = m
		return nil
	} else if strings.HasPrefix(s, "Mem:") {
		var m MemoryData
		_, err := fmt.Sscanf(s, "Mem:\t%dk total,\t%dk used,\t%dk free,\t%dk buffers", &m.Total, &m.Used, &m.Free, &m.Buffers)
		if err != nil {
			return fmt.Errorf("top: unable to parse memory data; %s", err)
		}
		d.Memory = m
		return nil
	}
	return errors.New("top: unable to parse memory data")
}

func (d *Data) parseSwapData(s string) error {
	if strings.HasPrefix(s, "KiB Swap:") {
		var sd SwapData
		_, err := fmt.Sscanf(s, "KiB Swap:\t%d total,\t%d used,\t%d free.\t%d cached Mem", &sd.Total, &sd.Used, &sd.Free, &d.Memory.Cached)
		if err != nil {
			return fmt.Errorf("top: unable to parse swap data; %s", err)
		}
		d.Swap = sd
		return nil
	} else if strings.HasPrefix(s, "Swap:") {
		var sd SwapData
		_, err := fmt.Sscanf(s, "Swap:\t%dk total,\t%dk used,\t%dk free,\t%dk cached", &sd.Total, &sd.Used, &sd.Free, &d.Memory.Cached)
		if err != nil {
			return fmt.Errorf("top: unable to parse swap data; %s", err)
		}
		d.Swap = sd
		return nil
	}
	return errors.New("top: unable to parse swap data")
}
