package entity

import "time"

// File represents the database model for storing files and their content
type File struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Filename     string    `json:"filename" gorm:"not null"`
	OriginalName string    `json:"original_name" gorm:"not null"`
	FileSize     int64     `json:"file_size" gorm:"not null"`
	ContentType  string    `json:"content_type" gorm:"not null"`
	FileData     []byte    `json:"-" gorm:"type:bytea;not null"`           // Store actual file content
	ParsedText   *string   `json:"parsed_text,omitempty" gorm:"type:text"` // Nullable for parsing failures
	ParseError   *string   `json:"parse_error,omitempty" gorm:"type:text"` // Store parsing errors
	UploadedAt   time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for the File model
func (File) TableName() string {
	return "files"
}
