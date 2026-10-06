package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	cmMeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"github.com/vultr/govultr/v3"
	"golang.org/x/oauth2"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8Meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GroupName ...
var GroupName = os.Getenv("GROUP_NAME")

const version = "v0.4.3"

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	// This will register our Vultr provider with the webhook serving
	// library, making it available as an API under the provided GroupName.
	cmd.RunWebhookServer(GroupName,
		&VultrSolver{},
	)
}

// VultrSolver implements the provider-specific logic needed to
// 'present' an ACME challenge TXT record for your own DNS provider.
type VultrSolver struct {
	k8Client *kubernetes.Clientset
}

// VultrProviderConfig is a structure that is used to decode into when
// solving a DNS01 challenge.
type VultrProviderConfig struct {
	APIKeySecretRef cmMeta.SecretKeySelector `json:"apiKeySecretRef"`
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource.
func (v *VultrSolver) Name() string {
	return "vultr"
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
func (v *VultrSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	vultrClient, err := v.newVultrClient(ctx, ch, cfg)
	if err != nil {
		return err
	}

	zoneName, err := util.FindZoneByFqdn(ctx, ch.ResolvedFQDN, util.RecursiveNameservers)
	if err != nil {
		return err
	}

	records, err := v.getRecords(ctx, ch, zoneName, vultrClient)
	if err != nil {
		return err
	}

	for _, v := range records {
		if v.Type == "TXT" && v.Data == fmt.Sprintf("\"%s\"", ch.Key) {
			return nil
		}
	}

	req := &govultr.DomainRecordCreateReq{
		Name: v.stripZone(ch.ResolvedFQDN, zoneName),
		Type: "TXT",
		Data: ch.Key,
		TTL:  60,
	}

	_, _, err = vultrClient.DomainRecord.Create(ctx, util.UnFqdn(zoneName), req)
	if err != nil {
		return err
	}

	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
func (v *VultrSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	vultrClient, err := v.newVultrClient(ctx, ch, cfg)
	if err != nil {
		return err
	}

	zoneName, err := util.FindZoneByFqdn(ctx, ch.ResolvedFQDN, util.RecursiveNameservers)
	if err != nil {
		return err
	}

	records, err := v.getRecords(ctx, ch, zoneName, vultrClient)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Type == "TXT" && record.Data == fmt.Sprintf("\"%s\"", ch.Key) {
			if err := vultrClient.DomainRecord.Delete(ctx, util.UnFqdn(zoneName), record.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

// Initialize will be called when the webhook first starts.
func (v *VultrSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}
	v.k8Client = cl
	return nil
}

// loadConfig is a small helper function that decodes JSON configuration into the typed config struct.
func loadConfig(cfgJSON *extapi.JSON) (VultrProviderConfig, error) {
	cfg := VultrProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	return cfg, nil
}

func (v *VultrSolver) newVultrClient(ctx context.Context, ch *v1alpha1.ChallengeRequest, cfg VultrProviderConfig) (*govultr.Client, error) {
	ref := cfg.APIKeySecretRef
	if ref.Name == "" || ref.Key == "" {
		return nil, fmt.Errorf("apiKeySecretRef.name and apiKeySecretRef.key must both be set")
	}

	secret, err := v.k8Client.CoreV1().Secrets(ch.ResourceNamespace).Get(ctx, ref.Name, k8Meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	keyBytes, ok := secret.Data[ref.Key]
	if !ok {
		return nil, fmt.Errorf("no key %s in secret %s", ref.Key, ref.Name)
	}

	config := &oauth2.Config{}
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: string(keyBytes)})
	vultrClient := govultr.NewClient(oauth2.NewClient(ctx, ts))
	vultrClient.SetUserAgent(fmt.Sprintf("cert-manager-webhook-vultr/%s", version))

	return vultrClient, nil
}

func (v *VultrSolver) getRecords(ctx context.Context, ch *v1alpha1.ChallengeRequest, zone string, vultrClient *govultr.Client) ([]govultr.DomainRecord, error) {
	domain := util.UnFqdn(zone)
	log.Printf("[DEBUG] Looking up records for domain: %s", domain)

	var records []govultr.DomainRecord
	targetName := v.stripZone(ch.ResolvedFQDN, zone)
	listOptions := &govultr.ListOptions{PerPage: 100}
	for {
		recordsList, meta, _, err := vultrClient.DomainRecord.List(ctx, domain, listOptions)
		if err != nil {
			return nil, err
		}

		for _, record := range recordsList {
			if record.Name == targetName {
				records = append(records, record)
			}
		}
		if meta.Links.Next == "" {
			break
		}

		listOptions.Cursor = meta.Links.Next
	}

	return records, nil
}

func (v *VultrSolver) stripZone(resolvedFQDN, zone string) string {
	fqdn := util.UnFqdn(resolvedFQDN)
	zone = util.UnFqdn(zone)
	if fqdn == zone {
		return ""
	}
	if strings.HasSuffix(fqdn, "."+zone) {
		return strings.TrimSuffix(fqdn, "."+zone)
	}
	return fqdn
}
