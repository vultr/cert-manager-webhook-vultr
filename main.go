package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/cmd"
	cmMeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/jetstack/cert-manager/pkg/issuer/acme/dns/util"
	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8Meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GroupName ...
var GroupName = os.Getenv("GROUP_NAME")

const version = "v0.3.0"

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
	k8Client    *kubernetes.Clientset
	vultrClient *govultr.Client
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
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	if v.vultrClient == nil {
		if err := v.setVultrClient(ch, cfg); err != nil {
			return err
		}
	}

	zoneName, err := util.FindZoneByFqdn(ch.ResolvedFQDN, util.RecursiveNameservers)
	if err != nil {
		return nil
	}

	records, err := v.getRecords(ch)
	if err != nil {
		return err
	}

	for _, v := range records {
		if v.Type == "TXT" && v.Data == fmt.Sprintf("\"%s\"", ch.Key) {
			return nil
		}
	}

	req := &govultr.DomainRecordReq{
		Name: v.stripZone(ch.ResolvedFQDN, zoneName),
		Type: "TXT",
		Data: ch.Key,
		TTL:  60,
	}

	_, err = v.vultrClient.DomainRecord.Create(context.Background(), util.UnFqdn(zoneName), req)
	if err != nil {
		return err
	}

	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
func (v *VultrSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	if v.vultrClient == nil {
		if err := v.setVultrClient(ch, cfg); err != nil {
			return err
		}
	}

	zoneName, err := util.FindZoneByFqdn(ch.ResolvedFQDN, util.RecursiveNameservers)
	if err != nil {
		return nil
	}

	records, err := v.getRecords(ch)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Type == "TXT" && record.Data == fmt.Sprintf("\"%s\"", ch.Key) {
			if err := v.vultrClient.DomainRecord.Delete(context.Background(), util.UnFqdn(zoneName), record.ID); err != nil {
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

func (v *VultrSolver) setVultrClient(ch *v1alpha1.ChallengeRequest, cfg VultrProviderConfig) error {
	ref := cfg.APIKeySecretRef
	if ref.Name == "" || ref.Key == "" {
		return fmt.Errorf("key not set in secret : %s", ref.Name)
	}

	secret, err := v.k8Client.CoreV1().Secrets(ch.ResourceNamespace).Get(context.Background(), ref.Name, k8Meta.GetOptions{})
	if err != nil {
		return err
	}

	keyBytes, ok := secret.Data[ref.Key]
	if !ok {
		return fmt.Errorf("no key %s in secret : %s'", ref.Key, ref.Name)
	}

	config := &oauth2.Config{}
	ctx := context.Background()
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: string(keyBytes)})
	v.vultrClient = govultr.NewClient(oauth2.NewClient(ctx, ts))
	v.vultrClient.SetUserAgent(fmt.Sprintf("cert-manager-webhook-vultr/%s", version))

	return nil
}

func (v *VultrSolver) getRecords(ch *v1alpha1.ChallengeRequest) ([]govultr.DomainRecord, error) {
	zone, err := util.FindZoneByFqdn(ch.ResolvedFQDN, util.RecursiveNameservers)
	if err != nil {
		return nil, err
	}

	test := util.UnFqdn(zone)
	fmt.Println(test)
	var records []govultr.DomainRecord
	//todo fill in the list options + meta
	recordsList, _, err := v.vultrClient.DomainRecord.List(context.Background(), test, nil)
	if err != nil {
		return nil, err
	}

	targetName := v.stripZone(ch.ResolvedFQDN, zone)
	for _, record := range recordsList {
		if record.Name == targetName {
			records = append(records, record)
		}
	}

	return records, err
}

func (v *VultrSolver) stripZone(resolvedFQDN, zone string) string {
	targetName := resolvedFQDN
	if strings.HasSuffix(resolvedFQDN, zone) {
		targetName = resolvedFQDN[:len(resolvedFQDN)-len(zone)-1]
	}
	return targetName
}
