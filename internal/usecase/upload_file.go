package usecase

import (
	"leiia/internal/database/repository"
	"leiia/internal/entity"
)

type UploadFile struct {
	fileRepository repository.File
}

func (u *UploadFile) Execute(file *entity.File) error {
	// receive file from request
	// save file to database with status "received"
	// save file to object storage
	// read file
	// break file into chunks
	// update file status to "chunked"
	// return chunks

	return nil
}
