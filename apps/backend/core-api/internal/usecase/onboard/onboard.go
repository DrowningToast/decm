package usecase

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/cockroachdb/errors"
	ethutils "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	"apps/backend/common"
	customerror "apps/backend/common/customerror"
	"apps/backend/common/encryptutils"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	oauth_services "apps/backend/services/oauth"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"golang.org/x/oauth2"
)

type OnboardUsecase struct {
	authService         *auth.AuthService
	googleOAuthService  *oauth_services.GoogleOAuthService
	registerSignMessage string

	AuthenticationCredentialDg datagateway.AuthenticationCredentialDataGateway
	ProfileDg                  datagateway.ProfileDataGateway
}

func NewOnboardUsecase(authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway, profileDg datagateway.ProfileDataGateway, authService *auth.AuthService, googleOAuthService *oauth_services.GoogleOAuthService) *OnboardUsecase {
	return &OnboardUsecase{
		authService:                authService,
		googleOAuthService:         googleOAuthService,
		registerSignMessage:        "Please sign this message to prove your ownership of the wallet",
		AuthenticationCredentialDg: authenticationCredentialDg,
		ProfileDg:                  profileDg,
	}
}

func (u *OnboardUsecase) GetRegisterSignMessage() string {
	return u.registerSignMessage
}

// Register authentication credential with wallet address, and generate JWT token
func (u *OnboardUsecase) RegisterWithWalletAddress(ctx context.Context, signedMsg string) (*uuid.UUID, *string, error) {
	walletAddress, err := cyptoutils.GetAddressFromSignature(u.registerSignMessage, signedMsg)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	if walletAddress == (ethcommon.Address{}) {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid signed message"))
	}

	isValid, errCode := common.ValidateEvmAddress(walletAddress.Hex())
	if errCode != nil {
		errCode := string(*errCode)
		return nil, nil, customerror.ParseWithMessage(&customerror.ErrInvalidArgument, errors.New("invalid wallet address"), errCode)
	}
	if !isValid {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid wallet address"))
	}

	// Validate the signed message
	address := ethcommon.HexToAddress(walletAddress.Hex())
	result, err := cyptoutils.VerifySignedMessageByAddress(address, u.registerSignMessage, signedMsg)
	if err != nil {
		return nil, nil, customerror.ParseWithMessage(&customerror.ErrInternalServer, err, "an error has occured while verifying signed message")
	}
	if !result {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("verification failed"))
	}

	// Check for already existing wallet address
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialByWalletAddress(ctx, walletAddress.Hex())
	if credential != nil {
		return nil, nil, customerror.Parse(&customerror.ErrDuplicateEntry, errors.New("wallet address already exists"))
	}
	if err == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to get authentication credential by wallet address")
	}
	var errNotFound *customerror.Err
	if !errors.As(err, &errNotFound) {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to get authentication credential by wallet address")
	}
	if *errNotFound.Code != customerror.ErrNotFound.Code {
		return nil, nil, errNotFound.Extend("failed to get authentication credential by wallet address")
	}

	// Create new credential
	credential = &entity.AuthenticationCredential{
		WalletAddress:  walletAddress.Hex(),
		SolutionStatus: common.SolutionStatusBYOK,
	}

	credential, err = u.AuthenticationCredentialDg.CreateAuthenticationCredential(ctx, *credential)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create new credential")
	}

	sessionToken, err := u.authService.CreateToken(auth.JwtPayload{
		UserId:        credential.Id,
		WalletAddress: credential.WalletAddress,
	})
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return &credential.Id, &sessionToken, nil
}

// Register with credential id, return JWT token
func (u *OnboardUsecase) RegisterWithGoogle(ctx context.Context, token *oauth2.Token, password string) (*uuid.UUID, *string, error) {
	if password == "" {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("password is required"))
	}
	if len(password) < 6 {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("password must be at least 8 characters long"))
	}

	userInfo, customerr := u.googleOAuthService.GetUserInfo(ctx, token)
	if customerr != nil {
		var customErr *customerror.Err
		if errors.As(customerr, &customErr) {
			return nil, nil, customErr.Extend("failed to get user info from google")
		}
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, customerr).Extend("failed to get user info from google")
	}

	// Look for duplicate credentials
	// Must returns not found error
	credential, customerr := u.AuthenticationCredentialDg.GetAuthenticationCredentialByGoogleConnectorRef(ctx, userInfo.Email)
	var cusErr *customerror.Err
	if !errors.As(customerr, &cusErr) {
		return nil, nil, cusErr.Extend("failed to check for existing google connector ref")
	}
	if *cusErr.Code != customerror.ErrNotFound.Code {
		return nil, nil, cusErr.Extend("failed to check for existing google connector ref")
	}
	if credential != nil {
		return nil, nil, customerror.Parse(&customerror.ErrDuplicateEntry, errors.New("credential already exists"))
	}

	// Look for duplicate profile
	profile, customerr := u.ProfileDg.GetProfileByEmail(ctx, userInfo.Email)
	if !errors.As(customerr, &cusErr) {
		return nil, nil, cusErr.Extend("failed to check for existing profile")
	}
	if *cusErr.Code != customerror.ErrNotFound.Code {
		return nil, nil, cusErr.Extend("failed to check for existing profile")
	}
	if customerr == nil || profile != nil {
		return nil, nil, customerror.Parse(&customerror.ErrDuplicateEntry, errors.New("profile already exists"))
	}

	// Generate mnemonic
	wordCount := 12
	mnemonic, customerr := cyptoutils.GenerateMnemonic(&wordCount)
	if customerr != nil {
		var customErr *customerror.Err
		if errors.As(customerr, &customErr) {
			return nil, nil, customErr.Extend("failed to generate mnemonic")
		}
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, customerr)
	}
	seed, customerr := cyptoutils.GenerateSeedFromMnemonic(mnemonic)
	if customerr != nil {
		var customErr *customerror.Err
		if errors.As(customerr, &customErr) {
			return nil, nil, customErr.Extend("failed to generate seed")
		}
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, customerr)
	}

	// Generate wallet address
	privateKey, customerr := cyptoutils.GeneratePrivateKeyFromSeed(seed)
	if customerr != nil {
		var customErr *customerror.Err
		if errors.As(customerr, &customErr) {
			return nil, nil, customErr.Extend("failed to generate private key")
		}
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, customerr)
	}
	privateKeyAsString := hex.EncodeToString(ethutils.FromECDSA(privateKey))

	// Encrypt private key
	encryptedPrivateKey, err := encryptutils.EncryptAESGCM(privateKeyAsString, password)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, customerr).Extend("failed to encrypt private key")
	}
	walletAddress, err := cyptoutils.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, customerr).Extend("failed to get wallet address")
	}

	hashedPassword, err := hashutils.HashPassword(password)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, customerr).Extend("failed to hash password")
	}

	// Create new credential
	credential = &entity.AuthenticationCredential{
		GoogleConnectorRef:  &userInfo.Email,
		SolutionStatus:      common.SolutionStatusManaged,
		WalletAddress:       walletAddress.Hex(),
		EncryptedPrivateKey: &encryptedPrivateKey,
		HashedPassword:      &hashedPassword,
	}
	credential, customerr = u.AuthenticationCredentialDg.CreateAuthenticationCredential(ctx, *credential)
	if customerr != nil {
		var customErr *customerror.Err
		if errors.As(customerr, &customErr) {
			return nil, nil, customErr.Extend("failed to create new credential")
		}
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, customerr)
	}

	sessionToken, err := u.authService.CreateToken(auth.JwtPayload{
		UserId:        credential.Id,
		WalletAddress: credential.WalletAddress,
	})
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return &credential.Id, &sessionToken, nil
}

func (u *OnboardUsecase) CheckOnboardStatusWithJwt(ctx context.Context, currentUser auth.JwtClaims) (*auth.JwtPayload, *uuid.UUID, error) {
	jwt := &auth.JwtPayload{
		UserId:        currentUser.UserId,
		WalletAddress: currentUser.WalletAddress,
	}

	profile, err := u.ProfileDg.GetProfileByAuthenticationCredentialId(ctx, currentUser.UserId)
	if err != nil {
		var notFoundErr *customerror.Err
		if errors.As(err, &notFoundErr) {
			return jwt, nil, nil
		}
		return jwt, nil, errors.Wrap(err, "failed to check for existing profile")
	}

	return jwt, &profile.Id, nil
}

// Return: JWT token, authentication credential id, profile id, error
func (u *OnboardUsecase) CheckOnboardStatusWithGoogleConnectorRef(ctx context.Context, accessToken string, expiresIn int) (*auth.JwtPayload, *uuid.UUID, error) {
	token := &oauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
	userInfo, err := u.googleOAuthService.GetUserInfo(ctx, token)
	if err != nil {
		var customErr *customerror.Err
		if errors.As(err, &customErr) {
			if customErr.Code != &customerror.ErrUnauthenticated.Code {
				return nil, nil, customErr.Extend("failed to get user info from google")
			}
			return nil, nil, customerror.Parse(&customerror.ErrUnauthenticated, err)
		}
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	// check for existing credential
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialByGoogleConnectorRef(ctx, userInfo.Email)
	if err != nil {
		var notFoundErr *customerror.Err
		if errors.As(err, &notFoundErr) {
			return nil, nil, nil
		}
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err).Extend("failed to check for existing google connector ref")
	}
	jwt := &auth.JwtPayload{
		UserId:        credential.Id,
		WalletAddress: credential.WalletAddress,
	}

	profile, err := u.ProfileDg.GetProfileByAuthenticationCredentialId(ctx, credential.Id)
	if err != nil {
		var notFoundErr *customerror.Err
		if errors.As(err, &notFoundErr) {
			return jwt, nil, nil
		}
		return jwt, nil, errors.Wrap(err, "failed to check for existing profile")
	}

	return jwt, &profile.Id, nil
}

// Return: is account exists, is profile exists, error
// Return: JWT token, authentication credential id, profile id, error
func (u *OnboardUsecase) CheckOnboardStatusWithWalletAddress(ctx context.Context, messageSignature string) (*auth.JwtPayload, *uuid.UUID, error) {
	walletAddress, err := cyptoutils.GetAddressFromSignature(u.registerSignMessage, messageSignature)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get wallet address from signature")
	}

	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialByWalletAddress(ctx, walletAddress.Hex())
	if err != nil {
		var notFoundErr *customerror.Err
		if !errors.As(err, &notFoundErr) {
			return nil, nil, errors.Wrap(err, "failed to check for existing wallet address")
		}
	}
	if credential == nil {
		return nil, nil, nil
	}

	profile, err := u.ProfileDg.GetProfileByAuthenticationCredentialId(ctx, credential.Id)
	if err != nil {
		var notFoundErr *customerror.Err
		if !errors.As(err, &notFoundErr) {
			return nil, nil, errors.Wrap(err, "failed to check for existing profile")
		}
	}
	jwt := &auth.JwtPayload{
		UserId:        credential.Id,
		WalletAddress: credential.WalletAddress,
	}
	if profile == nil {
		return jwt, nil, nil
	}
	return jwt, &profile.Id, nil
}
