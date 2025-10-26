package common

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type ValidateEvmAddressErrCode string

var (
	ErrNotHexString    ValidateEvmAddressErrCode = "NOT_HEX_STRING"
	ErrInvalidChecksum ValidateEvmAddressErrCode = "INVALID_CHECKSUM"
)

func ValidateEvmAddress(address string) (bool, *ValidateEvmAddressErrCode) {
	if !common.IsHexAddress(address) {
		return false, &ErrNotHexString
	}

	addr := common.HexToAddress(address)
	checksumAddr := addr.Hex()

	// Verify checksum if mixed case
	if address != strings.ToLower(address) && address != strings.ToUpper(address) {
		if address != checksumAddr {
			return false, &ErrInvalidChecksum
		}
		return true, nil
	}

	return true, nil
}

type SolutionStatus int32

const (
	SolutionStatusManaged SolutionStatus = 0
	SolutionStatusBYOK    SolutionStatus = 1
)
