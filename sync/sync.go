package sync

import (
	"fmt"
	"os"
	"time"

	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/path"
	"github.com/knqyf263/pet/snippet"
	"github.com/pkg/errors"
)

// Client manages communication with the remote Snippet repository
type Client interface {
	GetSnippet() (*Snippet, error)
	UploadSnippet(string) error
}

// Snippet is the remote snippet
type Snippet struct {
	Content   string
	UpdatedAt time.Time
}

// AutoSync syncs snippets automatically
func AutoSync(filePath path.AbsolutePath) error {
	client, err := NewSyncClient()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize API client")
	}

	snippet, err := client.GetSnippet()
	if err != nil {
		return err
	}

	fi, err := os.Stat(filePath.Get())
	if os.IsNotExist(err) || fi.Size() == 0 {
		return download(snippet.Content)
	} else if err != nil {
		return errors.Wrap(err, "Failed to get a FileInfo")
	}

	local := fi.ModTime().UTC()
	remote := snippet.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return upload(client)
	case remote.After(local):
		return download(snippet.Content)
	default:
		return nil
	}
}

// NewSyncClient returns Client
func NewSyncClient() (Client, error) {
	if config.Conf.General.Backend == "gitlab" {
		client, err := NewGitLabClient()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to initialize GitLab client")
		}
		return client, nil
	} else if config.Conf.General.Backend == "ghe" {
		client, err := NewGHEGistClient()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to initialize GHE client")
		}
		return client, nil
	}
	client, err := NewGistClient()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize Gist client")
	}
	return client, nil
}

// upload uploads snippets from the main snippet file
// to the remote repository - directories are ignored
func upload(client Client) (err error) {
	var snippets snippet.Snippets
	if err := snippets.Load(false); err != nil {
		return errors.Wrap(err, "Failed to load the local snippets")
	}

	body, err := snippets.ToString()
	if err != nil {
		return err
	}

	if err = client.UploadSnippet(body); err != nil {
		return errors.Wrap(err, "Failed to upload snippet")
	}

	fmt.Println("Upload success")
	return nil
}

// download downloads snippets from the remote repository
// and saves them to the main snippet file - directories ignored
func download(content string) error {
	var snippets snippet.Snippets
	if err := snippets.Load(false); err != nil {
		return err
	}

	body, err := snippets.ToString()
	if err != nil {
		return err
	}

	if content == body {
		// no need to download
		fmt.Println("Already up-to-date")
		return nil
	}

	fmt.Println("Download success")
	absPath, err := path.NewAbsolutePath(config.Conf.General.SnippetFile)
	if err != nil {
		return err
	}

	return os.WriteFile(
		absPath.Get(),
		[]byte(content),
		os.ModePerm,
	)
}
