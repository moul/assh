package config

import (
	"sort"
	"strings"
)

// HostsMap is a map of **Host).Name -> *Host
type HostsMap map[string]*Host

// HostsList is a list of *Host
type HostsList []*Host

// ToList returns a slice of *Hosts
func (hm *HostsMap) ToList() HostsList {
	list := HostsList{}
	for _, host := range *hm {
		list = append(list, host)
	}
	return list
}

func (hl HostsList) Len() int           { return len(hl) }
func (hl HostsList) Swap(i, j int)      { hl[i], hl[j] = hl[j], hl[i] }
func (hl HostsList) Less(i, j int) bool { return strings.Compare(hl[i].name, hl[j].name) < 0 }

// SortedList returns a list of hosts sorted by their name
func (hm *HostsMap) SortedList() HostsList {
	sortedList := hm.ToList()
	sort.Sort(sortedList)
	return sortedList
}
