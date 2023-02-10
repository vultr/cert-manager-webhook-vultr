package main

import (
	"os"
	"testing"
	"time"

	"github.com/cert-manager/cert-manager/test/acme/dns"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	fixture := dns.NewFixture(&VultrSolver{},
		dns.SetResolvedZone(zone),
		dns.SetAllowAmbientCredentials(false),
		dns.SetManifestPath("testdata/vultr"),
		dns.SetDNSName(zone),
		dns.SetPropagationLimit(time.Minute*20),
	)
	fixture.RunConformance(t)
}
