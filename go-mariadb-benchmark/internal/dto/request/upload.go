package request

import "mime/multipart"

type Upload struct {
	File *multipart.FileHeader `json:"file"`
}
