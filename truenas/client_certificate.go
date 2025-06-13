package truenas

import (
	"context"
	"fmt"
	"time"
)

// CertificateClient provides methods for certificate management
type CertificateClient struct {
	client *Client
}

// NewCertificateClient creates a new certificate client
func NewCertificateClient(client *Client) *CertificateClient {
	return &CertificateClient{client: client}
}

// Certificate represents a TLS/SSL certificate
type Certificate struct {
	ID                 int            `json:"id"`
	Type               int            `json:"type"`
	Name               string         `json:"name"`
	Certificate        string         `json:"certificate"`
	Privatekey         string         `json:"privatekey"`
	CSR                string         `json:"CSR"`
	Acme               any            `json:"acme"`
	CertificateUsed    bool           `json:"certificate_used"`
	Revoked            bool           `json:"revoked"`
	Internal           string         `json:"internal"`
	CA                 bool           `json:"ca"`
	Cert               string         `json:"cert"`
	Chain              []string       `json:"chain"`
	Country            string         `json:"country"`
	State              string         `json:"state"`
	City               string         `json:"city"`
	Organization       string         `json:"organization"`
	OrganizationalUnit string         `json:"organizational_unit"`
	Common             string         `json:"common"`
	SAN                []string       `json:"san"`
	Email              string         `json:"email"`
	DN                 string         `json:"DN"`
	Subject            map[string]any `json:"subject"`
	Extensions         map[string]any `json:"extensions"`
	NotBefore          time.Time      `json:"not_before"`
	NotAfter           time.Time      `json:"not_after"`
	Issuer             map[string]any `json:"issuer"`
	Digest             string         `json:"digest"`
	Serial             int            `json:"serial"`
	KeyLength          int            `json:"key_length"`
	KeyType            string         `json:"key_type"`
	Fingerprint        string         `json:"fingerprint"`
	RootPath           string         `json:"root_path"`
	CertType           string         `json:"cert_type"`
	CertTypeFull       string         `json:"cert_type_full"`
	SignedBy           any            `json:"signed_by"`
	Lifetime           int            `json:"lifetime"`
	From               string         `json:"from"`
	Parsed             bool           `json:"parsed"`
}

// CertificateExtensions represents X.509v3 certificate extensions
type CertificateExtensions struct {
	BasicConstraints       *BasicConstraints       `json:"BasicConstraints,omitempty"`
	AuthorityKeyIdentifier *AuthorityKeyIdentifier `json:"AuthorityKeyIdentifier,omitempty"`
	ExtendedKeyUsage       *ExtendedKeyUsage       `json:"ExtendedKeyUsage,omitempty"`
	KeyUsage               *KeyUsage               `json:"KeyUsage,omitempty"`
}

// BasicConstraints represents basic constraints extension
type BasicConstraints struct {
	CA                bool `json:"ca"`
	Enabled           bool `json:"enabled"`
	PathLength        *int `json:"path_length"`
	ExtensionCritical bool `json:"extension_critical"`
}

// AuthorityKeyIdentifier represents authority key identifier extension
type AuthorityKeyIdentifier struct {
	AuthorityCertIssuer bool `json:"authority_cert_issuer"`
	Enabled             bool `json:"enabled"`
	ExtensionCritical   bool `json:"extension_critical"`
}

// ExtendedKeyUsage represents extended key usage extension
type ExtendedKeyUsage struct {
	Usages            []string `json:"usages"`
	Enabled           bool     `json:"enabled"`
	ExtensionCritical bool     `json:"extension_critical"`
}

// KeyUsage represents key usage extension
type KeyUsage struct {
	Enabled           bool `json:"enabled"`
	DigitalSignature  bool `json:"digital_signature"`
	ContentCommitment bool `json:"content_commitment"`
	KeyEncipherment   bool `json:"key_encipherment"`
	DataEncipherment  bool `json:"data_encipherment"`
	KeyAgreement      bool `json:"key_agreement"`
	KeyCertSign       bool `json:"key_cert_sign"`
	CRLSign           bool `json:"crl_sign"`
	EncipherOnly      bool `json:"encipher_only"`
	DecipherOnly      bool `json:"decipher_only"`
	ExtensionCritical bool `json:"extension_critical"`
}

// CertificateCreateRequest represents parameters for certificate.create
type CertificateCreateRequest struct {
	// Common fields
	Name       string `json:"name"`
	CreateType string `json:"create_type"`

	// Internal certificate fields
	KeyLength          int                    `json:"key_length,omitempty"`
	KeyType            string                 `json:"key_type,omitempty"`
	ECCurve            string                 `json:"ec_curve,omitempty"`
	DigestAlgorithm    string                 `json:"digest_algorithm,omitempty"`
	Lifetime           int                    `json:"lifetime,omitempty"`
	Country            string                 `json:"country,omitempty"`
	State              string                 `json:"state,omitempty"`
	City               string                 `json:"city,omitempty"`
	Organization       string                 `json:"organization,omitempty"`
	OrganizationalUnit string                 `json:"organizational_unit,omitempty"`
	Email              string                 `json:"email,omitempty"`
	Common             string                 `json:"common,omitempty"`
	SAN                []string               `json:"san,omitempty"`
	SignedBy           int                    `json:"signedby,omitempty"`
	CertExtensions     *CertificateExtensions `json:"cert_extensions,omitempty"`

	// Import certificate fields
	Certificate string `json:"certificate,omitempty"`
	Privatekey  string `json:"privatekey,omitempty"`
	Passphrase  string `json:"passphrase,omitempty"`

	// CSR fields
	CSR string `json:"CSR,omitempty"`

	// ACME fields
	TOS              bool              `json:"tos,omitempty"`
	CSRID            int               `json:"csr_id,omitempty"`
	AcmeDirectoryURI string            `json:"acme_directory_uri,omitempty"`
	DNSMapping       map[string]string `json:"dns_mapping,omitempty"`
	RenewDays        int               `json:"renew_days,omitempty"`

	// Other fields
	Type   int `json:"type,omitempty"`
	Serial int `json:"serial,omitempty"`
}

// CertificateUpdateRequest represents parameters for certificate.update
type CertificateUpdateRequest struct {
	Name    string `json:"name,omitempty"`
	Revoked bool   `json:"revoked,omitempty"`
}

// Certificate create types
type CertificateCreateType string

const (
	CertificateCreateInternal    CertificateCreateType = "CERTIFICATE_CREATE_INTERNAL"
	CertificateCreateImported    CertificateCreateType = "CERTIFICATE_CREATE_IMPORTED"
	CertificateCreateCSR         CertificateCreateType = "CERTIFICATE_CREATE_CSR"
	CertificateCreateImportedCSR CertificateCreateType = "CERTIFICATE_CREATE_IMPORTED_CSR"
	CertificateCreateACME        CertificateCreateType = "CERTIFICATE_CREATE_ACME"
)

// Key types
type CertificateKeyType string

const (
	CertificateKeyTypeRSA CertificateKeyType = "RSA"
	CertificateKeyTypeEC  CertificateKeyType = "EC"
)

// EC curves
type CertificateECCurve string

const (
	CertificateECCurveBrainpoolP512R1 CertificateECCurve = "BrainpoolP512R1"
	CertificateECCurveBrainpoolP384R1 CertificateECCurve = "BrainpoolP384R1"
	CertificateECCurveBrainpoolP256R1 CertificateECCurve = "BrainpoolP256R1"
	CertificateECCurveSECP256K1       CertificateECCurve = "SECP256K1"
	CertificateECCurveSECP384R1       CertificateECCurve = "SECP384R1"
	CertificateECCurveSECP521R1       CertificateECCurve = "SECP521R1"
	CertificateECCurveEd25519         CertificateECCurve = "ed25519"
)

// Digest algorithms
type CertificateDigestAlgorithm string

const (
	CertificateDigestSHA1   CertificateDigestAlgorithm = "SHA1"
	CertificateDigestSHA224 CertificateDigestAlgorithm = "SHA224"
	CertificateDigestSHA256 CertificateDigestAlgorithm = "SHA256"
	CertificateDigestSHA384 CertificateDigestAlgorithm = "SHA384"
	CertificateDigestSHA512 CertificateDigestAlgorithm = "SHA512"
)

// List returns all certificates
func (c *CertificateClient) List(ctx context.Context) ([]Certificate, error) {
	var result []Certificate
	err := c.client.Call(ctx, "certificate.query", []any{}, &result)
	return result, err
}

// Get returns a specific certificate by ID
func (c *CertificateClient) Get(ctx context.Context, id int) (*Certificate, error) {
	var result []Certificate
	err := c.client.Call(ctx, "certificate.query", []any{[]any{[]any{"id", "=", id}}}, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, NewNotFoundError("certificate", fmt.Sprintf("ID %d", id))
	}
	return &result[0], nil
}

// Create creates a new certificate
func (c *CertificateClient) Create(ctx context.Context, req *CertificateCreateRequest) (*Certificate, error) {
	var result Certificate
	err := c.client.CallJob(ctx, "certificate.create", []any{*req}, &result)
	return &result, err
}

// Update updates an existing certificate
func (c *CertificateClient) Update(ctx context.Context, id int, req *CertificateUpdateRequest) (*Certificate, error) {
	var result Certificate
	err := c.client.CallJob(ctx, "certificate.update", []any{id, *req}, &result)
	return &result, err
}

// Delete deletes a certificate
func (c *CertificateClient) Delete(ctx context.Context, id int, force bool) error {
	return c.client.CallJob(ctx, "certificate.delete", []any{id, force}, nil)
}

// Certificate Configuration Choices

// GetCountryChoices returns available country choices for certificates
func (c *CertificateClient) GetCountryChoices(ctx context.Context) (map[string]string, error) {
	var result map[string]string
	err := c.client.Call(ctx, "certificate.country_choices", []any{}, &result)
	return result, err
}

// GetKeyTypeChoices returns supported key types for certificates
func (c *CertificateClient) GetKeyTypeChoices(ctx context.Context) (map[string]string, error) {
	var result map[string]string
	err := c.client.Call(ctx, "certificate.key_type_choices", []any{}, &result)
	return result, err
}

// GetECCurveChoices returns supported EC curves
func (c *CertificateClient) GetECCurveChoices(ctx context.Context) (map[string]string, error) {
	var result map[string]string
	err := c.client.Call(ctx, "certificate.ec_curve_choices", []any{}, &result)
	return result, err
}

// GetExtendedKeyUsageChoices returns choices for ExtendedKeyUsage extension
func (c *CertificateClient) GetExtendedKeyUsageChoices(ctx context.Context) (map[string]string, error) {
	var result map[string]string
	err := c.client.Call(ctx, "certificate.extended_key_usage_choices", []any{}, &result)
	return result, err
}

// GetProfiles returns predefined certificate profiles for specific use cases
func (c *CertificateClient) GetProfiles(ctx context.Context) (map[string]any, error) {
	var result map[string]any
	err := c.client.Call(ctx, "certificate.profiles", []any{}, &result)
	return result, err
}

// GetACMEServerChoices returns popular ACME servers with their directory URIs
func (c *CertificateClient) GetACMEServerChoices(ctx context.Context) (map[string]string, error) {
	var result map[string]string
	err := c.client.Call(ctx, "certificate.acme_server_choices", []any{}, &result)
	return result, err
}
