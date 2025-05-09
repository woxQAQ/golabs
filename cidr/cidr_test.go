package cidr

import (
	"net"
	"testing"
)

func TestAddGroup(t *testing.T) {
	tests := []struct {
		name    string
		cidr    string
		wantErr bool
	}{
		{"valid cidr", "192.168.1.0/24", false},
		{"invalid cidr", "invalid.cidr", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AddGroup(tt.cidr, "test-group")
			if (err != nil) != tt.wantErr {
				t.Errorf("AddGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindGroup(t *testing.T) {
	// Setup test data
	_, err := AddGroup("10.0.0.0/8", "group-a")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []struct {
		name    string
		ip      string
		want    string
		wantErr bool
	}{
		{"found in group", "10.1.2.3", "group-a", false},
		{"not found", "192.168.1.1", "", false},
		{"invalid ip", "invalid.ip", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindGroup(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.name != tt.want {
				t.Errorf("FindGroup() = %v, want %v", got.name, tt.want)
			}
		})
	}
}

func TestGetIpUint32(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want uint32
	}{
		{"ipv4", "192.168.1.1", 3232235777},
		{"ipv4", "10.0.0.1", 167772161},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if got := getIpUint32(ip); got != tt.want {
				t.Errorf("getIpUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}
