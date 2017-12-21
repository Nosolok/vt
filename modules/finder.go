package modules

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Find(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {

		if !(f.IsDir()) {
			fmt.Println(f.Name())
			fmt.Println(path + fmt.Sprintf("%c", os.PathSeparator) + f.Name())

			fh, _ := os.Open(path + fmt.Sprintf("%c", os.PathSeparator) + f.Name())
			defer fh.Close()
			hMd5 := CalcHashMd5(fh)
			hSha1 := CalcHashSha1(fh)
			fmt.Printf("%x", hMd5)
			fmt.Println()
			fmt.Printf("%x", hSha1)
			fmt.Println()
			fmt.Println()
		}
	}
}
