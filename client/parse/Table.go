package parse

import (
	"GedditQL/server/storage"
	"fmt"
	"github.com/Yong-L/go-texttable"
)

// Table parses the result into a table and outputs to the console
func Table(res storage.Response) error {

	if res.Err == "" {

		tbl := &texttable.TextTable{}
		if len(res.Names) != 0 {
			tbl.SetHeaderArr(res.Names)
		}

		if len(res.Data) != 0 {
			for _, v := range res.Data {
				tbl.AddRowArr(v)
			}
		}

		if len(res.Names) != 0 || len(res.Data) != 0 {
			fmt.Println(tbl.Draw())
		}

		fmt.Println(res.Result)

	} else {
		fmt.Println(res.Err)
	}

	return nil
}
