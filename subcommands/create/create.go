package create

import (
	"dpaste-cli/lib/dpaste"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
)

const (
	exitGeneralFailure = 1
	exitNoPaste        = 2
	exitErrRead        = 3
)

// Create is a cli.Command handler to create a new paste on dpaste
func Create(c *cli.Context, client *dpaste.Dpaste) error {
	var content []byte
	var err error

	info, _ := os.Stdin.Stat()

	// We default to stdin
	if info.Mode()&os.ModeCharDevice == os.ModeCharDevice && info.Size() > 0 {
		_, err = os.Stdin.Read(content)
	}

	// If we have a file flag we read from the file
	if file := c.String("file"); file != "" {
		fileStream, err := os.Open(file)
		if err != nil {
			str := "error reading file"
			if errors.Is(err, os.ErrNotExist) {
				str = "file does not exist"
			}
			return cli.Exit(str, exitErrRead)
		}
		content, err = ioutil.ReadAll(fileStream)
	}

	if err != nil {
		return err
	}

	if len(content) == 0 {
		return cli.Exit("nothing to paste", exitNoPaste)
	}

	request := dpaste.CreateRequest{
		Content:    string(content),
		ExpiryDays: c.Int("expiry_days"),
		Syntax:     c.String("syntax"),
		Title:      c.String("Title"),
	}

	response, err := client.Create(request)

	if err != nil {
		return cli.Exit(err.Error(), exitGeneralFailure)
	}

	writer := c.App.Writer

	if response.Success {
		_, _ = fmt.Fprintf(writer, "Link: %q.\nExpires In: %q", response.Location, response.Expiry)
		return nil
	}

	// Probably though, we should be printing custom here with error writer and yada yada
	return cli.Exit(fmt.Sprintf("Failed to create paste: %v", response.Response.Status), exitGeneralFailure)
}
