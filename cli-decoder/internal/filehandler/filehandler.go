package filehandler

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const md5hashesFilename = "md5_hashes"

type Filehandler interface {
	SaveBytes(bts []byte) error
	Release()
	getFilename(string) string
}

func loadhashes() ([]string, *os.File) {
	f, err := os.OpenFile(md5hashesFilename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	if info, err := f.Stat(); err != nil {
		panic(err)
	} else if info.Size() == 0 {
		return []string{}, f
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	hashes := strings.Split(string(data), ":")

	return hashes, f
}

type commonFilehandler struct {
	hashes        []string
	hashesFile    *os.File
	fileExtention string
}

func (fh *commonFilehandler) addHash(md5hash string) error {
	fh.hashes = append(fh.hashes, md5hash)

	_, err := fh.hashesFile.Write([]byte(md5hash + ":"))
	if err != nil {
		return err
	}

	return nil
}

func (fh *commonFilehandler) containsMD5Hash(hash string) bool {
	for i := range fh.hashes {
		if fh.hashes[i] == hash {
			return true
		}
	}

	return false
}

func (fh *commonFilehandler) Release() {
	fh.hashesFile.Close()
}

func (fh *commonFilehandler) getFilename(name string) string {
	return fmt.Sprintf("%v.%v", name, fh.fileExtention)
}

func (fh *commonFilehandler) SaveBytes(bts []byte) (deferErr error) {
	checksum := fmt.Sprintf("%x", md5.Sum(bts))

	if fh.containsMD5Hash(checksum) {
		return fmt.Errorf("the data is already saved to a file")
	}

	filename := fh.getFilename(checksum)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("%w: %v", err, filename)
	}

	defer func() {
		if err := f.Close(); err != nil {
			deferErr = err
		}
	}()

	_, err = f.Write(bts)
	if err != nil {
		return err
	}

	err = fh.addHash(checksum)
	if err != nil {
		return err
	}

	return
}

// JSON File Handler
type JSONFilehandler struct {
	commonFilehandler
}

func NewJSONFilehandler() *JSONFilehandler {
	h, f := loadhashes()

	fh := &JSONFilehandler{}
	fh.commonFilehandler.hashes = h
	fh.commonFilehandler.hashesFile = f
	fh.fileExtention = "json"

	return fh
}

// XML File Handler
type XMLFilehandler struct {
	commonFilehandler
}

func NewXMLFilehandler() *XMLFilehandler {
	h, f := loadhashes()

	fh := &XMLFilehandler{}
	fh.commonFilehandler.hashes = h
	fh.commonFilehandler.hashesFile = f
	fh.fileExtention = "xml"

	return fh
}
