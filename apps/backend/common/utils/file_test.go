package utils

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFileContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    string
	}{
		{
			name:        "JPEG",
			contentType: "image/jpeg",
			expected:    "image/jpeg",
		},
		{
			name:        "PNG",
			contentType: "image/png",
			expected:    "image/png",
		},
		{
			name:        "PDF",
			contentType: "application/pdf",
			expected:    "application/pdf",
		},
		{
			name:        "empty",
			contentType: "",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &multipart.FileHeader{
				Header: make(map[string][]string),
			}
			if tt.contentType != "" {
				file.Header.Set("Content-Type", tt.contentType)
			}

			result := GetFileContentType(file)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "JPEG file",
			filename: "image.jpg",
			expected: ".jpg",
		},
		{
			name:     "PNG file",
			filename: "image.png",
			expected: ".png",
		},
		{
			name:     "PDF file",
			filename: "document.pdf",
			expected: ".pdf",
		},
		{
			name:     "no extension",
			filename: "file",
			expected: "",
		},
		{
			name:     "multiple dots",
			filename: "file.backup.txt",
			expected: ".txt",
		},
		{
			name:     "hidden file",
			filename: ".hidden",
			expected: ".hidden",
		},
		{
			name:     "path with extension",
			filename: "/path/to/file.txt",
			expected: ".txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &multipart.FileHeader{
				Filename: tt.filename,
			}

			result := GetFileExtension(file)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOpenFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	content := []byte("test content")
	_, err = tmpFile.Write(content)
	require.NoError(t, err)
	_ = tmpFile.Close()

	// We need to create a proper multipart form to test OpenFile
	// Since multipart.FileHeader.Open() requires the file to be part of a form,
	// we'll test the basic functionality
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(tmpFile.Name()))
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	_ = writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(10 << 20) // 10MB max
	require.NoError(t, err)
	defer func() { _ = form.RemoveAll() }()

	files := form.File["file"]
	require.Len(t, files, 1)
	fileHeader := files[0]

	file, err := OpenFile(fileHeader)
	require.NoError(t, err)
	defer func() { _ = file.Close() }()

	// Read the file content
	buf := make([]byte, len(content))
	n, err := file.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, len(content), n)
	assert.Equal(t, content, buf[:n])
}

func TestGetFileContentType_CaseInsensitive(t *testing.T) {
	file := &multipart.FileHeader{
		Header: make(map[string][]string),
	}
	file.Header.Set("content-type", "image/jpeg") // lowercase

	result := GetFileContentType(file)
	// Note: Header.Get is case-insensitive, so this should work
	assert.Equal(t, "image/jpeg", result)
}

func TestGetFileExtension_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "empty filename",
			filename: "",
			expected: "",
		},
		{
			name:     "only extension",
			filename: ".txt",
			expected: ".txt",
		},
		{
			name:     "extension only with dot",
			filename: ".",
			expected: ".", // filepath.Ext(".") returns "."
		},
		{
			name:     "multiple extensions",
			filename: "file.tar.gz",
			expected: ".gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &multipart.FileHeader{
				Filename: tt.filename,
			}

			result := GetFileExtension(file)
			assert.Equal(t, tt.expected, result)
		})
	}
}
