package file

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"regexp"

	search "github.com/geoirb/face-search/internal/face-search"
)

// Facade for work with file
type Facade struct {
	regexpfileName *regexp.Regexp
	downloadDir    string
}

// NewFacade ...
func NewFacade(downloadDir string) *Facade {
	f := &Facade{
		downloadDir: downloadDir,
	}
	f.regexpfileName, _ = regexp.Compile(`\/([^\.\/]*\.[a-z]{3})[\W_]`)
	return f
}

// GetPath of file.
func (f *Facade) GetPath(file search.File) (path string, err error) {
	fileName := f.regexpfileName.FindAllStringSubmatch(file.URL, -1)
	if len(fileName) != 1 || len(fileName[0]) != 2 {
		err = search.ErrFileNameNotFound
		return
	}
	path = f.downloadDir + fileName[0][1]
	if err = f.download(path, string(file.URL)); err != nil {
		return
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = search.ErrFileNotFound
	}
	return
}

// Delete file.
func (f *Facade) Delete(file search.File) error {
	return os.Remove(file.Path)
}

// GetHah by file.
func (f *Facade) GetHash(fl search.File) (hash string, err error) {
	file, err := os.Open(fl.Path)
	if err != nil {
		return
	}
	defer file.Close()

	h := sha256.New()
	if _, err = io.Copy(h, file); err != nil {
		return
	}

	hashHex := h.Sum(nil)
	hashBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(hashHex)))
	base64.StdEncoding.Encode(hashBase64, hashHex)

	hash = string(hashBase64)
	return
}

func (f *Facade) download(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
