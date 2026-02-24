package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	uploadDir  string
	maxSizeMB int
}

func NewUploadHandler(uploadDir string, maxSizeMB int) *UploadHandler {
	os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{uploadDir: uploadDir, maxSizeMB: maxSizeMB}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(h.maxSizeMB)<<20)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		Error(c, http.StatusBadRequest, "failed to read uploaded file")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		Error(c, http.StatusBadRequest, "unsupported file type")
		return
	}

	filename := uuid.New().String() + ext
	destPath := filepath.Join(h.uploadDir, filename)

	if err := c.SaveUploadedFile(header, destPath); err != nil {
		Error(c, http.StatusInternalServerError, "failed to save file")
		return
	}

	url := fmt.Sprintf("/uploads/%s", filename)
	Success(c, gin.H{"url": url})
}
