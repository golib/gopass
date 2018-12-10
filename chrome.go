package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/golib/cli"
)

func chromeImporter() cli.ActionFunc {
	return func(ctx *cli.Context) (err error) {
		filename := path.Clean(ctx.String("csv"))
		if filename == "" {
			cli.ShowSubcommandHelp(ctx)
			return
		}

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
			iname   = -1
			iurl    = -1
		)
		for i, field := range fields {
			switch strings.ToUpper(field) {
			case "USER", "USERNAME":
				iuser = i

			case "PASS", "PASSWD", "PASSWORD":
				ipasswd = i

			case "NAME":
				iname = i

			case "URL", "HREF", "LINK":
				iurl = i
			}
		}
		if iuser == -1 || ipasswd == -1 || iurl == -1 {
			return errors.New("cannot find username, password and url fields")
		}

		var (
			records [][]string

			includeEmpty = ctx.Bool("include-empty")
		)
		for {
			fields, err := csvReader.Read()
			if err != nil {
				if err == io.EOF {
					break
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

			var name string
			if iname > -1 {
				name = fields[iname]
			}
			if name == "" {
				name = urlobj.Hostname()
			}

			absurl := fmt.Sprintf("%s://%s", urlobj.Scheme, urlobj.Hostname())
			if port := urlobj.Port(); len(port) > 0 {
				absurl += ":" + port
			}

			records = append(records, []string{
				name, absurl, fields[iuser], fields[ipasswd],
			})
		}

		sort.SliceStable(records, func(i int, j int) bool {
			return records[i][0] < records[j][0]
		})

		buf := bytes.NewBuffer(nil)

		csvWriter := csv.NewWriter(buf)
		csvWriter.WriteAll(records)
		csvWriter.Flush()

		chromeFilename := fmt.Sprintf("./Chrome-%s.csv", time.Now().Format("2006-01-02-15-04-05"))
		err = ioutil.WriteFile(chromeFilename, buf.Bytes(), 0644)
		if err != nil {
			log.Printf("Write %s: %v\n", chromeFilename)
		} else {
			log.Printf("Write %s: OK!\n", chromeFilename)
		}

		return err
	}
}
