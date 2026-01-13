package api

import (
	"acm-site/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListFolders lists folders under a parent.
// If parent_id is omitted, it returns root folders.
func ListFolders(c *gin.Context) {
	parentIDRaw := c.Query("parent_id")
	var parentID *uint
	if parentIDRaw != "" {
		parsed, err := strconv.ParseUint(parentIDRaw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid parent_id"})
			return
		}
		tmp := uint(parsed)
		parentID = &tmp
	}

	folders, err := service.GetFoldersByParentID(parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to list folders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"folders": folders})
}

// AddFolder is a backward-compatible endpoint for creating folders.
// It supports optional `parentId` to enable nested folders.
func AddFolder(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "ok", "folder": folder})
}

func DeleteFolder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	if err := service.RemoveFolder(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

