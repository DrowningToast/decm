package issuer

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"strings"
)

// GetVerifiedIssuersRequest contains parameters for fetching verified issuers
type GetVerifiedIssuersRequest struct {
	SearchQuery string
	Limit       int
	Offset      int
}

func (u *IssuerUsecase) GetVerifiedIssuers(ctx context.Context, req GetVerifiedIssuersRequest) ([]VerifiedIssuerViewModel, error) {
	var profiles []entity.Profile
	var err error

	fetchLimit := 1000

	// Fetch all issuer credentials to get Google connector refs and wallet addresses
	allCredentials, err := u.IssuerDg.ListAllIssuerCredentials(ctx, fetchLimit)
	if err != nil {
		return nil, err
	}

	// Create maps of credential ID -> Google connector ref and wallet address
	googleEmailByCredId := make(map[string]*string)
	walletAddressByCredId := make(map[string]string)
	for _, cred := range allCredentials {
		credId := cred.Id.String()
		googleEmailByCredId[credId] = cred.GoogleConnectorRef
		walletAddressByCredId[credId] = cred.WalletAddress
	}

	// If no search query, return paginated results from simple list
	if req.SearchQuery == "" {
		profiles, err = u.IssuerDg.ListVerifiedIssuerProfiles(ctx, req.Limit, req.Offset)
		if err != nil {
			return nil, err
		}
	} else {
		// For search, use two-pronged approach:
		// 1. Search wallet_address at DB level (not encrypted)
		// 2. Fetch all profiles and filter encrypted fields in-memory

		// Search by wallet address at DB level
		credentialsByWallet, err := u.IssuerDg.SearchIssuerCredentialsByWalletAddress(ctx, req.SearchQuery, fetchLimit, 0)
		if err != nil {
			return nil, err
		}

		// Get credential IDs that matched wallet address search
		walletMatchedCredIds := make(map[string]bool)
		for _, cred := range credentialsByWallet {
			walletMatchedCredIds[cred.Id.String()] = true
		}

		// Fetch all issuer profiles for in-memory filtering of encrypted fields
		allProfiles, err := u.IssuerDg.ListIssuerProfiles(ctx, fetchLimit, 0)
		if err != nil {
			return nil, err
		}

		// Filter decrypted profiles based on search query (case-insensitive)
		searchLower := strings.ToLower(req.SearchQuery)
		var filteredProfiles []entity.Profile
		seenIds := make(map[string]bool) // Prevent duplicates

		for _, profile := range allProfiles {
			credId := profile.AuthenticationCredentialId.String()

			// Skip if already added
			if seenIds[credId] {
				continue
			}

			matched := false

			// Check if matched by wallet address search
			if walletMatchedCredIds[credId] {
				matched = true
			}

			// Search in first name (encrypted - searched after pgmapper decryption)
			if !matched && profile.FirstName != nil && strings.Contains(strings.ToLower(*profile.FirstName), searchLower) {
				matched = true
			}

			// Search in last name (encrypted)
			if !matched && profile.LastName != nil && strings.Contains(strings.ToLower(*profile.LastName), searchLower) {
				matched = true
			}

			// Search in email (encrypted)
			if !matched && profile.Email != nil && strings.Contains(strings.ToLower(*profile.Email), searchLower) {
				matched = true
			}

			// Search in academic email (encrypted)
			if !matched && profile.AcademicEmail != nil && strings.Contains(strings.ToLower(*profile.AcademicEmail), searchLower) {
				matched = true
			}

			// Search in Google email (encrypted)
			if !matched {
				googleEmail := googleEmailByCredId[credId]
				if googleEmail != nil && strings.Contains(strings.ToLower(*googleEmail), searchLower) {
					matched = true
				}
			}

			if matched {
				filteredProfiles = append(filteredProfiles, profile)
				seenIds[credId] = true
			}
		}

		// Apply pagination to filtered results
		start := req.Offset
		end := req.Offset + req.Limit
		if start >= len(filteredProfiles) {
			profiles = []entity.Profile{}
		} else {
			if end > len(filteredProfiles) {
				end = len(filteredProfiles)
			}
			profiles = filteredProfiles[start:end]
		}
	}

	// Convert entities to ViewModels
	viewModels := make([]VerifiedIssuerViewModel, len(profiles))
	for i, profile := range profiles {
		credId := profile.AuthenticationCredentialId.String()
		viewModels[i] = VerifiedIssuerViewModel{
			Id:                          profile.Id,
			AuthenticationCredentialId:  profile.AuthenticationCredentialId,
			IsProfilePicturePublic:      profile.IsProfilePicturePublic,
			ProfilePictureUrl:           profile.ProfilePictureUrl,
			IsFirstNamePublic:           profile.IsFirstNamePublic,
			FirstName:                   profile.FirstName,
			IsLastNamePublic:            profile.IsLastNamePublic,
			LastName:                    profile.LastName,
			IsEmailPublic:               profile.IsEmailPublic,
			Email:                       profile.Email,
			IsBioPublic:                 profile.IsBioPublic,
			Bio:                         profile.Bio,
			IsPhoneNumberPublic:         profile.IsPhoneNumberPublic,
			PhoneNumber:                 profile.PhoneNumber,
			IsAddressPublic:             profile.IsAddressPublic,
			Address:                     profile.Address,
			IsAcademicInstitutionPublic: profile.IsAcademicInstitutionPublic,
			AcademicInstitution:         profile.AcademicInstitution,
			IsAcademicEmailPublic:       profile.IsAcademicEmailPublic,
			AcademicEmail:               profile.AcademicEmail,
			GoogleConnectorRef:          googleEmailByCredId[credId],
			WalletAddress:               walletAddressByCredId[credId],
			CreatedAt:                   profile.CreatedAt,
			UpdatedAt:                   profile.UpdatedAt,
		}
	}

	return viewModels, nil
}
