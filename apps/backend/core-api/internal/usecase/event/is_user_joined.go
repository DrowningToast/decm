package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserJoinStatus struct {
	IsBroadcasted      bool       `json:"is_broadcasted"`
	BroadcastedAt      *time.Time `json:"broadcasted_at,omitempty"`
	Signature          string     `json:"signature"`
	SignMessage        string     `json:"sign_message"`
	JoinedAt           time.Time  `json:"joined_at"`
	IsAccepted         bool       `json:"is_accepted"`
	DeadlineBlock      *int32     `json:"deadline_block,omitempty"`
	EstimatedDeadline  *time.Time `json:"estimated_deadline,omitempty"`
	SignatureExpiredAt *time.Time `json:"signature_expired_at,omitempty"`
}

func (uc *EventUsecase) IsUserJoined(ctx context.Context, eventID uuid.UUID, credentialID uuid.UUID) (*UserJoinStatus, error) {
	// 1. Checks the row in event_attendee and signature details in one query
	attendee, err := uc.EventAttendeeDg.GetEventAttendeeWithSignature(ctx, eventID, credentialID)
	if err != nil {
		return nil, err
	}

	// Returns nil if not found (not joined)
	if attendee == nil {
		return nil, nil
	}

	status := &UserJoinStatus{
		JoinedAt:   attendee.JoinedAt,
		IsAccepted: attendee.IsAccepted,
	}

	if attendee.UserSignature != nil {
		status.IsBroadcasted = attendee.UserSignature.BroadcastedAt != nil
		status.BroadcastedAt = attendee.UserSignature.BroadcastedAt
		status.Signature = attendee.UserSignature.Signature
		status.SignMessage = attendee.UserSignature.SignMessage
		status.DeadlineBlock = attendee.UserSignature.DeadlineBlock
		status.EstimatedDeadline = attendee.UserSignature.EstimatedDeadline
		status.SignatureExpiredAt = attendee.UserSignature.MarkAsExpiredAt
	} else {
		// No signature record means the broadcast has not occurred yet
		status.IsBroadcasted = false
		status.BroadcastedAt = nil
	}

	return status, nil
}
