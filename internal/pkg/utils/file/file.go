package file

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// UploadAndSaveFile uploads file from request with given key, saves it
// to dir and returns filename (without dirpath).
// Returns empty string if file hasn't been uploaded.
func UploadAndSaveFile(ctx *fiber.Ctx, key, dir string) (string, error) {
	// get file from request
	file, err := ctx.FormFile(key)
	if errors.Is(err, fasthttp.ErrMissingFile) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get file from request: %w", err)
	}
	return saveFile(file, dir)
}

// UploadAndSaveFiles uploads files from request with given key, saves them
// to dir and returns filenames (without dirpath).
// Returns nil if files haven't been uploaded.
func UploadAndSaveFiles(ctx *fiber.Ctx, key, dir string) ([]string, error) {
	// get files from request
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("get mpfd: %w", err)
	}
	files, ok := form.File[key]
	if !ok || len(files) == 0 { // no one file was uploaded
		return nil, nil
	}

	filenames := make([]string, len(files))
	for idx, file := range files {
		filename, err := saveFile(file, dir)
		if err != nil {
			return nil, err
		}
		filenames[idx] = filename
	}
	return filenames, nil
}

// saveFile saves given file to the given dir and returns filename (without dirpath).
func saveFile(file *multipart.FileHeader, dir string) (string, error) {
	// get file extension
	filenameList := strings.Split(file.Filename, ".")
	extension := filenameList[len(filenameList)-1]
	// collect path to file
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), extension)
	filepath := path.Join(dir, filename)

	src, err := file.Open()
	if err != nil {
		return filename, fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return filename, fmt.Errorf("create dest file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return filename, fmt.Errorf("save file data: %w", err)
	}
	return filename, nil
}
