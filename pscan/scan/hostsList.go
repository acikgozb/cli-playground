// Package scan provides types and functions to perform TCP port
// scans on a list of hosts.
package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrExists    = errors.New("host already in the list")
	ErrNotExists = errors.New("host not in the list")
)

// HostList represents a list of hosts to run port scan
type HostsList struct {
	Hosts []string
}

func (hl *HostsList) search(host string) (bool, int) {
	sort.Strings(hl.Hosts)

	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}

	return false, -1
}

func (hl *HostsList) Add(host string) error {
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w:%s", ErrExists, host)
	}

	hl.Hosts = append(hl.Hosts, host)
	return nil
}

func (hl *HostsList) Remove(host string) error {
	found, i := hl.search(host)
	if !found {
		return fmt.Errorf("%w:%s", ErrNotExists, host)
	}

	hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
	return nil
}

func (hl *HostsList) Load(hostsFile string) error {
	f, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}

	return nil
}

func (hl *HostsList) Save(hostsFile string) error {
	output := ""

	for _, host := range hl.Hosts {
		output += fmt.Sprintln(host)
	}

	return os.WriteFile(hostsFile, []byte(output), 0644)
}
