package usecase

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"

	ethutils "github.com/ethereum/go-ethereum/crypto"

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
	googleOAuthService  oauth_services.GoogleOAuthService
	registerSignMessage string

	AuthenticationCredentialDg datagateway.AuthenticationCredentialDataGateway
	ProfileDg                  datagateway.ProfileDataGateway
}

func NewOnboardUsecase(authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway, profileDg datagateway.ProfileDataGateway, authService *auth.AuthService) *OnboardUsecase {
	return &OnboardUsecase{
		authService:                authService,
		registerSignMessage:        "Please sign this message to prove your ownership of the wallet",
		AuthenticationCredentialDg: authenticationCredentialDg,
		ProfileDg:                  profileDg,
	}
}

func (u *OnboardUsecase) GetRegisterSignMessage() string {
	return u.registerSignMessage
}

// Register authentication credential with wallet address, and generate JWT token
func (u *OnboardUsecase) RegisterWithWalletAddress(ctx context.Context, signedMsg string, walletAddress string) (*string, *customerror.Err) {
	// Remove 0x prefix if it exists
	if len(walletAddress) >= 2 && walletAddress[:2] == "0x" {
		walletAddress = walletAddress[2:]
	}

	isValid, errCode := common.ValidateEvmAddress(walletAddress)
	if errCode != nil {
		errCode := string(*errCode)
		return nil, customerror.ParseWithMessage(&customerror.ErrInvalidArgument, errors.New("invalid wallet address"), errCode)
	}
	if !isValid {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid wallet address"))
	}

	// Validate the signed message
	address := ethcommon.HexToAddress(walletAddress)
	result, err := cyptoutils.VerifySignedMessageByAddress(address, u.registerSignMessage, signedMsg)
	if err != nil {
		return nil, customerror.ParseWithMessage(&customerror.ErrInternalServer, err, "an error has occured while verifying signed message")
	}
	if !result {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("verification failed"))
	}

	// Check for already existing wallet address
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialByWalletAddress(context.Background(), walletAddress)
	if err == nil || credential != nil {
		return nil, customerror.Parse(&customerror.ErrDuplicateEntry, errors.New("wallet address already exists"))
	}

	var cusErr customerror.Err
	if errors.As(err, &cusErr) {
		if cusErr.Code != &customerror.ErrNotFound.Code {
			return nil, cusErr.Extend("failed to check for existing wallet address")
		}
	}

	// Create new credential
	credential = &entity.AuthenticationCredential{
		WalletAddress:  walletAddress,
		SolutionStatus: entity.SolutionStatusBYOK,
	}

	credential, err = u.AuthenticationCredentialDg.CreateAuthenticationCredential(ctx, *credential)
	if err != nil {
		var cusErr customerror.Err
		if errors.As(err, &cusErr) {
			return nil, cusErr.Extend("failed to create new credential")
		}
		return nil, customerror.ParseWithMessage(&customerror.ErrInternalServer, err, "failed to create new credential")
	}

	sessionToken, err := u.authService.CreateToken(auth.JwtPayload{
		UserId:        credential.Id,
		WalletAddress: credential.WalletAddress,
	})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return &sessionToken, nil
}

// Register with Google OAuth token, return JWT token
func (u *OnboardUsecase) RegisterWithGoogle(ctx context.Context, token *oauth2.Token, password string) (*string, []string, *customerror.Err) {
	if password == "" {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("password is required"))
	}
	if len(password) < 6 {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("password must be at least 8 characters long"))
	}

	userInfo, customerr := u.googleOAuthService.GetUserInfo(ctx, token)
	if customerr != nil {
		return nil, nil, customerr.Extend("failed to get user info from google")
	}

	// Look for duplicate credentials
	// Must returns not found error
	credential, customerr := u.AuthenticationCredentialDg.GetAuthenticationCredentialByGoogleConnectorRef(ctx, userInfo.Id)
	var cusErr customerror.Err
	if !errors.As(customerr, &cusErr) {
		return nil, nil, cusErr.Extend("failed to check for existing google connector ref")
	}
	if cusErr.Code != &customerror.ErrNotFound.Code {
		return nil, nil, cusErr.Extend("failed to check for existing google connector ref")
	}
	if customerr == nil || credential != nil {
		return nil, nil, customerror.Parse(&customerror.ErrDuplicateEntry, errors.New("credential already exists"))
	}

	// Look for duplicate profile
	profile, customerr := u.ProfileDg.GetProfileByEmail(ctx, userInfo.Email)
	if !errors.As(customerr, &cusErr) {
		return nil, nil, cusErr.Extend("failed to check for existing profile")
	}
	if cusErr.Code != &customerror.ErrNotFound.Code {
		return nil, nil, cusErr.Extend("failed to check for existing profile")
	}
	if customerr == nil || profile != nil {
		return nil, nil, customerror.Parse(&customerror.ErrDuplicateEntry, errors.New("profile already exists"))
	}

	// Generate mnemonic
	wordCount := 12
	mnemonic, customerr := cyptoutils.GenerateMnemonic(&wordCount)
	if customerr != nil {
		return nil, nil, customerr.Extend("failed to generate mnemonic")
	}
	seed, customerr := cyptoutils.GenerateSeedFromMnemonic(mnemonic)
	if customerr != nil {
		return nil, nil, customerr.Extend("failed to generate seed")
	}

	// Generate wallet address
	privateKey, customerr := cyptoutils.GeneratePrivateKeyFromSeed(seed)
	if customerr != nil {
		return nil, nil, customerr.Extend("failed to generate private key")
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
		GoogleConnectorRef:  &userInfo.Id,
		SolutionStatus:      entity.SolutionStatusManaged,
		WalletAddress:       walletAddress.Hex(),
		EncryptedPrivateKey: &encryptedPrivateKey,
		HashedPassword:      &hashedPassword,
	}
	credential, customerr = u.AuthenticationCredentialDg.CreateAuthenticationCredential(ctx, *credential)
	if customerr != nil {
		return nil, nil, customerr.Extend("failed to create new credential")
	}

	sessionToken, err := u.authService.CreateToken(auth.JwtPayload{
		UserId:        credential.Id,
		WalletAddress: credential.WalletAddress,
	})
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return &sessionToken, strings.Fields(*mnemonic), nil
}
