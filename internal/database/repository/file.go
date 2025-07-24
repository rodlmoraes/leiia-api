package repository

import (
	"leiia/internal/entity"

	"gorm.io/gorm"
)

type File struct {
	db *gorm.DB
}

func (f *File) Create(file *entity.File) error {
	return f.db.Create(file).Error
}

func (f *File) Find(dest *entity.File, id string) error {
	return f.db.First(dest, "id = ?", id).Error
}
