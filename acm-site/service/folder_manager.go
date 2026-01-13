package service

import (
	"acm-site/database"
	"acm-site/model"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FolderTreeNode struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	ParentID *uint          `json:"parentId"`
	Children []FolderTreeNode `json:"children"`
}

type BreadcrumbItem struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FolderContentFolder struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	ParentID    *uint  `json:"parentId"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	FileCount   int64  `json:"fileCount"`
	FolderCount int64  `json:"folderCount"`
}

type FolderContentSubFolder struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	ParentID  *uint  `json:"parentId"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	FileCount int64  `json:"fileCount"`
}

type FolderContentFile struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	UploadedAt string `json:"uploadedAt"`
}

type FolderContent struct {
	Folder     FolderContentFolder      `json:"folder"`
	SubFolders []FolderContentSubFolder `json:"subFolders"`
	Files      []FolderContentFile      `json:"files"`
}

type SearchResult struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // "folder" | "file"
	Path string `json:"path"`
	Size int64  `json:"size,omitempty"`
	FolderID *uint `json:"folderId,omitempty"`
}

func BuildFolderTree() ([]FolderTreeNode, error) {
	var folders []model.Folder
	if err := database.DB.Order("created_at asc").Find(&folders).Error; err != nil {
		return nil, err
	}

	childrenByParent := make(map[uint][]model.Folder)
	for _, f := range folders {
		parentKey := uint(0)
		if f.ParentID != nil {
			parentKey = *f.ParentID
		}
		childrenByParent[parentKey] = append(childrenByParent[parentKey], f)
	}

	var build func(parentKey uint) []FolderTreeNode
	build = func(parentKey uint) []FolderTreeNode {
		children := childrenByParent[parentKey]
		nodes := make([]FolderTreeNode, 0, len(children))
		for _, child := range children {
			nodes = append(nodes, FolderTreeNode{
				ID:       child.ID,
				Name:     child.Name,
				ParentID: child.ParentID,
				Children: build(child.ID),
			})
		}
		return nodes
	}

	return build(0), nil
}

func BuildBreadcrumb(folderID uint) ([]BreadcrumbItem, error) {
	if folderID == 0 {
		return []BreadcrumbItem{{ID: 0, Name: "根目录"}}, nil
	}

	all, byID, err := loadFolderIndex()
	if err != nil {
		return nil, err
	}
	if _, ok := byID[folderID]; !ok {
		_ = all
		return nil, errors.New("folder not found")
	}

	var crumbs []BreadcrumbItem
	cur := folderID
	for cur != 0 {
		f := byID[cur]
		crumbs = append(crumbs, BreadcrumbItem{ID: f.ID, Name: f.Name})
		if f.ParentID == nil {
			break
		}
		cur = *f.ParentID
	}

	// Reverse, and prepend root.
	for i, j := 0, len(crumbs)-1; i < j; i, j = i+1, j-1 {
		crumbs[i], crumbs[j] = crumbs[j], crumbs[i]
	}

	return append([]BreadcrumbItem{{ID: 0, Name: "根目录"}}, crumbs...), nil
}

func BuildFolderPath(folderID uint) (string, []BreadcrumbItem, error) {
	crumbs, err := BuildBreadcrumb(folderID)
	if err != nil {
		return "", nil, err
	}
	if len(crumbs) == 1 {
		return "/", crumbs, nil
	}
	parts := make([]string, 0, len(crumbs)-1)
	for _, c := range crumbs[1:] {
		parts = append(parts, c.Name)
	}
	return "/" + strings.Join(parts, "/"), crumbs, nil
}

func GetFolderContent(folderID uint) (*FolderContent, error) {
	if folderID == 0 {
		subFolders, err := GetFoldersByParentID(nil)
		if err != nil {
			return nil, err
		}

		subFolderDTOs, err := toSubFolderDTOs(subFolders)
		if err != nil {
			return nil, err
		}

		return &FolderContent{
			Folder: FolderContentFolder{
				ID:          0,
				Name:        "根目录",
				ParentID:    nil,
				Path:        "/",
				Size:        0,
				FileCount:   0,
				FolderCount: int64(len(subFolders)),
			},
			SubFolders: subFolderDTOs,
			Files:      []FolderContentFile{},
		}, nil
	}

	folder, err := GetFolderByID(folderID)
	if err != nil {
		return nil, err
	}

	path, _, err := BuildFolderPath(folderID)
	if err != nil {
		return nil, err
	}

	subFolders, err := GetFoldersByParentID(&folderID)
	if err != nil {
		return nil, err
	}

	files, err := GetFilesByFolderID(folderID)
	if err != nil {
		return nil, err
	}

	fileDTOs := make([]FolderContentFile, 0, len(files))
	var totalSize int64
	for _, f := range files {
		totalSize += f.Size
		fileDTOs = append(fileDTOs, FolderContentFile{
			ID:         f.ID,
			Name:       f.Name,
			Size:       f.Size,
			UploadedAt: f.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	subFolderDTOs, err := toSubFolderDTOs(subFolders)
	if err != nil {
		return nil, err
	}

	var folderCount int64
	if err := database.DB.Model(&model.Folder{}).Where("parent_id = ?", folderID).Count(&folderCount).Error; err != nil {
		return nil, err
	}

	return &FolderContent{
		Folder: FolderContentFolder{
			ID:          folder.ID,
			Name:        folder.Name,
			ParentID:    folder.ParentID,
			Path:        path,
			Size:        totalSize,
			FileCount:   int64(len(files)),
			FolderCount: folderCount,
		},
		SubFolders: subFolderDTOs,
		Files:      fileDTOs,
	}, nil
}

func Search(q string, itemType string, limit int) ([]SearchResult, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return nil, errors.New("q is required")
	}
	if limit <= 0 {
		limit = 50
	}

	_, byID, err := loadFolderIndex()
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, 16)

	like := "%" + q + "%"
	if itemType == "" || itemType == "folder" {
		var folders []model.Folder
		if err := database.DB.Where("name LIKE ?", like).Limit(limit).Find(&folders).Error; err != nil {
			return nil, err
		}
		for _, f := range folders {
			p, _, perr := BuildFolderPath(f.ID)
			if perr != nil {
				p = "/"
			}
			results = append(results, SearchResult{
				ID:   f.ID,
				Name: f.Name,
				Type: "folder",
				Path: p,
			})
			if len(results) >= limit {
				return results, nil
			}
		}
	}

	if itemType == "" || itemType == "file" {
		var files []model.File
		if err := database.DB.Where("name LIKE ?", like).Limit(limit).Find(&files).Error; err != nil {
			return nil, err
		}
		for _, f := range files {
			folderPath := "/"
			if folder, ok := byID[f.FolderID]; ok {
				p, _, perr := BuildFolderPath(folder.ID)
				if perr == nil {
					folderPath = p
				}
			}
			results = append(results, SearchResult{
				ID:   f.ID,
				Name: f.Name,
				Type: "file",
				Path: strings.TrimRight(folderPath, "/") + "/" + f.Name,
				Size: f.Size,
				FolderID: func() *uint { v := f.FolderID; return &v }(),
			})
			if len(results) >= limit {
				return results, nil
			}
		}
	}

	return results, nil
}

func MoveFolder(folderID uint, targetParentID *uint) error {
	if targetParentID != nil && *targetParentID == folderID {
		return errors.New("cannot move folder into itself")
	}

	// Ensure folder exists.
	var folder model.Folder
	if err := database.DB.First(&folder, folderID).Error; err != nil {
		return err
	}

	// Validate target parent.
	if targetParentID != nil {
		var parent model.Folder
		if err := database.DB.First(&parent, *targetParentID).Error; err != nil {
			return err
		}
	}

	// Prevent cycles: target cannot be within folder's subtree.
	if targetParentID != nil {
		_, byID, err := loadFolderIndex()
		if err != nil {
			return err
		}

		cur := *targetParentID
		for cur != 0 {
			if cur == folderID {
				return errors.New("cannot move folder into its descendant")
			}
			f := byID[cur]
			if f.ParentID == nil {
				break
			}
			cur = *f.ParentID
		}
	}

	return database.DB.Model(&model.Folder{}).Where("id = ?", folderID).Update("parent_id", targetParentID).Error
}

func MoveFile(fileID uint, targetFolderID uint) error {
	// Ensure target folder exists.
	var targetFolder model.Folder
	if err := database.DB.First(&targetFolder, targetFolderID).Error; err != nil {
		return err
	}

	var file model.File
	if err := database.DB.First(&file, fileID).Error; err != nil {
		return err
	}

	targetDir := filepath.Join("..", "uploads", "folders", fmt.Sprintf("%d", targetFolder.ID))
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return err
	}

	oldPath := file.Path
	base := filepath.Base(oldPath)
	newPath := filepath.Join(targetDir, base)

	if oldPath != newPath {
		if err := os.Rename(oldPath, newPath); err != nil {
			// Fallback to copy+remove (e.g., across devices).
			if copyErr := copyFile(oldPath, newPath); copyErr != nil {
				return err
			}
			_ = os.Remove(oldPath)
		}
	}

	return database.DB.Model(&model.File{}).Where("id = ?", fileID).Updates(map[string]any{
		"folder_id": targetFolderID,
		"path":      newPath,
	}).Error
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func toSubFolderDTOs(subFolders []model.Folder) ([]FolderContentSubFolder, error) {
	if len(subFolders) == 0 {
		return []FolderContentSubFolder{}, nil
	}

	ids := make([]uint, 0, len(subFolders))
	for _, f := range subFolders {
		ids = append(ids, f.ID)
	}

	type aggRow struct {
		FolderID   uint  `gorm:"column:folder_id"`
		FileCount  int64 `gorm:"column:file_count"`
		TotalSize  int64 `gorm:"column:total_size"`
	}

	fileAgg := make(map[uint]aggRow, len(ids))
	var rows []aggRow
	if err := database.DB.Model(&model.File{}).
		Select("folder_id, COUNT(*) as file_count, COALESCE(SUM(size), 0) as total_size").
		Where("folder_id IN ?", ids).
		Group("folder_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		fileAgg[r.FolderID] = r
	}

	dtos := make([]FolderContentSubFolder, 0, len(subFolders))
	for _, f := range subFolders {
		p, _, perr := BuildFolderPath(f.ID)
		if perr != nil {
			p = "/"
		}
		agg := fileAgg[f.ID]
		dtos = append(dtos, FolderContentSubFolder{
			ID:        f.ID,
			Name:      f.Name,
			ParentID:  f.ParentID,
			Path:      p,
			Size:      agg.TotalSize,
			FileCount: agg.FileCount,
		})
	}

	return dtos, nil
}

func loadFolderIndex() ([]model.Folder, map[uint]model.Folder, error) {
	var all []model.Folder
	if err := database.DB.Find(&all).Error; err != nil {
		return nil, nil, err
	}
	byID := make(map[uint]model.Folder, len(all))
	for _, f := range all {
		byID[f.ID] = f
	}
	return all, byID, nil
}
