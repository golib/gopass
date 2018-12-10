package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/golib/cli"
)

func keychainImporter() cli.ActionFunc {
	return func(ctx *cli.Context) (err error) {
		filename := path.Clean(ctx.String("csv"))
		if filename == "" {
			cli.ShowSubcommandHelp(ctx)
			return
		}

		isUpdate := ctx.Bool("update")

		reader, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		if err != nil {
			return
		}

		csvReader := csv.NewReader(reader)

		fields, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil

				log.Println("Success!")
			}
			return err
		}

		var (
			iuser   = -1
			ipasswd = -1
			iurl    = -1
		)
		for i, field := range fields {
			switch strings.ToUpper(field) {
			case "USER", "USERNAME":
				iuser = i

			case "PASS", "PASSWD", "PASSWORD":
				ipasswd = i

			case "URL", "HREF", "LINK":
				iurl = i
			}
		}
		if iuser == -1 || ipasswd == -1 || iurl == -1 {
			return errors.New("cannot find username, password and url fields")
		}

		var (
			includeEmpty = ctx.Bool("include-empty")
		)
		for {
			fields, err := csvReader.Read()
			if err != nil {
				if err == io.EOF {
					err = nil

					log.Println("Success!")
				}
				return err
			}

			if !includeEmpty && (len(fields[iuser]) == 0 || len(fields[ipasswd]) == 0) {
				log.Printf("[%s](%s): ignore credential with empty username or password\n", fields[iuser], fields[iurl])
				continue
			}

			urlobj, err := url.Parse(fields[iurl])
			if err != nil {
				log.Printf("[%s](%s): cannot parse url\n", fields[iuser], fields[iurl])
				continue
			}

			err = addSecurityOfInternetPassword(fields[iuser], fields[ipasswd], urlobj, isUpdate)
			if err != nil {
				log.Printf("[%s](%s): add security with %v\n", fields[iuser], fields[iurl], err)
			} else {
				log.Printf("[%s](%s): OK!\n", fields[iuser], fields[iurl])
			}
		}

		return nil
	}
}

func addSecurityOfInternetPassword(username, password string, urlobj *url.URL, isUpdate bool) error {
	if len(username) == 0 || len(password) == 0 {
		return errors.New("both username and password are required")
	}

	cmd := exec.Command(
		"security",
		"add-internet-password",
		"-a", username,
		"-w", password,
		"-l", fmt.Sprintf("%s (%s)", urlobj.Hostname(), username),
		"-s", urlobj.Hostname(),
		"-p", urlobj.Path,
		"-t", "form",
		"-r", fmt.Sprintf("%-4s", urlobj.Scheme[:4]),
		"-T", "/Applications/Safari.app",
	)
	if isUpdate {
		cmd.Args = append(cmd.Args, "-U")
	}

	_, err := cmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 45") {
			err = fmt.Errorf("Duplicated")
		}

		return err
	}

	return nil
}
