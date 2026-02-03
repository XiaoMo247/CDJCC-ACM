package database

import (
	"fmt"
)

type folderIndexRow struct {
	IndexName string `gorm:"column:index_name"`
}

type folderCountRow struct {
	Cnt int `gorm:"column:cnt"`
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

	// 添加自引用外键（若已存在则跳过，避免启动时重复报错刷屏）。
	var fkCnt folderCountRow
	if err := DB.Raw(`
SELECT COUNT(*) AS cnt
FROM information_schema.REFERENTIAL_CONSTRAINTS
WHERE CONSTRAINT_SCHEMA = DATABASE()
  AND TABLE_NAME = 'folders'
  AND CONSTRAINT_NAME = 'fk_folders_parent'
`).Scan(&fkCnt).Error; err == nil && fkCnt.Cnt == 0 {
		_ = DB.Exec(`
ALTER TABLE folders
  ADD CONSTRAINT fk_folders_parent
  FOREIGN KEY (parent_id) REFERENCES folders(id)
  ON DELETE CASCADE
`).Error
	}

	// 确保 parent_id 索引存在（若已存在则跳过）。
	var idxCnt folderCountRow
	if err := DB.Raw(`
SELECT COUNT(*) AS cnt
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'folders'
  AND INDEX_NAME = 'idx_parent_id'
`).Scan(&idxCnt).Error; err == nil && idxCnt.Cnt == 0 {
		_ = DB.Exec(`CREATE INDEX idx_parent_id ON folders(parent_id)`).Error
	}
}

