import {
    CommonSolutionStatus,
    type EntityProfile,
    type ProfileGetMyProfileViewModel,
} from "@decm/api";
import type { Profile, ProfileWithAuth, SolutionStatus } from "./AuthService";

export const mapSolutionStatus = (solutionStatus: CommonSolutionStatus): SolutionStatus => {
    if (solutionStatus === CommonSolutionStatus.SolutionStatusBYOK) {
        return "BYOK";
    }
    return "SYSTEM_MANAGED";
};

export const mapProfileViewModel = (profile: EntityProfile): Profile => {
    return {
        id: profile.id,
        authenticationCredentialId: profile.authentication_credential_id,
        isProfilePicturePublic: profile.is_profile_picture_public,
        academicEmail: profile.academic_email,
        academicInstitution: profile.academic_institution,
        address: profile.address,
        bio: profile.bio,
        email: profile.email,
        firstName: profile.first_name,
        lastName: profile.last_name,
        phoneNumber: profile.phone_number,
        profilePictureUrl: profile.profile_picture_url,
        isAcademicEmailPublic: profile.is_academic_email_public,
        isAcademicInstitutionPublic: profile.is_academic_institution_public,
        isAddressPublic: profile.is_address_public,
        isBioPublic: profile.is_bio_public,
        isEmailPublic: profile.is_email_public,
        isFirstNamePublic: profile.is_first_name_public,
        isLastNamePublic: profile.is_last_name_public,
        isPhoneNumberPublic: profile.is_phone_number_public,
    };
};

export const mapProfileWithAuthViewModel = (
    profile: ProfileGetMyProfileViewModel,
): ProfileWithAuth => {
    return {
        id: profile.profile_id,
        authenticationCredentialId: profile.authentication_credential_id,
        isProfilePicturePublic: profile.is_profile_picture_public,
        academicEmail: profile.academic_email,
        academicInstitution: profile.academic_institution,
        address: profile.address,
        bio: profile.bio,
        email: profile.email,
        firstName: profile.first_name,
        lastName: profile.last_name,
        phoneNumber: profile.phone_number,
        profilePictureUrl: profile.profile_picture_url,
        isAcademicEmailPublic: profile.is_academic_email_public,
        isAcademicInstitutionPublic: profile.is_academic_institution_public,
        isAddressPublic: profile.is_address_public,
        isBioPublic: profile.is_bio_public,
        isEmailPublic: profile.is_email_public,
        isFirstNamePublic: profile.is_first_name_public,
        isLastNamePublic: profile.is_last_name_public,
        isPhoneNumberPublic: profile.is_phone_number_public,

        walletAddress: profile.wallet_address,
        solutionStatus: mapSolutionStatus(profile.solution_status),
        googleConnectorRef: profile.google_connector_ref,
        githubConnectorRef: profile.github_connector_ref,

        profileId: profile.profile_id,
    };
};
