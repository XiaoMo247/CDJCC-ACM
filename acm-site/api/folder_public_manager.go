package api

import (
	"acm-site/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFolderTree(c *gin.Context) {
	tree, err := service.BuildFolderTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get folder tree"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tree": tree})
}

func GetFolderContent(c *gin.Context) {
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

func GetFolderBreadcrumb(c *gin.Context) {
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

func SearchPublic(c *gin.Context) {
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

