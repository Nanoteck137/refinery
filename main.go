package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/kr/pretty"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

// Program Dependencies: git nix skopeo

type TagInfoCommit struct {
	Url     string    `json:"url"`
	Sha     string    `json:"sha"`
	Created time.Time `json:"created"`
}

type TagInfoArchiveDownloadCount struct {
	Zip   int `json:"zip"`
	TarGz int `json:"tar_gz"`
}

type TagInfo struct {
	Name                 string                      `json:"name"`
	Message              string                      `json:"message"`
	Id                   string                      `json:"id"`
	Commit               TagInfoCommit               `json:"commit"`
	ZipballUrl           string                      `json:"zipball_url"`
	TarballUrl           string                      `json:"tarball_url"`
	ArchiveDownloadCount TagInfoArchiveDownloadCount `json:"archive_download_count"`
}

func fetchTags(url string) ([]TagInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var res []TagInfo
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// NOTE(patrik): This is not everything
type BranchInfoCommit struct {
	Id        string    `json:"id"`
	Message   string    `json:"message"`
	Url       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
}

// NOTE(patrik): This is not everything
type BranchInfo struct {
	Name   string           `json:"name"`
	Commit BranchInfoCommit `json:"commit"`
}

func fetchBranches(url string) ([]BranchInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var res []BranchInfo
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type Repo struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`

	HtmlUrl     string `json:"html_url"`
	Url         string `json:"url"`
	SshUrl      string `json:"ssh_url"`
	CloneUrl    string `json:"clone_url"`
	OriginalUrl string `json:"original_url"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ArchivedAt time.Time `json:"archived_at"`

	// "id": 2,
	//	"owner": {
	//	  "id": 1,
	//	  "login": "nanoteck137",
	//	  "login_name": "",
	//	  "source_id": 0,
	//	  "full_name": "Patrik M. Rosenström",
	//	  "email": "nanoteck137@noreply.localhost",
	//	  "avatar_url": "https://forgejo.nanoteck137.net/avatars/13eae1986176a31526e4530737c647064855d4268bb6359932f5d76c9d9fcf0d",
	//	  "html_url": "https://forgejo.nanoteck137.net/nanoteck137",
	//	  "language": "",
	//	  "is_admin": false,
	//	  "last_login": "0001-01-01T00:00:00Z",
	//	  "created": "2026-04-09T20:44:48+02:00",
	//	  "restricted": false,
	//	  "active": false,
	//	  "prohibit_login": false,
	//	  "location": "",
	//	  "pronouns": "",
	//	  "website": "",
	//	  "description": "",
	//	  "visibility": "public",
	//	  "followers_count": 0,
	//	  "following_count": 0,
	//	  "starred_repos_count": 0,
	//	  "username": "nanoteck137"
	//	},
	//
	// "name": "tunebook",
	// "full_name": "nanoteck137/tunebook",
	// "description": "",
	// "empty": false,
	// "private": false,
	// "fork": false,
	// "template": false,
	// "parent": null,
	// "mirror": false,
	// "size": 4299,
	// "language": "Go",
	// "languages_url": "https://forgejo.nanoteck137.net/api/v1/repos/nanoteck137/tunebook/languages",
	// "html_url": "https://forgejo.nanoteck137.net/nanoteck137/tunebook",
	// "url": "https://forgejo.nanoteck137.net/api/v1/repos/nanoteck137/tunebook",
	// "link": "",
	// "ssh_url": "ssh://git@git.nanoteck137.net/nanoteck137/tunebook.git",
	// "clone_url": "https://forgejo.nanoteck137.net/nanoteck137/tunebook.git",
	// "original_url": "https://github.com/Nanoteck137/tunebook",
	// "website": "",
	// "stars_count": 0,
	// "forks_count": 0,
	// "watchers_count": 1,
	// "open_issues_count": 0,
	// "open_pr_counter": 0,
	// "release_counter": 0,
	// "default_branch": "main",
	// "archived": false,
	// "created_at": "2026-04-09T21:03:22+02:00",
	// "updated_at": "2026-04-14T00:27:56+02:00",
	// "archived_at": "1970-01-01T01:00:00+01:00",
	//
	//	"permissions": {
	//	  "admin": false,
	//	  "push": false,
	//	  "pull": true
	//	},
	//
	// "has_issues": true,
	//
	//	"internal_tracker": {
	//	  "enable_time_tracker": true,
	//	  "allow_only_contributors_to_track_time": true,
	//	  "enable_issue_dependencies": true
	//	},
	//
	// "has_wiki": true,
	// "wiki_branch": "main",
	// "globally_editable_wiki": false,
	// "has_pull_requests": true,
	// "has_projects": true,
	// "has_releases": true,
	// "has_packages": true,
	// "has_actions": true,
	// "ignore_whitespace_conflicts": false,
	// "allow_merge_commits": true,
	// "allow_rebase": true,
	// "allow_rebase_explicit": true,
	// "allow_squash_merge": true,
	// "allow_fast_forward_only_merge": true,
	// "allow_rebase_update": true,
	// "default_delete_branch_after_merge": false,
	// "default_merge_style": "merge",
	// "default_allow_maintainer_edit": false,
	// "default_update_style": "merge",
	// "avatar_url": "",
	// "internal": false,
	// "mirror_interval": "",
	// "object_format_name": "sha1",
	// "mirror_updated": "0001-01-01T00:00:00Z",
	// "repo_transfer": null,
	// "topics": []
}

func fetchRepo(url string) (*Repo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var res Repo
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func randomDirName(name string) string {
	return fmt.Sprintf("%s-%d-%06d", name, time.Now().Unix(), rand.IntN(1_000_000))
}

type Config struct {
	Owner    string `toml:"owner"`
	RepoName string `toml:"repoName"`

	RegistryImage    string `toml:"registryImage"`
	RegistryUsername string `toml:"registryUsername"`
	RegistryPassword string `toml:"registryPassword"`
}

var rootCmd = &cobra.Command{
	Use: "refinery",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var newCmd = &cobra.Command{
	Use: "new <NAME>",
	Short: "Create new config",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		conf := Config{
			Owner:            "<REPO OWNER>",
			RepoName:         "<REPO NAME>",
			RegistryImage:    "<NAME OF THE FINAL IMAGE>",
			RegistryUsername: "<USERNAME FOR REGISTRY AUTH>",
			RegistryPassword: "<PASSWORD FOR REGISTRY AUTH>",
		}

		d, err := toml.Marshal(conf)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(name + ".toml", d, 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var publishCmd = &cobra.Command{
	Use: "publish <CONFIG>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := args[0]

		d, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		var config Config
		err = toml.Unmarshal(d, &config)
		if err != nil {
			log.Fatal(err)
		}

		// TODO(patrik): Config
		base := "https://forgejo.nanoteck137.net"

		owner := config.Owner
		repoName := config.RepoName

		// TODO(patrik): Config
		registryImage := config.RegistryImage

		// TODO(patrik): Config
		registryUsername := config.RegistryUsername
		registryPassword := config.RegistryPassword

		url := fmt.Sprintf("%s/api/v1/repos/%s/%s", base, owner, repoName)
		repo, err := fetchRepo(url)
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(repo)

		// Fetch tags/branches
		url = fmt.Sprintf("%s/api/v1/repos/%s/%s/tags", base, owner, repoName)
		tags, err := fetchTags(url)
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(tags)

		url = fmt.Sprintf("%s/api/v1/repos/%s/%s/branches", base, owner, repoName)
		branches, err := fetchBranches(url)
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(branches)

		// Select tag/branch
		// TODO():
		type Choice struct {
			Type  string
			Value string
		}

		var options []huh.Option[Choice]

		for _, branch := range branches {
			options = append(options, huh.NewOption("branch:"+branch.Name, Choice{
				Type:  "branch",
				Value: branch.Name,
			}))
		}

		for _, tag := range tags {
			options = append(options, huh.NewOption("tag:"+tag.Name, Choice{
				Type:  "tag",
				Value: tag.Name,
			}))
		}

		var choice Choice

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[Choice]().
					Options(options...).
					Value(&choice),
			),
		)

		err = form.Run()
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(choice)

		// Clone repo
		workDir := "./work"

		buildDir := path.Join(workDir, randomDirName(repo.Name))
		fmt.Printf("buildDir: %v\n", buildDir)
		defer func() {
			os.RemoveAll(buildDir)
		}()

		runGit := func(dir string, params ...string) error {
			cmd := exec.Command("git", params...)
			cmd.Dir = dir

			env := os.Environ()
			env = append(env,
				"GIT_CONFIG_SYSTEM=/dev/null",
				"GIT_CONFIG_GLOBAL=/dev/null",
				"GIT_CONFIG_NOSYSTEM=1",

				"GIT_TERMINAL_PROMPT=0",
				"GIT_ASKPASS=/bin/false",
			)

			cmd.Env = env
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				return err
			}

			return nil
		}

		// git clone --depth 1 "$REPO_URL" "$BUILD_DIR"
		err = runGit("", "clone", "--depth", "1", repo.CloneUrl, buildDir)
		if err != nil {
			log.Fatal(err)
		}

		// if [ "$TYPE" = "tag" ]; then
		//   git fetch --tags
		//   git checkout "tags/$NAME"
		// else
		//   git checkout "$NAME"
		// fi

		tag := ""

		switch choice.Type {
		case "tag":
			err = runGit(buildDir, "fetch", "--tags")
			if err != nil {
				log.Fatal(err)
			}

			err = runGit(buildDir, "checkout", "tags/"+choice.Value)
			if err != nil {
				log.Fatal(err)
			}

			tag = choice.Value
		case "branch":
			err = runGit(buildDir, "fetch", "origin", choice.Value)
			if err != nil {
				log.Fatal(err)
			}

			err = runGit(buildDir, "checkout", "FETCH_HEAD")
			if err != nil {
				log.Fatal(err)
			}

			tag = choice.Value
		default:
			panic("unknown choice type: " + choice.Type)
		}

		// Fetch version
		// TODO(patrik)

		nixBuild := func() (string, error) {
			var res bytes.Buffer

			cmd := exec.Command("nix", "build", "--no-link", "--print-out-paths", ".#docker")
			cmd.Dir = buildDir
			cmd.Stdout = &res
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				return "", err
			}

			return strings.TrimSpace(res.String()), nil
		}

		// Build project
		imagePath, err := nixBuild()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("imagePath: %v\n", imagePath)

		// Push to registry

		runSkopeo := func(params ...string) error {
			p := []string{
				"--insecure-policy",
			}

			p = append(p, params...)

			cmd := exec.Command("skopeo", p...)
			cmd.Dir = buildDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				return err
			}

			return err
		}

		// err = runSkopeo("inspect", "docker-archive:"+imagePath)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		targetImage := registryImage + ":" + tag
		err = runSkopeo(
			"copy",
			"docker-archive:"+imagePath,
			"docker://"+targetImage,
			"--dest-username", registryUsername,
			"--dest-password", registryPassword,
		)
		if err != nil {
			log.Fatal(err)
		}

		// Tag latest in registry

		tagLatest := false
		form = huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Tag Latest?").
					Value(&tagLatest),
			),
		)

		err = form.Run()
		if err != nil {
			log.Fatal(err)
		}

		if tagLatest {
			dstImage := registryImage + ":latest"

			err = runSkopeo(
				"copy",
				"docker://"+targetImage,
				"docker://"+dstImage,
				"--dest-username", registryUsername,
				"--dest-password", registryPassword,
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd, publishCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}
