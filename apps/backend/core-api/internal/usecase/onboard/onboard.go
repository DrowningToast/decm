package usecase

import (
	"context"
	"errors"

	"apps/backend/common"
	customerror "apps/backend/common/customerror"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/core-api/internal/usecase/jwtutils"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type OnboardUsecase struct {
	registerSignMessage string

	AuthenticationCredentialDg datagateway.AuthenticationCredentialDataGateway
}

func NewOnboardUsecase(authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway) *OnboardUsecase {
	return &OnboardUsecase{
		registerSignMessage:        "Please sign this message to prove your ownership of the wallet",
		AuthenticationCredentialDg: authenticationCredentialDg,
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
		return nil, customerror.TryParseAsCustomErrWithMsg(&customerror.ErrInvalidArgument, errors.New("invalid wallet address"), errCode)
	}
	if !isValid {
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("invalid wallet address"))
	}

	// Validate the signed message
	address := ethcommon.HexToAddress(walletAddress)
	result, err := cyptoutils.VerifySignedMessageByAddress(address, u.registerSignMessage, signedMsg)
	if err != nil {
		return nil, customerror.TryParseAsCustomErrWithMsg(&customerror.ErrInternalServer, err, "an error has occured while verifying signed message")
	}
	if !result {
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("verification failed"))
	}

	// Check for already existing wallet address
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialByWalletAddress(context.Background(), walletAddress)
	if err == nil || credential != nil {
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrDuplicateEntry, errors.New("wallet address already exists"))
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
		return nil, customerror.TryParseAsCustomErrWithMsg(&customerror.ErrInternalServer, err, "failed to create new credential")
	}

	sessionToken, err := jwtutils.CreateToken(jwtutils.JwtPayload{
		UserId:        credential.Id,
		WalletAddress: credential.WalletAddress,
	})
	if err != nil {
		return nil, customerror.TryParseAsCustomErr(&customerror.ErrInternalServer, err)
	}

	return &sessionToken, nil
}
