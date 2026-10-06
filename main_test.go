package main

import (
	"testing"
)

func TestStripZone(t *testing.T) {
	tests := []struct {
		name string
		fqdn string
		zone string
		want string
	}{
		{name: "challenge below zone", fqdn: "_acme-challenge.example.com.", zone: "example.com.", want: "_acme-challenge"},
		{name: "nested challenge below zone", fqdn: "_acme-challenge.foo.example.com.", zone: "example.com.", want: "_acme-challenge.foo"},
		{name: "zone suffix is label bounded", fqdn: "_acme-challenge.notexample.com.", zone: "example.com.", want: "_acme-challenge.notexample.com"},
	}

	solver := &VultrSolver{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := solver.stripZone(tt.fqdn, tt.zone); got != tt.want {
				t.Fatalf("stripZone(%q, %q) = %q, want %q", tt.fqdn, tt.zone, got, tt.want)
			}
		})
	}
}
