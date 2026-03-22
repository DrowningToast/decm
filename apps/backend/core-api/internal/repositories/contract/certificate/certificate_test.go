package contract

import (
	"apps/backend/core-api/internal/entity"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixIssuerAddresses(t *testing.T) {
	t.Run("single address", func(t *testing.T) {
		input := `{"issuerAddresses": "["0xabc123def456abc123def456abc123def456abc1"]"}`
		want := `{"issuerAddresses": "0xabc123def456abc123def456abc123def456abc1"}`
		assert.Equal(t, want, fixIssuerAddresses(input))
	})

	t.Run("multiple addresses", func(t *testing.T) {
		input := `{"issuerAddresses": "["0xaaa","0xbbb","0xccc"]"}`
		want := `{"issuerAddresses": "0xaaa,0xbbb,0xccc"}`
		assert.Equal(t, want, fixIssuerAddresses(input))
	})

	t.Run("empty array becomes empty string", func(t *testing.T) {
		input := `{"issuerAddresses": "[]"}`
		want := `{"issuerAddresses": ""}`
		assert.Equal(t, want, fixIssuerAddresses(input))
	})

	t.Run("already valid comma-separated string is untouched", func(t *testing.T) {
		input := `{"issuerAddresses": "0xaaa,0xbbb"}`
		assert.Equal(t, input, fixIssuerAddresses(input))
	})

	t.Run("unrelated fields are untouched", func(t *testing.T) {
		input := `{"eventName": "Test","issuerAddresses": "["0xaaa"]","status": "VALID"}`
		want := `{"eventName": "Test","issuerAddresses": "0xaaa","status": "VALID"}`
		assert.Equal(t, want, fixIssuerAddresses(input))
	})

	t.Run("two addresses — real contract output that caused the 500", func(t *testing.T) {
		// Exact pattern the contract emits. Before the fix this caused:
		// "invalid character '0' after object key:value pair"
		input := `{"issuerAddresses": "["0x1111111111111111111111111111111111111111","0x2222222222222222222222222222222222222222"]"}`
		want := `{"issuerAddresses": "0x1111111111111111111111111111111111111111,0x2222222222222222222222222222222222222222"}`
		assert.Equal(t, want, fixIssuerAddresses(input))
	})
}

func TestFixFreeTextField(t *testing.T) {
	t.Run("no embedded quote — value untouched", func(t *testing.T) {
		input := `"eventName": "Hello World","eventDescription": "desc"`
		result := fixFreeTextField(input, `"eventName": "`, `","eventDescription":`)
		var m map[string]string
		require.NoError(t, json.Unmarshal([]byte("{"+result+"}"), &m))
		assert.Equal(t, "Hello World", m["eventName"])
	})

	t.Run("embedded quote followed by e — the exact 500 trigger", func(t *testing.T) {
		// eventName = `Hello "everyone` causes "invalid character 'e' after object key:value pair"
		// because the parser sees "Hello " as the value and then finds 'e'.
		input := `"eventName": "Hello "everyone","eventDescription": "desc"`
		result := fixFreeTextField(input, `"eventName": "`, `","eventDescription":`)
		var m map[string]string
		require.NoError(t, json.Unmarshal([]byte("{"+result+"}"), &m))
		assert.Equal(t, `Hello "everyone`, m["eventName"])
	})

	t.Run("multiple embedded quotes", func(t *testing.T) {
		input := `"eventName": "say "hello" world","eventDescription": "desc"`
		result := fixFreeTextField(input, `"eventName": "`, `","eventDescription":`)
		var m map[string]string
		require.NoError(t, json.Unmarshal([]byte("{"+result+"}"), &m))
		assert.Equal(t, `say "hello" world`, m["eventName"])
	})

	t.Run("field not present — string unchanged", func(t *testing.T) {
		input := `{"status": "VALID"}`
		assert.Equal(t, input, fixFreeTextField(input, `"eventName": "`, `","eventDescription":`))
	})
}

func TestFixFreeTextFields(t *testing.T) {
	t.Run("fixes all five free-text fields with embedded quotes", func(t *testing.T) {
		input := `"eventName": "Hello "everyone","eventDescription": "a "great" event","certificateTokenId": "1","certificateId": "c1","userId": "u1","issuerId": "i1","issuedAt": "2024-01-01","issuerAddresses": "0xaaa","receiverAddress": "0xbbb","encryptedUserData": "enc","backendEncryptedUserData": "benc","certificateTitle": "Title "with" quotes","certificateSubtitle": "Sub "title"","status": "VALID"`
		result := fixFreeTextFields(`{"data": {` + input + `}}`)

		var out struct {
			Data map[string]string `json:"data"`
		}
		require.NoError(t, json.Unmarshal([]byte(result), &out))
		assert.Equal(t, `Hello "everyone`, out.Data["eventName"])
		assert.Equal(t, `a "great" event`, out.Data["eventDescription"])
		assert.Equal(t, `Title "with" quotes`, out.Data["certificateTitle"])
		assert.Equal(t, `Sub "title"`, out.Data["certificateSubtitle"])
	})

	t.Run("signMessage embedded JSON — the actual 500 cause", func(t *testing.T) {
		// The contract stores signMessage as a raw JSON object without escaping, e.g.:
		// "signMessage": "{"eventContractAddress":"0x...","receivers":["0x..."]}",
		// Parser reads `"signMessage": "{"` → string ends → sees 'e' → 500.
		signMsg := `{"eventContractAddress":"0x428e6c4258cA560bbb56D83F0cfbAD8dc33f6696","receivers":["0xabc"]}`
		input := `"signMessage": "` + signMsg + `","host": {"signature": "0xsig","publicKey": "0xpk"}`
		result := fixFreeTextFields(`{"proof": {` + input + `}}`)

		var out struct {
			Proof struct {
				SignMessage string `json:"signMessage"`
			} `json:"proof"`
		}
		require.NoError(t, json.Unmarshal([]byte(result), &out))
		assert.Equal(t, signMsg, out.Proof.SignMessage)
	})
}

func TestGetTokenData_EmbeddedQuoteInEventName_Parses(t *testing.T) {
	// Simulates the 500 that occurred when eventName contains a `"e` sequence:
	// "invalid character 'e' after object key:value pair"
	contractOutput := strings.Join([]string{
		tokenDataPrefix,
		`{`,
		`"header": {`,
		`"@context": ["https://www.w3.org/2018/credentials/v1"],`,
		`"type": ["VerifiableCredential","EventCertificate"],`,
		`"id": "cert-1",`,
		`"issuer": "0xissuer",`,
		`"issuanceDate": "2024-01-01T00:00:00Z"`,
		`},`,
		`"data": {`,
		`"eventName": "Hello "everyone",`,
		`"eventDescription": "desc",`,
		`"certificateTokenId": "1",`,
		`"certificateId": "cert-1",`,
		`"userId": "user-1",`,
		`"issuerId": "issuer-1",`,
		`"issuedAt": "2024-01-01T00:00:00Z",`,
		`"issuerAddresses": "0xabc",`,
		`"receiverAddress": "0xreceiver",`,
		`"encryptedUserData": "encdata",`,
		`"backendEncryptedUserData": "backendencdata",`,
		`"certificateTitle": "Certificate of Completion",`,
		`"certificateSubtitle": "subtitle",`,
		`"status": "VALID"`,
		`},`,
		`"proof": {`,
		`"encryptedByUserRawData": "userraw",`,
		`"encryptedByBackendRawData": "backendraw",`,
		`"hash": "0xhash",`,
		`"signMessage": "0xwallet,0xcontract,12345",`,
		`"host": {"signature": "0xhostsig","publicKey": "hostpubkey"},`,
		`"issuers": [{"issuerSignature": "0xissuersig","issuerPublicKey": "issuerpubkey"}]`,
		`}`,
		`}`,
	}, "")

	jsonStr := strings.TrimPrefix(contractOutput, tokenDataPrefix)
	jsonStr = fixIssuerAddresses(jsonStr)
	jsonStr = fixFreeTextFields(jsonStr)

	var payload entity.CertificatePayload
	err := json.Unmarshal([]byte(jsonStr), &payload)
	require.NoError(t, err, "payload with embedded quote in eventName must parse without error")
	assert.Equal(t, `Hello "everyone`, payload.Data.EventName)
	assert.Equal(t, "VALID", payload.Data.Status)
}

func TestGetTokenData_FixIssuerAddresses_ParsesFullPayload(t *testing.T) {
	// Simulates what the contract returns for a minted certificate with one issuer
	// address — the payload that caused the 500 before the fix.
	contractOutput := strings.Join([]string{
		tokenDataPrefix,
		`{`,
		`"header":{`,
		`"@context":["https://www.w3.org/2018/credentials/v1"],`,
		`"type":["VerifiableCredential","EventCertificate"],`,
		`"id":"cert-1",`,
		`"issuer":"0xissuer",`,
		`"issuanceDate":"2024-01-01T00:00:00Z"`,
		`},`,
		`"data":{`,
		`"eventName":"Test Event",`,
		`"eventDescription":"desc",`,
		`"certificateTokenId":"1",`,
		`"certificateId":"cert-1",`,
		`"userId":"user-1",`,
		`"issuerId":"issuer-1",`,
		`"issuedAt":"2024-01-01T00:00:00Z",`,
		`"issuerAddresses": "["0xabc123def456abc123def456abc123def456abc1"]",`,
		`"receiverAddress":"0xreceiver",`,
		`"encryptedUserData":"encdata",`,
		`"backendEncryptedUserData":"backendencdata",`,
		`"certificateTitle":"Certificate of Completion",`,
		`"certificateSubtitle":"subtitle",`,
		`"status":"VALID"`,
		`},`,
		`"proof":{`,
		`"encryptedByUserRawData":"userraw",`,
		`"encryptedByBackendRawData":"backendraw",`,
		`"hash":"0xhash",`,
		`"signMessage":"0xwallet,0xcontract,12345",`,
		`"host":{"signature":"0xhostsig","publicKey":"hostpubkey"},`,
		`"issuers":[{"issuerSignature":"0xissuersig","issuerPublicKey":"issuerpubkey"}]`,
		`}`,
		`}`,
	}, "")

	jsonStr := strings.TrimPrefix(contractOutput, tokenDataPrefix)
	jsonStr = fixIssuerAddresses(jsonStr)

	var payload entity.CertificatePayload
	err := json.Unmarshal([]byte(jsonStr), &payload)
	require.NoError(t, err, "fixed JSON must parse without error")

	assert.Equal(t, "Test Event", payload.Data.EventName)
	assert.Equal(t, "Certificate of Completion", payload.Data.CertificateTitle)
	assert.Equal(t, "0xabc123def456abc123def456abc123def456abc1", payload.Data.IssuerAddresses)
	assert.Equal(t, "VALID", payload.Data.Status)
	assert.Equal(t, "0xhostsig", payload.Proof.Host.Signature)
	require.Len(t, payload.Proof.Issuers, 1)
	assert.Equal(t, "0xissuersig", payload.Proof.Issuers[0].IssuerSignature)
}
