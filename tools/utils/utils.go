package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gosimple/slug"
	"github.com/nrednav/cuid2"
	"github.com/pelletier/go-toml/v2"
)

var CreateId = CreateIdGenerator(32)
var CreateSmallId = CreateIdGenerator(8)

var CreateUserId = CreateIdGenerator(10)

var CreateArtistId = CreateIdGenerator(10)
var CreateAlbumId = CreateIdGenerator(16)
var CreateTrackId = CreateIdGenerator(32)
var CreateTrackMediaId = CreateIdGenerator(32)

var CreatePlaylistFilterId = CreateIdGenerator(8)
var CreateTrackFilterId = CreateIdGenerator(8)

var CreateVirtualPlaylistId = CreateIdGenerator(16)

var CreateApiTokenId = CreateIdGenerator(32)

var CreateUserListeningEventId = CreateIdGenerator(32)

func CreateIdGenerator(length int) func() string {
	res, err := cuid2.Init(cuid2.WithLength(length))
	if err != nil {
		log.Fatal("Failed to create id generator", "err", err)
	}

	return res
}

func ParseAuthHeader(authHeader string) string {
	splits := strings.Split(authHeader, " ")
	if len(splits) != 2 {
		return ""
	}

	if splits[0] != "Bearer" {
		return ""
	}

	return splits[1]
}

// TODO(patrik): Move to ImageService
func CreateSquareImage(src, dest string) error {
	cmd := exec.Command(
		"magick", src,
		"-gravity", "Center",
		"-extent", "%[fx:min(w,h)]x%[fx:min(w,h)]",
		dest,
	)
	// TODO(patrik): Make this configureble
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Move to ImageService
func CreateResizedImage(src string, dest string, width, height int) error {
	args := []string{
		src,
		"-resize", fmt.Sprintf("%dx%d^", width, height),
		"-gravity", "Center",
		"-extent", fmt.Sprintf("%dx%d", width, height),
		dest,
	}

	cmd := exec.Command("magick", args...)
	// TODO(patrik): Make this configureble
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Move to ImageService
func GeneratePlaylistCover(images [4]string, output string, tileSize int) error {
	if len(images) == 0 {
		return fmt.Errorf("at least one image is required")
	}

	size := fmt.Sprintf("%dx%d", tileSize, tileSize)

	buildTile := func(img string) []string {
		if img == "" {
			return []string{"(", "xc:black", "-resize", size, ")"}
		}
		return []string{"(", img, "-resize", size + "^", "-gravity", "center", "-extent", size, ")"}
	}

	args := []string{}
	for _, img := range images {
		args = append(args, buildTile(img)...)
	}

	args = append(args,
		"(", "-clone", "0-1", "+append", ")",
		"(", "-clone", "2-3", "+append", ")",
		"-delete", "0-3",
		"-append",
		output,
	)

	cmd := exec.Command("magick", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// TODO(patrik): Move to ImageService
func ConvertImage(src string, dest string) error {
	args := []string{
		"convert",
		src,
		dest,
	}

	cmd := exec.Command("magick", args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func Slug(s string) string {
	return slug.Make(s)
}

func SplitString(s string) []string {
	tags := []string{}
	if s != "" {
		tags = strings.Split(s, ",")
	}

	return tags
}

func TotalPages(perPage, totalItems int) int {
	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}

func ExtractNumber(s string) int {
	n := ""
	for _, c := range s {
		if unicode.IsDigit(c) {
			n += string(c)
		} else {
			break
		}
	}

	if len(n) == 0 {
		return 0
	}

	i, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return 0
	}

	return int(i)
}

var validImageExts = []string{
	".png",
	".jpg",
	".jpeg",
}

func IsValidImageExt(ext string) bool {
	for _, e := range validImageExts {
		if ext == e {
			return true
		}
	}

	return false
}

// TODO(patrik): Update this
var validExts []string = []string{
	".wav",
	".flac",
	".opus",
}

func IsValidTrackExt(ext string) bool {
	for _, valid := range validExts {
		if valid == ext {
			return true
		}
	}

	return false
}

func SqlNullToStringPtr(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}

	return nil
}

func SqlNullToInt64Ptr(value sql.NullInt64) *int64 {
	if value.Valid {
		return &value.Int64
	}

	return nil
}

func SqlNullToFloat64Ptr(value sql.NullFloat64) *float64 {
	if value.Valid {
		return &value.Float64
	}

	return nil
}

// TODO(patrik): Move to auth service
const (
	letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
)

func randomString(charset string, length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

func GenerateCode() (string, error) {
	part1, err := randomString(letters, 4)
	if err != nil {
		return "", err
	}

	part2, err := randomString(digits, 4)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", part1, part2), nil
}

func GenerateAuthChallenge() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func CreateDirectories(dirs []string) error {
	for _, dir := range dirs {
		err := os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	return nil
}

func Pointer[T any](val T) *T {
	return &val
}

// TODO(patrik): Remove these
func StringPtr(val string) *string {
	return &val
}

func IntPtr(val int) *int {
	return &val
}

func UIntPtr(val uint) *uint {
	return &val
}

func Int8Ptr(val int8) *int8 {
	return &val
}

func Int16Ptr(val int16) *int16 {
	return &val
}

func Int32Ptr(val int32) *int32 {
	return &val
}

func Int64Ptr(val int64) *int64 {
	return &val
}

func UInt8Ptr(val uint8) *uint8 {
	return &val
}

func UInt16Ptr(val uint16) *uint16 {
	return &val
}

func UInt32Ptr(val uint32) *uint32 {
	return &val
}

func UInt64Ptr(val uint64) *uint64 {
	return &val
}

func PrettyDuration(d time.Duration) string {
	switch {
	case d < time.Millisecond:
		return d.String() // show microseconds as-is
	case d < time.Second:
		return d.Truncate(time.Millisecond).String() // "123ms"
	case d < time.Minute:
		return d.Truncate(time.Second).String() // "42s"
	default:
		return d.Truncate(time.Second).String() // "2h35m42s"
	}
}

func ReadToml[T any](p string) (T, error) {
	var res T

	d, err := os.ReadFile(p)
	if err != nil {
		return res, fmt.Errorf("read file: %w", err)
	}

	err = toml.Unmarshal(d, &res)
	if err != nil {
		return res, fmt.Errorf("unmarshal: %w", err)
	}

	return res, nil
}

func ReadJson[T any](p string) (T, error) {
	var res T

	d, err := os.ReadFile(p)
	if err != nil {
		return res, fmt.Errorf("read file: %w", err)
	}

	err = json.Unmarshal(d, &res)
	if err != nil {
		return res, fmt.Errorf("unmarshal: %w", err)
	}

	return res, nil
}
