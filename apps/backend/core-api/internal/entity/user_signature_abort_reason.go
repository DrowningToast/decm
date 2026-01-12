package entity

// UserSignatureAbortReason is the typed reason code stored in user_signature.aborted_reason.
type UserSignatureAbortReason int32

const (
	// AbortReasonParticipantNotJoined is set when the worker detects that the
	// certificate claimant has never joined the event on-chain and there is no
	// pending join in the queue, so the claim can never succeed.
	AbortReasonParticipantNotJoined UserSignatureAbortReason = 1
)
