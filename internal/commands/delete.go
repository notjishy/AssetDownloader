package commands

import (
	"fmt"
	"strconv"

	"golang.org/x/sync/errgroup"
	"jishe.wtf/assetdownloader/internal/records"
)

func Delete(rSlice []string) error {
	uString := "Usage: assetdownloader delete <number>"
	uString += "\nExample: assetdownloader delete 1 2 3 ..."

	if len(rSlice) < 1 {
		return fmt.Errorf("invalid arguments: %s", uString)
	}

	rs, err := records.GetRecords()
	if err != nil {
		return fmt.Errorf("error occured while getting records data: %v", err)
	}

	g := new(errgroup.Group)
	for _, r := range rSlice {
		g.Go(func(r string, rs []records.Record) func() error {
			return func() error {
				// convert to int
				rInt, err := strconv.Atoi(r)
				if err != nil {
					return fmt.Errorf("error converting to int: %v", err)
				}

				// delete the record
				err = records.Remove(rInt, rs)
				if err != nil {
					return err
				}

				return nil
			}
		}(r, rs))
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error occured while attempting to delete records: %v", err)
	}

	return nil // no errors
}
