package modules

import (
	"fmt"
	"io/ioutil"
	"os"
)

type File struct {
	Filename string
	HashMd5  []uint8
	HashSha1 []uint8
}

func Find(path string) []*File {
	var fileHashes []*File
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if !(f.IsDir()) {
			hMd5 := CalcHashMd5(path + fmt.Sprintf("%c", os.PathSeparator) + f.Name())
			hSha1 := CalcHashSha1(path + fmt.Sprintf("%c", os.PathSeparator) + f.Name())

			var fMeta = new(File)
			fMeta.Filename = f.Name()
			fMeta.HashMd5 = hMd5
			fMeta.HashSha1 = hSha1
			fileHashes = append(fileHashes, fMeta)
		}
	}

	return fileHashes
}
