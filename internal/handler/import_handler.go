package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"azh/internal/service"
	"github.com/julienschmidt/httprouter"
	"os"
)

// ImportHandler handles CSV import requests
type ImportHandler struct {
	importService *service.ImportService
}

// NewImportHandler creates a new ImportHandler
func NewImportHandler(importService *service.ImportService) *ImportHandler {
	return &ImportHandler{importService: importService}
}

// ImportCSV handles POST /api/import for uploading and processing CSV files
func (h *ImportHandler) ImportCSV(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse multipart form with a 32MB max memory limit
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Process uploaded files
	form := r.MultipartForm
	files := form.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var processedFiles []string
	for _, fileHeader := range files {
		if h.extractSingleCSVFile(w, fileHeader) {
			return
		}
		processedFiles = append(processedFiles, fileHeader.Filename)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "Successfully processed files: %v"}`, processedFiles)))
}

func (h *ImportHandler) extractSingleCSVFile(w http.ResponseWriter, fileHeader *multipart.FileHeader) bool {
	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to open file %s", fileHeader.Filename), http.StatusInternalServerError)
		return true
	}
	defer file.Close()

	// Create a temporary file using os.CreateTemp
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("csv-upload-%s-*", fileHeader.Filename))
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to create temporary file for %s", fileHeader.Filename), http.StatusInternalServerError)
		return true
	}
	tmpFilePath := tmpFile.Name()
	defer tmpFile.Close()
	defer os.Remove(tmpFilePath) // Clean up after processing

	// Write uploaded file content to temporary file
	if _, err := io.Copy(tmpFile, file); err != nil {
		http.Error(w, fmt.Sprintf("Unable to write to temporary file for %s", fileHeader.Filename), http.StatusInternalServerError)
		return true
	}

	// Process the CSV file
	err = h.importService.ProcessCSV(tmpFilePath, fileHeader.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing file %s: %v", fileHeader.Filename, err), http.StatusInternalServerError)
		return true
	}
	return false
}
