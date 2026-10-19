package file

import (
	"errors"
	"fmt"
	"path"
	"slices"
	"strings"
	"uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// TODO: rewrite func using io.Copy
// TODO: write func to upload multiple files into dir
// UploadAndSaveFile uploads file from request, saves it
// to disk and returns filename (without dirpath).
// Returns empty string if file hasn't been uploaded.
// AcceptExtensions param is a slice with accepted extensions (e.g. []string{"png", "jpeg"}).
// If acceptExtensions is nil then all extensions are allowed.
func UploadAndSaveFile(ctx *fiber.Ctx, dirpath string, key string,
	acceptExtensions []string) (string, error) {

	// upload file from request
	file, err := ctx.FormFile(key)
	if errors.Is(err, fasthttp.ErrMissingFile) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("upload file from request: %w", err)
	}
	// check file extension
	filenameList := strings.Split(file.Filename, ".")
	extension := filenameList[len(filenameList)-1]
	if acceptExtensions != nil && !slices.Contains(acceptExtensions, extension) {
		return "", fmt.Errorf("invalid file extension: expected one of %v", acceptExtensions)
	}

	// collect path to file
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), extension)
	filepath := path.Join(dirpath, filename)

	// save file to disk
	if err := ctx.SaveFile(file, filepath); err != nil {
		return "", fmt.Errorf("save file to disk: %w", err)
	}
	return filename, nil
}
