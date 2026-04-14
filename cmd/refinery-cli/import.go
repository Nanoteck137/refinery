package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/refinery/cmd/refinery-cli/api"
	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3"
)

var importCmd = &cobra.Command{
	Use: "import",
}

var dialect = ember.SqliteDialect()

func runImport(dbFile, apiAddr, apiToken string) error {
	dbUrl := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=ON&_serialized=1&_synchronous=NORMAL", dbFile)
	db, err := ember.OpenDatabase("sqlite3", dbUrl)
	if err != nil {
		return err
	}

	client := api.New(apiAddr)
	client.Headers.Set("X-Api-Token", apiToken)

	ctx := context.Background()

	type Playlist struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	query := dialect.From("playlists").
		Select(
			"playlists.id",
			"playlists.name",
		)

	playlists, err := ember.Multiple[Playlist](db, ctx, query)
	if err != nil {
		return err
	}

	fmt.Println("Found playlists:")
	for i, playlist := range playlists {
		fmt.Printf("[%d] %s\n", i+1, playlist.Name)
	}

	fmt.Println("\nEnter playlist numbers to import (comma separated, e.g. 1,3,5) or 'all':")
	var input string
	fmt.Scanln(&input)

	var selectedIndices []int
	if input == "all" {
		for i := range playlists {
			selectedIndices = append(selectedIndices, i)
		}
	} else {
		for part := range strings.SplitSeq(input, ",") {
			part = strings.TrimSpace(part)
			num, err := strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("invalid input: %s", part)
			}
			if num < 1 || num > len(playlists) {
				return fmt.Errorf("invalid playlist number: %d", num)
			}
			selectedIndices = append(selectedIndices, num-1)
		}
	}

	type ImportMode int
	const (
		ImportModePlaylist ImportMode = iota
		ImportModeFavorites
	)

	type SelectedPlaylist struct {
		Index int
		Mode  ImportMode
	}

	selectedPlaylists := make([]SelectedPlaylist, 0, len(selectedIndices))
	for _, idx := range selectedIndices {
		playlist := playlists[idx]
		fmt.Printf("\nPlaylist: %s\n", playlist.Name)
		fmt.Println("How do you want to import?")
		fmt.Println("  [p] - Create a new playlist")
		fmt.Println("  [f] - Add tracks to favorites")
		fmt.Printf("Choice: ")

		var choice string
		fmt.Scanln(&choice)
		choice = strings.TrimSpace(strings.ToLower(choice))

		mode := ImportModePlaylist
		if choice == "f" {
			mode = ImportModeFavorites
		}

		selectedPlaylists = append(selectedPlaylists, SelectedPlaylist{
			Index: idx,
			Mode:  mode,
		})
	}

	type PlaylistItem struct {
		TrackId   string `db:"track_id"`
		TrackName string `db:"track_name"`
	}

	for _, sp := range selectedPlaylists {
		playlist := playlists[sp.Index]

		fmt.Printf("\nImporting playlist: %s\n", playlist.Name)

		query := dialect.From("playlist_items").
			Select(
				"playlist_items.track_id",
				goqu.I("tracks.name").As("track_name"),
			).
			Join(
				goqu.I("tracks"),
				goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
			).
			Where(goqu.I("playlist_items.playlist_id").Eq(playlist.Id)).
			Order(goqu.I("playlist_items.rowid").Asc())

		items, err := ember.Multiple[PlaylistItem](db, ctx, query)
		if err != nil {
			return err
		}

		if sp.Mode == ImportModeFavorites {
			for _, item := range items {
				slog.Info("Favoriting track", "track", item.TrackName, "trackId", item.TrackId)

				_, err := client.FavoriteTrack(item.TrackId, api.Options{})
				if err != nil {
					var apiErr *api.ApiError[any]
					if errors.As(err, &apiErr) {
						if apiErr.Type == "TRACK_NOT_FOUND" {
							slog.Warn("Track not found", "track", item.TrackName, "trackId", item.TrackId)
							continue
						}
					}

					return err
				}

				slog.Info("Success", "track", item.TrackName, "trackId", item.TrackId)
			}
		} else {
			serverPlaylist, err := client.CreatePlaylist(api.CreatePlaylistBody{
				Name: playlist.Name,
			}, api.Options{})
			if err != nil {
				return err
			}

			for _, item := range items {
				slog.Info("Trying to add track", "track", item.TrackName, "trackId", item.TrackId)

				_, err := client.AddItemToPlaylist(serverPlaylist.Id, api.AddItemToPlaylistBody{
					TrackId: item.TrackId,
				}, api.Options{})
				if err != nil {
					var apiErr *api.ApiError[any]
					if errors.As(err, &apiErr) {
						if apiErr.Type == "TRACK_NOT_FOUND" {
							slog.Warn("Track not found", "track", item.TrackName, "trackId", item.TrackId)
							continue
						}
					}

					return err
				}

				slog.Info("Success", "track", item.TrackName, "trackId", item.TrackId)
			}
		}

		fmt.Printf("Imported playlist: %s\n", playlist.Name)
	}

	fmt.Println("\nDone!")
	return nil
}

var importPlaylistCmd = &cobra.Command{
	Use:  "playlist <OLD_DATABASE>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiAddr, _ := cmd.Flags().GetString("api-addr")
		apiToken, _ := cmd.Flags().GetString("api-token")

		err := runImport(args[0], apiAddr, apiToken)
		if err != nil {
			slog.Error("failed to import", "err", err)
			return
		}
	},
}

func init() {
	importPlaylistCmd.Flags().String("api-addr", "", "API address")
	importPlaylistCmd.MarkFlagRequired("api-addr")

	importPlaylistCmd.Flags().String("api-token", "", "API token")
	importPlaylistCmd.MarkFlagRequired("api-token")

	importCmd.AddCommand(importPlaylistCmd)

	rootCmd.AddCommand(importCmd)
}
