package apis

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/nanoteck137/refinery/core"
	"github.com/nanoteck137/pyrin"
)

func InstallFilesHandlers(app core.App, g pyrin.Group) {
	g.Register(
		pyrin.NormalHandler{
			Name:   "GetUserImage",
			Method: http.MethodGet,
			Path:   "/users/images/:userId/:image",
			HandlerFunc: func(c pyrin.Context) error {
				userId := c.Param("userId")
				image := c.Param("image")

				ext := path.Ext(image)
				name := strings.TrimRight(image, ext)

				imageType, ok := app.ImageService().GetImageTypeFromExt(ext)
				if !ok {
					// TODO(patrik): Better error
					return errors.New("unsupported image ext")
				}

				ctx := c.Request().Context()

				p, err := app.ImageService().GetUserImage(ctx, userId, name, imageType)
				if err != nil {
					return err
				}

				f := os.DirFS(path.Dir(p))
				return pyrin.ServeFile(c, f, path.Base(p))
			},
		},
	)
}
