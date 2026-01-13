package model

import "time"

type Folder struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Logical folder name (unique under the same parent).
	Name string `gorm:"type:varchar(100);not null;index:idx_folder_parent_name,unique" json:"name"`

	// Parent folder (NULL means root).
	ParentID *uint   `gorm:"index;index:idx_folder_parent_name,unique" json:"parent_id"`
	Parent   *Folder `gorm:"constraint:OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
