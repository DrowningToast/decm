package auth

import (
	"apps/backend/common"
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"errors"

	gocommon "github.com/ethereum/go-ethereum/common"
)

func (u *AuthUsecase) SecuredSignStringForManagedUser(ctx context.Context, auth *auth.JwtClaims, message string, password string, deadlineBlock *uint64) ([]byte, *gocommon.Hash, error) {
	if auth == nil {
		return nil, nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}
	if auth.UserId == [16]byte{} {
		return nil, nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("invalid user id"))
	}

	// check type of user
	user, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, auth.UserId)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if user == nil {
		return nil, nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user not found"))
	}
	if user.SolutionStatus != common.SolutionStatusManaged {
		return nil, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("user is not a managed user"))
	}

	// check password
	match, err := hashutils.CompareHash(password, *user.HashedPassword)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if !match {
		return nil, nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("invalid password"))
	}

	// get private key
	privateKey, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, auth.UserId)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if privateKey == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("private key not found"))
	}

	// decrypt private key
	decryptedPk, _, err := cyptoutils.DecryptPrivateKey(*privateKey.EncryptedPrivateKey, password)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	messageHash := cyptoutils.HashEthereumMessage(message)
	signature, err := cyptoutils.Sign(messageHash.Bytes(), decryptedPk)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return signature, &messageHash, nil
}

func (u *AuthUsecase) SecuredSignActionForManagedUser(ctx context.Context, auth *auth.JwtClaims, password string, hostAddress gocommon.Address, contractAddress gocommon.Address, deadlineBlock *uint64) ([]byte, *gocommon.Hash, *uint64, error) {
	if auth == nil {
		return nil, nil, nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}
	if auth.UserId == [16]byte{} {
		return nil, nil, nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("invalid user id"))
	}

	if deadlineBlock == nil {
		calculatedDeadlineBlock, err := u.BlockchainClientDg.GetCalculatedDeadlineBlock(ctx)
		if err != nil {
			return nil, nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		deadlineBlock = &calculatedDeadlineBlock
	}

	signMessage, err := cyptoutils.GetSignMessage(hostAddress, contractAddress, *deadlineBlock)
	if err != nil {
		return nil, nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	encryptedPrivateKey, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, auth.UserId)
	if err != nil {
		return nil, nil, nil, err
	}
	if encryptedPrivateKey == nil {
		return nil, nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("encrypted private key not found"))
	}

	privateKey, _, err := cyptoutils.DecryptPrivateKey(*encryptedPrivateKey.EncryptedPrivateKey, password)
	if err != nil {
		return nil, nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	signature, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return signature, &messageHash, deadlineBlock, nil
}
