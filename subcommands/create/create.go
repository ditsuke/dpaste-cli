package create

import (
	"dpaste-cli/lib/dpaste"
	"errors"
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
	"io"
	"os"
)

const (
	exitGeneralFailure = 1
	exitNoPaste        = 2
	exitErrRead        = 3
)

// Create is a cli.Command handler to create a new paste on dpaste
func Create(c *cli.Context, client *dpaste.Dpaste) error {
	content, err := getContent(c)
	if err != nil {
		return cli.Exit("nothing to read: "+err.Error(), exitErrRead)
	}

	creationRequest := dpaste.CreateRequest{
		Content:    content,
		ExpiryDays: c.Int("expiry_days"),
		Syntax:     c.String("syntax"),
		Title:      c.String("Title"),
	}

	creationResponse, err := client.Create(creationRequest)
	if err != nil {
		return cli.Exit(err.Error(), exitGeneralFailure)
	}

	writer := c.App.Writer

	if creationResponse.Success {
		_, _ = fmt.Fprintf(writer, "Link: %q.\nExpires In: %q", creationResponse.Location, creationResponse.Expiry)
		return nil
	}

	// Probably though, we should be printing custom here with error writer and yada yada
	return cli.Exit(fmt.Sprintf("Failed to create paste: %v", creationResponse.Response.Status), exitGeneralFailure)
}

// getContent returns an io.Reader for the content to upload, defaulting to standard
// input before checking for a file from the "file" flag.
func getContent(c *cli.Context) (io.Reader, error) {
	// we default to stdin, so we check if we have a pipe to read from
	if stdinIsTty := isatty.IsTerminal(os.Stdin.Fd()); !stdinIsTty {
		return os.Stdin, nil
	}

	// else, we try to open the file from the "file" flag
	if file := c.String("file"); file != "" {
		if fInfo, err := os.Stat(file); err != nil || fInfo.Size() == 0 {
			if err == nil {
				err = errors.New("file is empty")
			}
			return nil, err
		}

		fileStream, err := os.Open(file)
		return fileStream, err
	}
	return nil, errors.New("no file provided")
}
