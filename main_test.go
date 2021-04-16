package main

import (
	"os"
	"testing"

	"github.com/jetstack/cert-manager/test/acme/dns"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	fixture := dns.NewFixture(&VultrSolver{},
		dns.SetResolvedZone(zone),
		dns.SetAllowAmbientCredentials(false),
		dns.SetManifestPath("testdata/vultr"),
		dns.SetBinariesPath("_test/kubebuilder/bin"),
		dns.SetResolvedZone("freightapp.co."),
	)
	fixture.RunConformance(t)
}
