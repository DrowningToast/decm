package entity

// CertificatePayloadHeader mirrors CertificateVCStructs.Header.
// @context and type are hardcoded by the contract:
//
//	@context: ["https://www.w3.org/2018/credentials/v1"]
//	type:     ["VerifiableCredential", "EventCertificate"]
type CertificatePayloadHeader struct {
	Context      []string `json:"@context"`
	Type         []string `json:"type"`
	Id           string   `json:"id"`
	Issuer       string   `json:"issuer"`
	IssuanceDate string   `json:"issuanceDate"`
}

// CertificatePayloadData mirrors CertificateVCStructs.Data.
// EncryptedUserData / BackendEncryptedUserData: full attendee profile JSON
// (8 fields from event_attendees), ECIES-encrypted with user / backend public key.
// IssuerAddresses: comma-separated string (Solidity limitation, not a JSON array).
type CertificatePayloadData struct {
	EventName                string `json:"eventName"`
	EventDescription         string `json:"eventDescription"`
	CertificateTokenId       string `json:"certificateTokenId"`
	CertificateId            string `json:"certificateId"`
	UserId                   string `json:"userId"`
	IssuerId                 string `json:"issuerId"`
	IssuedAt                 string `json:"issuedAt"`
	IssuerAddresses          string `json:"issuerAddresses"`
	ReceiverAddress          string `json:"receiverAddress"`
	EncryptedUserData        string `json:"encryptedUserData"`
	BackendEncryptedUserData string `json:"backendEncryptedUserData"`
	CertificateTitle         string `json:"certificateTitle"`
	CertificateSubtitle      string `json:"certificateSubtitle"`
	Status                   string `json:"status"` // "VALID" | "REVOKED"
}

// CertificatePayloadHostProof mirrors CertificateVCStructs.HostProof.
type CertificatePayloadHostProof struct {
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

// CertificatePayloadIssuerProof mirrors CertificateVCStructs.IssuerProof
// and the generated Go binding type CertificateVCStructsIssuerProof.
type CertificatePayloadIssuerProof struct {
	IssuerSignature string `json:"issuerSignature"`
	IssuerPublicKey string `json:"issuerPublicKey"`
}

// CertificatePayloadProof mirrors CertificateVCStructs.Proof.
// EncryptedByUserRawData / EncryptedByBackendRawData: certificate PII CSV
// (name,academic_institution,certificate_title,certificate_subtitle),
// ECIES-encrypted with user / backend public key.
// Hash: Keccak256 of the PII CSV (same value as certificate_digest in DB).
// SignMessage: participant's claim message "<walletAddr>,<contractAddr>,<deadlineBlock>".
type CertificatePayloadProof struct {
	EncryptedByUserRawData    string                          `json:"encryptedByUserRawData"`
	EncryptedByBackendRawData string                          `json:"encryptedByBackendRawData"`
	Hash                      string                          `json:"hash"`
	SignMessage               string                          `json:"signMessage"`
	Host                      CertificatePayloadHostProof     `json:"host"`
	Issuers                   []CertificatePayloadIssuerProof `json:"issuers"`
}

// AttendeeProfileData is the attendee PII that is ECIES-encrypted into
// BackendEncryptedUserData (backend public key) and EncryptedUserData (user public key)
// at certificate claim time.
type AttendeeProfileData struct {
	FirstName           *string `json:"first_name"`
	LastName            *string `json:"last_name"`
	Email               *string `json:"email"`
	Bio                 *string `json:"bio"`
	PhoneNumber         *string `json:"phone_number"`
	Address             *string `json:"address"`
	AcademicInstitution *string `json:"academic_institution"`
	AcademicEmail       *string `json:"academic_email"`
}

// CertificateRawData is the typed representation of the PII CSV stored in
// proof.EncryptedByBackendRawData. The CSV format is:
//
//	"{name},{academic_institution},{certificate_title},{certificate_subtitle}"
//
// where name = "{firstName} {lastName}" (space-joined, as written at claim time).
type CertificateRawData struct {
	Name                *string `json:"name"`
	AcademicInstitution *string `json:"academic_institution"`
	CertificateTitle    *string `json:"certificate_title"`
	CertificateSubtitle *string `json:"certificate_subtitle"`
}

// CertificatePayload is the fully typed domain representation of the on-chain VC.
// Returned by CertificateContractDataGateway.GetTokenData after stripping the
// "data:application/json;utf8," prefix and unmarshalling the JSON.
type CertificatePayload struct {
	Header CertificatePayloadHeader `json:"header"`
	Data   CertificatePayloadData   `json:"data"`
	Proof  CertificatePayloadProof  `json:"proof"`
}
