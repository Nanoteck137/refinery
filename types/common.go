package types

import "path"

type DataDir string

func (d DataDir) String() string {
	return string(d)
}

func (d DataDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d DataDir) Playlists() string {
	return path.Join(d.String(), "playlists")
}

func (d DataDir) Playlist(id string) string {
	return path.Join(d.Playlists(), id)
}

func (d DataDir) Users() string {
	return path.Join(d.String(), "users")
}

func (d DataDir) User(id string) string {
	return path.Join(d.Users(), id)
}

func (d DataDir) Cache() string {
	return path.Join(d.String(), "cache")
}

func (d DataDir) CacheImages() ImageCacheDir {
	return ImageCacheDir(path.Join(d.Cache(), "images"))
}

func (d DataDir) CacheTranscoding() string {
	return path.Join(d.Cache(), "transcoding")
}

type ImageCacheDir string

func (d ImageCacheDir) String() string {
	return string(d)
}

func (d ImageCacheDir) Artists() string {
	return path.Join(d.String(), "artists")
}

func (d ImageCacheDir) Artist(id string) string {
	return path.Join(d.Artists(), id)
}

func (d ImageCacheDir) Albums() string {
	return path.Join(d.String(), "albums")
}

func (d ImageCacheDir) Album(id string) string {
	return path.Join(d.Albums(), id)
}

func (d ImageCacheDir) Tracks() string {
	return path.Join(d.String(), "tracks")
}

func (d ImageCacheDir) Track(id string) string {
	return path.Join(d.Tracks(), id)
}

func (d ImageCacheDir) Playlists() string {
	return path.Join(d.String(), "playlists")
}

func (d ImageCacheDir) Playlist(id string) string {
	return path.Join(d.Playlists(), id)
}

func (d ImageCacheDir) Users() string {
	return path.Join(d.String(), "users")
}

func (d ImageCacheDir) User(id string) string {
	return path.Join(d.Users(), id)
}

type Change[T any] struct {
	Value   T
	Changed bool
}

type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

type ErrorList []*Error

func (p *ErrorList) Add(message string) {
	*p = append(*p, &Error{
		Message: message,
	})
}
