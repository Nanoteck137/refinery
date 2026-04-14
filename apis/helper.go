package apis

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/refinery/core"
	"github.com/nanoteck137/refinery/database"
	"github.com/nanoteck137/refinery/tools/utils"
	"github.com/nanoteck137/refinery/types"
)

type UserCheckFunc func(user *database.User) error

func RequireAdmin(user *database.User) error {
	if user.Role != types.RoleSuperUser && user.Role != types.RoleAdmin {
		return InvalidAuth("user requires 'super_user' or 'admin' role")
	}

	return nil
}

func User(app core.App, c pyrin.Context, checks ...UserCheckFunc) (*database.User, error) {
	user, err := getUser(app, c)
	if err != nil {
		return nil, err
	}

	for _, check := range checks {
		err := check(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// TODO(patrik): Cleanup
func getUser(app core.App, c pyrin.Context) (*database.User, error) {
	ctx := c.Request().Context()

	apiTokenHeader := c.Request().Header.Get("X-Api-Token")
	if apiTokenHeader != "" {
		token, err := app.DB().GetApiTokenById(ctx, apiTokenHeader)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				return nil, InvalidAuth("invalid api token")
			}

			return nil, err
		}

		user, err := app.DB().GetUserById(ctx, token.UserId)
		if err != nil {
			return nil, InvalidAuth("invalid api token")
		}

		return &user, nil
	}

	authHeader := c.Request().Header.Get("Authorization")
	tokenString := utils.ParseAuthHeader(authHeader)
	if tokenString == "" {
		return nil, InvalidAuth("invalid authorization header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.Config().JwtSecret), nil
	})

	if err != nil {
		// TODO(patrik): Handle error better
		return nil, InvalidAuth("invalid authorization token")
	}

	jwtValidator := jwt.NewValidator(jwt.WithIssuedAt())

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := jwtValidator.Validate(token.Claims); err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		userId := claims["userId"].(string)
		user, err := app.DB().GetUserById(ctx, userId)
		if err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		return &user, nil
	}

	return nil, InvalidAuth("invalid authorization token")
}

func ConvertURL(c pyrin.Context, path string) string {
	host := c.Request().Host

	scheme := "http"

	h := c.Request().Header.Get("X-Forwarded-Proto")
	if h != "" {
		scheme = h
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

const (
	IMAGE_ORIGINAL = "original.png"
	IMAGE_SMALL    = "128.png"
	IMAGE_MEDIUM   = "256.png"
	IMAGE_LARGE    = "512.png"
)

func ConvertArtistCoverURL(c pyrin.Context, artistId string) types.Images {
	first := "/files/artists/images/" + artistId + "/"
	return types.Images{
		Original: ConvertURL(c, first+IMAGE_ORIGINAL),
		Small:    ConvertURL(c, first+IMAGE_SMALL),
		Medium:   ConvertURL(c, first+IMAGE_MEDIUM),
		Large:    ConvertURL(c, first+IMAGE_LARGE),
	}
}

func ConvertAlbumCoverURL(c pyrin.Context, albumId string) types.Images {
	first := "/files/albums/images/" + albumId + "/"
	return types.Images{
		Original: ConvertURL(c, first+IMAGE_ORIGINAL),
		Small:    ConvertURL(c, first+IMAGE_SMALL),
		Medium:   ConvertURL(c, first+IMAGE_MEDIUM),
		Large:    ConvertURL(c, first+IMAGE_LARGE),
	}
}

func ConvertPlaylistCoverURL(c pyrin.Context, playlistId string) types.Images {
	first := "/files/playlists/images/" + playlistId + "/"
	return types.Images{
		Original: ConvertURL(c, first+IMAGE_ORIGINAL),
		Small:    ConvertURL(c, first+IMAGE_SMALL),
		Medium:   ConvertURL(c, first+IMAGE_MEDIUM),
		Large:    ConvertURL(c, first+IMAGE_LARGE),
	}
}

func ConvertUserPictureURL(c pyrin.Context, userId string) types.Images {
	first := "/files/users/images/" + userId + "/"
	return types.Images{
		Original: ConvertURL(c, first+IMAGE_ORIGINAL),
		Small:    ConvertURL(c, first+IMAGE_SMALL),
		Medium:   ConvertURL(c, first+IMAGE_MEDIUM),
		Large:    ConvertURL(c, first+IMAGE_LARGE),
	}
}

func getPageParams(q url.Values, defaultPerPage int) types.PageParams {
	perPage := defaultPerPage
	page := 0

	if s := q.Get("perPage"); s != "" {
		i, _ := strconv.Atoi(s)
		if i > 0 {
			perPage = i
		}
	}

	if s := q.Get("page"); s != "" {
		i, _ := strconv.Atoi(s)
		page = i
	}

	return types.PageParams{
		PerPage: perPage,
		Page:    page,
	}
}

func getFilterParams(q url.Values) types.FilterParams {
	return types.FilterParams{
		Filter: q.Get("filter"),
		Sort:   q.Get("sort"),
	}
}

func formatTime(ms int64) string {
	return time.UnixMilli(ms).UTC().Format(time.RFC3339)
}
