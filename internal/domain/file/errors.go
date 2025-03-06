package file

import "errors"

var (
	ErrFileNotFound          = errors.New("file not found")
	ErrStorageQuotaExceeded  = errors.New("storage quota exceeded")
	ErrInvalidFile           = errors.New("invalid file")
	ErrAccessDenied          = errors.New("access denied")
	ErrStorageUploadFailed   = errors.New("storage upload failed")
	ErrStorageDownloadFailed = errors.New("storage download failed")
)

