package api

import (
	"acm-site/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminGetFolderTree(c *gin.Context) {
	tree, err := service.BuildFolderTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get folder tree"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tree": tree})
}

func AdminGetFolderContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	content, err := service.GetFolderContent(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, content)
}

func AdminGetFolderBreadcrumb(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	breadcrumb, err := service.BuildBreadcrumb(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"breadcrumb": breadcrumb})
}

func AdminCreateFolder(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		ParentID *uint  `json:"parentId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	folder, err := service.CreateFolder(req.Name, req.ParentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"folder": folder})
}

func AdminUploadFileToFolder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file is required"})
		return
	}

	newFile, err := service.SaveUploadedFile(file, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok", "file": newFile})
}

func AdminSearch(c *gin.Context) {
	q := c.Query("q")
	t := c.Query("type")
	limitRaw := c.Query("limit")
	limit := 50
	if limitRaw != "" {
		if parsed, err := strconv.Atoi(limitRaw); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	results, err := service.Search(q, t, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func AdminMoveItem(c *gin.Context) {
	var req struct {
		Type           string `json:"type" binding:"required"`
		ID             uint   `json:"id" binding:"required"`
		TargetFolderID *uint  `json:"targetFolderId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	switch req.Type {
	case "folder":
		if err := service.MoveFolder(req.ID, req.TargetFolderID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	case "file":
		if req.TargetFolderID == nil || *req.TargetFolderID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "targetFolderId is required for file move"})
			return
		}
		if err := service.MoveFile(req.ID, *req.TargetFolderID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

