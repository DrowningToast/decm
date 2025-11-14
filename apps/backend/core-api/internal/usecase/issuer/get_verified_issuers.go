package issuer

import (
	"context"
	"strings"

	"apps/backend/core-api/internal/entity"
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

		fetchLimit := 1000

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
			CreatedAt:                   profile.CreatedAt,
			UpdatedAt:                   profile.UpdatedAt,
		}
	}

	return viewModels, nil
}
