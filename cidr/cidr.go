package cidr

import (
	"fmt"
	"maps"
	"net"
	"slices"
	"sync"
)

type Group struct {
	name string
}

type ipallocator struct {
	mu             sync.RWMutex
	subnetToGroups map[uint8]map[uint32]*Group
}

func NewIPAllocator() *ipallocator {
	return &ipallocator{
		subnetToGroups: make(map[uint8]map[uint32]*Group),
	}
}

type cidr interface {
	string | net.IPNet
}

var a *ipallocator

func init() {
	a = NewIPAllocator()
}

func AddGroup(cidr string, name string) (*Group, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	maskLen, _ := ipnet.Mask.Size()
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.subnetToGroups[uint8(maskLen)]; !ok {
		a.subnetToGroups[uint8(maskLen)] = make(map[uint32]*Group)
	}
	g := &Group{name}
	a.subnetToGroups[uint8(maskLen)][getIpUint32(ipnet.IP)] = g
	return g, nil
}

func FindGroup(ipStr string) (*Group, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid ip: %s", ipStr)
	}
	ipUint32 := getIpUint32(ip)
	a.mu.RLock()
	defer a.mu.RUnlock()
	keys := slices.Sorted(maps.Keys(a.subnetToGroups))
	slices.Reverse(keys)
	for _, key := range keys {
		if _, ok := a.subnetToGroups[key][ipUint32>>uint(32-key)]; ok {
			return a.subnetToGroups[key][ipUint32>>uint(32-key)], nil
		}
	}
	return nil, nil
}

func getIpUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
