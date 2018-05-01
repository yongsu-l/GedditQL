package parse

import (
	"GedditQL/server/storage"
	"fmt"
	"github.com/syohex/go-texttable"
)

// Table parses the result into a table and outputs to the console
func Table(res storage.Response) error {

	if res.Err == "" {
		tbl := &texttable.TextTable{}
		tbl.SetHeaderArr(res.Names)

		for _, v := range res.Data {
			tbl.AddRowArr(v)
		}

		fmt.Println(tbl.Draw())

	} else {
		fmt.Println(res.Err)
	}

	return nil
}
