package create

import (
	"errors"
	"fmt"
	"github.com/ditsuke/dpaste-cli/dpaste"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
)

const (
	exitGeneralFailure = 1
	exitNoPaste        = 2
	exitErrRead        = 3
)

// Create is a cli.Command handler to create a new paste on dpaste
func Create(c *cli.Context, client *dpaste.Dpaste) error {
	stdoutLogger := log.New(c.App.Writer, "", 0)
	errLogger := log.New(c.App.ErrWriter, "", 0)

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

	if creationResponse.Success {
		stdoutLogger.Println("paste created")

		link := creationResponse.Location
		if link == "" {
			return cli.Exit("failed to get paste link. file a bug report?", exitNoPaste)
		}
		stdoutLogger.Println("link: ", link)

		expires := creationResponse.Expiry
		if expires == "" {
			errLogger.Println("failed to get paste expiry date")
			return nil
		} else {
			stdoutLogger.Println("expires: ", expires)
		}
		return nil
	}

	return cli.Exit(fmt.Sprintf("failed to create paste: %v", creationResponse.Response.Status), exitGeneralFailure)
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
