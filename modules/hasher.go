package modules

import (
	"crypto/md5"
	"crypto/sha1"
	"io"
	"log"
	"os"
)

func CalcHashMd5(file string) []uint8 {
	fh, err := os.Open(file)
	defer fh.Close()
	if err != nil {
		log.Fatal(err)
	}

	hash := md5.New()
	io.Copy(hash, fh)
	return hash.Sum(nil)
}

func CalcHashSha1(file string) []uint8 {
	fh, err := os.Open(file)
	defer fh.Close()
	if err != nil {
		log.Fatal(err)
	}

	hash := sha1.New()
	io.Copy(hash, fh)
	return hash.Sum(nil)
}
