package pkg

import (
	"mime/multipart"
)

func IsImage(file *multipart.FileHeader) bool {
	allowedTypes := []string{"image/jpg", "image/png", "image/gif"}
	fileType := file.Header.Get("Content-Type")

	for _, allowedType := range allowedTypes {
		if fileType == allowedType {
			return true
		}
	}

	return false
}
