package modules

import (
	"crypto/md5"
	"crypto/sha1"
	"io"
	"os"
)

func CalcHashMd5(fh *os.File) []uint8 {
	hash := md5.New()
	io.Copy(hash, fh)
	return hash.Sum(nil)
}

func CalcHashSha1(fh *os.File) []uint8 {
	hash := sha1.New()
	io.Copy(hash, fh)
	return hash.Sum(nil)
}
