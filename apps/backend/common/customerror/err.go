package customerror

type Err struct {
	HttpStatus *int     `json:"http_status"`
	Code       *ErrCode `json:"code"`

	Message string
	Inner   error
}

func (e Err) Error() string {
	return e.Message + e.Inner.Error()
}

func TryParseAsCustomErr(preset *ErrSignature, err error) *Err {
	var defaultErrSignature *ErrSignature = &ErrInternalServer
	if preset != nil {
		defaultErrSignature = &ErrSignature{
			Code:           preset.Code,
			DefaultMessage: preset.DefaultMessage,
			HttpStatus:     preset.HttpStatus,
		}
	}

	return AsPresetError(*defaultErrSignature, err)
}

func TryParseAsCustomErrWithMsg(preset *ErrSignature, err error, message string) *Err {
	var defaultErrSignature *ErrSignature = &ErrInternalServer
	if preset != nil {
		defaultErrSignature = &ErrSignature{
			Code:           preset.Code,
			DefaultMessage: preset.DefaultMessage,
			HttpStatus:     preset.HttpStatus,
		}
	}

	return &Err{
		HttpStatus: &defaultErrSignature.HttpStatus,
		Code:       &defaultErrSignature.Code,
		Message:    message,
		Inner:      err,
	}
}

func AsPresetError(preset ErrSignature, err error) *Err {
	httpStatus := preset.HttpStatus
	code := preset.Code

	return &Err{
		HttpStatus: &httpStatus,
		Code:       &code,
		Message:    preset.DefaultMessage,
		Inner:      err,
	}
}
