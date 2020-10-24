package statikdata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	_ "github.com/rakyll/statik/fs"

	"github.com/shakahl/gqlassist/internal/utils"
)

var (
	statikFs http.FileSystem
)

func init() {

}

func Init() {
	if statikFs != nil {
		return
	}
	sfs, err := fs.New()
	if err != nil {
		panic(errors.Wrapf(err, "statikdata: error during intialization"))
	}
	statikFs = sfs

	dumpRegisteredAssets()
}

func GetFileSystem() http.FileSystem {
	if statikFs != nil {
		return statikFs
	}
	Init()
	return statikFs
}

func Open(name string) (http.File, error) {
	return GetFileSystem().Open(name)
}

func MustOpen(name string) http.File {
	file, err := GetFileSystem().Open(name)
	if err != nil {
		panic(err)
	}
	return file
}

// ReadDir reads the contents of the directory associated with file and
// returns a slice of FileInfo values in directory order.
func ReadDir(name string) ([]os.FileInfo, error) {
	f, err := GetFileSystem().Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdir(-1)
}

// Stat returns the FileInfo structure describing file.
func Stat(name string) (os.FileInfo, error) {
	f, err := GetFileSystem().Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

// ReadFile reads the file named by path from fs and returns the contents.
func ReadFile(path string) ([]byte, error) {
	rc, err := GetFileSystem().Open(path)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

func MustReadFile(path string) []byte {
	s, err := ReadFile(path)
	if err != nil {
		panic(err)
	}
	return s
}

// ReadFileString reads the file named by path from fs and returns the contents.
func ReadFileString(path string) (string, error) {
	buf, err := ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func MustReadFileString(path string) string {
	s, err := ReadFileString(path)
	if err != nil {
		panic(err)
	}
	return s
}

// Exists reports whether the named file or directory exists in http.FileSystem.
func Exists(name string) bool {
	if _, err := Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func MustExists(name string) bool {
	if !Exists(name) {
		panic(errors.New("statikdata: required asset does not exists: " + name))
	}
	return true
}

func dumpRegisteredAssets() {
	var err error
	var files []string
	walkerFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("  - %s (size: %s)\n", path, utils.FormatFileSize(info.Size()))
		files = append(files, path)
		return nil
	}
	err = fs.Walk(statikFs, "/", walkerFn)
	if err != nil {
		panic(errors.Wrapf(err, "statikdata: error while trying to list registered assets"))
	}
}
