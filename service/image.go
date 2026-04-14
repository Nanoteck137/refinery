package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/nanoteck137/refinery/assets"
	"github.com/nanoteck137/refinery/database"
	"github.com/nanoteck137/refinery/tools/utils"
	"github.com/nanoteck137/refinery/types"
)

var magickImageMapping = map[string]ImageType{
	"PNG":  ImageTypePng,
	"JPEG": ImageTypeJpeg,
}

type ImageService struct {
	logger *slog.Logger

	db      *database.Database
	dataDir types.DataDir
}

func NewImageService(
	logger *slog.Logger,
	db *database.Database,
	dataDir types.DataDir,
) *ImageService {
	return &ImageService{
		logger:  logger,
		db:      db,
		dataDir: dataDir,
	}
}

func (s *ImageService) convertImage(input, outputDir, name string, size int) (string, error) {
	p := path.Join(outputDir, name)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			err := utils.CreateResizedImage(input, p, size, size)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return p, nil
}

func (s *ImageService) convertSquareImage(input, outputDir, name string) (string, error) {
	p := path.Join(outputDir, name)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			err := utils.CreateSquareImage(input, p)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return p, nil
}

func (s *ImageService) copyDefaultToTemp(filename string) (string, error) {
	ext := path.Ext(filename)

	dest, err := os.CreateTemp("", "default*"+ext)
	if err != nil {
		return "", err
	}
	defer dest.Close()

	src, err := assets.DefaultImagesFS.Open(filename)
	if err != nil {
		return "", err
	}
	defer src.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return "", err
	}

	return dest.Name(), nil
}

// TODO(patrik): Rename to ImageFormat
// TODO(patrik): Move to types package
type ImageType string

const (
	ImageTypeEmpty   ImageType = ""
	ImageTypeUnknown ImageType = "unknown"
	ImageTypePng     ImageType = "png"
	ImageTypeJpeg    ImageType = "jpeg"
)

func (t ImageType) IsValid() bool {
	switch t {
	case ImageTypePng:
		return true
	case ImageTypeJpeg:
		return true
	}

	return false
}

func (t ImageType) ToExt() (string, bool) {
	switch t {
	case ImageTypePng:
		return ".png", true
	case ImageTypeJpeg:
		return ".jpeg", true
	}

	return "", false
}

func (s *ImageService) GetImageTypeFromExt(ext string) (ImageType, bool) {
	switch ext {
	case ".png":
		return ImageTypePng, true
	case ".jpg", ".jpeg":
		return ImageTypeJpeg, true
	}

	return "", false
}

func (s *ImageService) GetUserImage(ctx context.Context, userId, typ string, imageType ImageType) (string, error) {
	user, err := s.db.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return "", errors.New("user not found")
		}

		return "", err
	}

	cacheDir := s.dataDir.CacheImages()
	userCache := cacheDir.User(user.Id)

	// Make sure that the cache directory is setup
	dirs := []string{
		cacheDir.String(),
		cacheDir.Users(),
		userCache,
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}

	ext, ok := imageType.ToExt()
	if !ok {
		// TODO(patrik): Better error
		return "", errors.New("unknown image type")
	}

	input := ""

	if user.Picture.Valid {
		dir := s.dataDir.User(user.Id)
		input = path.Join(dir, user.Picture.String)
	} else {
		// TODO(patrik): Create a default user picture
		input, err = s.copyDefaultToTemp("default_album.png")
		if err != nil {
			return "", err
		}
		defer os.Remove(input)
	}

	switch typ {
	case "original":
		return s.convertSquareImage(input, userCache, "original_square"+ext)
	case "128":
		return s.convertImage(input, userCache, "128"+ext, 128)
	case "256":
		return s.convertImage(input, userCache, "256"+ext, 256)
	case "512":
		return s.convertImage(input, userCache, "512"+ext, 512)
	}

	return "", errors.New("unknown type")
}

func (s *ImageService) ValidateImage(p string) (ImageType, error) {
	cmd := exec.Command("magick", "identify", "-ping", "-format", "%m", p)

	var out bytes.Buffer
	cmd.Stdout = &out

	var errOut bytes.Buffer
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		details := strings.TrimSpace(errOut.String())

		s.logger.Error("failed to validate image",
			slog.Any("err", err),
			slog.String("output", details),
		)

		var execErr *exec.ExitError
		if errors.As(err, &execErr) {
			if !execErr.Success() {
				return ImageTypeUnknown, nil
			}
		}

		return ImageTypeUnknown, err
	}

	ty := strings.TrimSpace(out.String())

	res, exists := magickImageMapping[ty]
	if !exists {
		return ImageTypeUnknown, errors.New("no mapping: " + ty)
	}

	return res, nil
}

type DownloadPictureForUserParams struct {
	UserId string
	Url    string
}

// TODO(patrik): Cleanup
// TODO(patrik): Hash for files
func (s *ImageService) DownloadPictureForUser(
	ctx context.Context,
	params DownloadPictureForUserParams,
) (string, error) {
	// TODO(patrik): Cleanup, move to utils
	getImageExtFromContentType := func(contentType string) (string, error) {
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return "", fmt.Errorf("failed to parse content type: %w", err)
		}

		// TODO(patrik): Add support for more exts
		switch mediaType {
		case "image/png":
			return ".png", nil
		case "image/jpeg":
			return ".jpeg", nil
		default:
			return "", fmt.Errorf("unsupported media type: %s", mediaType)
		}
	}

	resp, err := http.Get(params.Url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	ext, err := getImageExtFromContentType(contentType)
	if err != nil {
		return "", err
	}

	// TODO(patrik): The tmp dir should be inside the work dir
	tmp, err := os.CreateTemp("", "tmp-image-*"+ext)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer tmp.Close()

	// always clean up temp file if something goes wrong
	defer func() {
		_, err := os.Stat(tmpPath)
		if err == nil {
			os.Remove(tmpPath)
		}
	}()

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		return "", err
	}

	tmp.Close()

	imageType, err := s.ValidateImage(tmpPath)
	if err != nil {
		return "", err
	}

	userDir := s.dataDir.User(params.UserId)

	err = utils.CreateDirectories([]string{
		userDir,
	})
	if err != nil {
		return "", err
	}

	imageExt, ok := imageType.ToExt()
	if !ok {
		return "", errors.New("invalid image type")
	}

	picture := "uploaded" + imageExt
	output := path.Join(userDir, picture)
	err = os.Rename(tmpPath, output)
	if err != nil {
		return "", fmt.Errorf("failed to promote temp file: %w", err)
	}

	return picture, nil
}
