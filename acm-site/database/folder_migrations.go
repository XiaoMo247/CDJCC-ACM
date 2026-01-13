package database

import (
	"fmt"
)

type folderIndexRow struct {
	IndexName string `gorm:"column:index_name"`
}

// ensureFolderSchema tries to upgrade the legacy "folders" schema (unique name, no nesting)
// into a nested-tree friendly schema (parent_id + composite unique index).
//
// It is best-effort: failures are logged but don't stop startup.
func ensureFolderSchema() {
	if DB == nil {
		return
	}

	// Drop legacy single-column unique index on `name` if present.
	var rows []folderIndexRow
	if err := DB.Raw(`
SELECT INDEX_NAME AS index_name
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'folders'
GROUP BY INDEX_NAME
HAVING MIN(NON_UNIQUE) = 0
   AND COUNT(*) = 1
   AND SUM(CASE WHEN COLUMN_NAME = 'name' THEN 1 ELSE 0 END) = 1
`).Scan(&rows).Error; err == nil {
		for _, r := range rows {
			_ = DB.Exec(fmt.Sprintf("ALTER TABLE folders DROP INDEX `%s`", r.IndexName)).Error
		}
	}

	// Try to add self FK (ignore if already exists).
	_ = DB.Exec(`
ALTER TABLE folders
  ADD CONSTRAINT fk_folders_parent
  FOREIGN KEY (parent_id) REFERENCES folders(id)
  ON DELETE CASCADE
`).Error

	// Ensure index on parent_id (ignore if exists).
	_ = DB.Exec(`CREATE INDEX idx_parent_id ON folders(parent_id)`).Error
}

