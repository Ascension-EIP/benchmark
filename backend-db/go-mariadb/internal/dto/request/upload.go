package request

import "../../../../../go-mariadb-benchmark/internal/dto/request/mime/multipart"

type Upload struct {
	File *multipart.FileHeader `form:"file"`
}
