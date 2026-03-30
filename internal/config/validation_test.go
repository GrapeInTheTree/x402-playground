package config

import "testing"

func TestValidateEthAddress(t *testing.T) {
	tests := []struct {
		addr    string
		wantErr bool
	}{
		{"0x70997970C51812dc3A010C7d01b50e0d17dc79C8", false},
		{"0x036CbD53842c5426634e7929541eC2318f3dCF7e", false},
		{"0x1234", true},                      // too short
		{"70997970C51812dc3A010C7d01b50e0d17dc79C8", true}, // missing 0x
		{"0xGGGG970C51812dc3A010C7d01b50e0d17dc79C8", true}, // invalid hex
		{"", true},
	}
	for _, tt := range tests {
		err := validateEthAddress(tt.addr, "TEST_ADDR")
		if (err != nil) != tt.wantErr {
			t.Errorf("validateEthAddress(%q) error=%v, wantErr=%v", tt.addr, err, tt.wantErr)
		}
	}
}

func TestValidateNetwork(t *testing.T) {
	tests := []struct {
		network string
		wantErr bool
	}{
		{"eip155:84532", false},
		{"eip155:8453", false},
		{"eip155:1", false},
		{"eip155:", true},
		{"eip155:abc", true},
		{"invalid", true},
		{"", true},
	}
	for _, tt := range tests {
		err := validateNetwork(tt.network)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateNetwork(%q) error=%v, wantErr=%v", tt.network, err, tt.wantErr)
		}
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"http://localhost:4022", false},
		{"https://sepolia.base.org", false},
		{"ftp://example.com", true},
		{"not-a-url", true},
		{"", true},
	}
	for _, tt := range tests {
		err := validateURL(tt.url, "TEST_URL")
		if (err != nil) != tt.wantErr {
			t.Errorf("validateURL(%q) error=%v, wantErr=%v", tt.url, err, tt.wantErr)
		}
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		port    string
		wantErr bool
	}{
		{"4022", false},
		{"1", false},
		{"65535", false},
		{"0", true},
		{"65536", true},
		{"abc", true},
		{"", true},
	}
	for _, tt := range tests {
		err := validatePort(tt.port)
		if (err != nil) != tt.wantErr {
			t.Errorf("validatePort(%q) error=%v, wantErr=%v", tt.port, err, tt.wantErr)
		}
	}
}

func TestValidateTransferMethod(t *testing.T) {
	tests := []struct {
		method  string
		wantErr bool
	}{
		{"eip3009", false},
		{"permit2", false},
		{"invalid", true},
		{"", true},
	}
	for _, tt := range tests {
		err := validateTransferMethod(tt.method)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateTransferMethod(%q) error=%v, wantErr=%v", tt.method, err, tt.wantErr)
		}
	}
}
