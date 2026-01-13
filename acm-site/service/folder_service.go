package service

import (
	"acm-site/database"
	"acm-site/model"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func GetFoldersByParentID(parentID *uint) ([]model.Folder, error) {
	var folders []model.Folder
	query := database.DB.Order("created_at desc")
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	if err := query.Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func GetAllFolders() ([]model.Folder, error) {
	var folders []model.Folder
	if err := database.DB.Order("created_at desc").Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func GetFolderByID(id uint) (*model.Folder, error) {
	var folder model.Folder
	if err := database.DB.First(&folder, id).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func CreateFolder(name string, parentID *uint) (*model.Folder, error) {
	if name == "" {
		return nil, errors.New("folder name is required")
	}

	// Enforce unique folder name under the same parent at application level.
	var count int64
	q := database.DB.Model(&model.Folder{}).Where("name = ?", name)
	if parentID == nil {
		q = q.Where("parent_id IS NULL")
	} else {
		q = q.Where("parent_id = ?", *parentID)
	}
	if err := q.Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("folder already exists")
	}

	folder := model.Folder{
		Name:     name,
		ParentID: parentID,
	}
	if err := database.DB.Create(&folder).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func RemoveFolder(folderID uint) error {
	// Collect all descendants (including self).
	folderIDs, foldersByID, err := collectFolderDescendants(folderID)
	if err != nil {
		return err
	}

	// Delete file records and physical files.
	var files []model.File
	if err := database.DB.Where("folder_id IN ?", folderIDs).Find(&files).Error; err != nil {
		return fmt.Errorf("query files: %w", err)
	}
	for _, f := range files {
		_ = os.Remove(f.Path)
	}
	if err := database.DB.Where("folder_id IN ?", folderIDs).Delete(&model.File{}).Error; err != nil {
		return fmt.Errorf("delete file records: %w", err)
	}

	// Delete folder records.
	if err := database.DB.Where("id IN ?", folderIDs).Delete(&model.Folder{}).Error; err != nil {
		return fmt.Errorf("delete folder records: %w", err)
	}

	// Delete physical directories.
	for _, id := range folderIDs {
		_ = os.RemoveAll(filepath.Join("..", "uploads", "folders", fmt.Sprintf("%d", id)))

		// Legacy directory: ../uploads/<folder.Name> (best-effort; only for root folders).
		if folder, ok := foldersByID[id]; ok && folder.ParentID == nil {
			_ = os.RemoveAll(filepath.Join("..", "uploads", folder.Name))
		}
	}

	return nil
}

func collectFolderDescendants(folderID uint) ([]uint, map[uint]model.Folder, error) {
	var all []model.Folder
	if err := database.DB.Find(&all).Error; err != nil {
		return nil, nil, err
	}

	foldersByID := make(map[uint]model.Folder, len(all))
	childrenByParent := make(map[uint][]uint, len(all))

	exists := false
	for _, f := range all {
		foldersByID[f.ID] = f
		if f.ID == folderID {
			exists = true
		}
		if f.ParentID != nil {
			childrenByParent[*f.ParentID] = append(childrenByParent[*f.ParentID], f.ID)
		}
	}
	if !exists {
		return nil, nil, gorm.ErrRecordNotFound
	}

	var ids []uint
	queue := []uint{folderID}
	seen := map[uint]bool{folderID: true}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		ids = append(ids, cur)
		for _, child := range childrenByParent[cur] {
			if seen[child] {
				continue
			}
			seen[child] = true
			queue = append(queue, child)
		}
	}

	return ids, foldersByID, nil
}
