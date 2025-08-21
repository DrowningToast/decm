package usecase

import (
	"errors"
	"regexp"

	"apps/backend/common"
	customerror "apps/backend/common/CustomError"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/usecase/cyptoutils"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type OnboardUsecase struct {
	RegisterSignMessage string

	AuthenticationCredentialDg *datagateway.AuthenticationCredentialDataGateway
}

func NewOnboardUsecase(registerSignMessage string, authenticationCredentialDg *datagateway.AuthenticationCredentialDataGateway) *OnboardUsecase {
	return &OnboardUsecase{
		RegisterSignMessage:        registerSignMessage,
		AuthenticationCredentialDg: authenticationCredentialDg,
	}
}

var ethAddressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

func (u *OnboardUsecase) GetRegisterSignMessage() string {
	return u.RegisterSignMessage
}

func (u *OnboardUsecase) RegisterWithPublicKey(signedMsg string, walletAddress string) (bool, *customerror.Err) {
	// Remove 0x prefix if it exists
	if len(walletAddress) >= 2 && walletAddress[:2] == "0x" {
		walletAddress = walletAddress[2:]
	}

	isValid, errCode := common.ValidateEvmAddress(walletAddress)
	if errCode != nil {
		errCode := string(*errCode)
		return false, customerror.TryParseAsCustomErrWithMsg(&customerror.ErrInvalidArgument, errors.New("invalid public key"), errCode)
	}
	if !isValid {
		return false, customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("invalid public key"))
	}

	// Validate the signed message
	address := ethcommon.HexToAddress(walletAddress)
	result, err := cyptoutils.VerifySignedMessageByAddress(address, u.RegisterSignMessage, signedMsg)
	if err != nil {
		return false, customerror.TryParseAsCustomErrWithMsg(&customerror.ErrInternalServer, err, "an error has occured while verifying signed message")
	}
	if !result {
		return false, customerror.TryParseAsCustomErr(&customerror.ErrInvalidArgument, errors.New("verification failed"))
	}
}
