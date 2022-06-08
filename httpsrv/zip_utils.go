package httpsrv

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func addFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "opening file %q", filename)
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return errors.Wrap(err, "performing stat()")
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return errors.Wrap(err, "populating zip fileinfo header")
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filepath.Base(filename)

	// Change to store to avoid compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Store

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return errors.Wrap(err, "creating zip file header")
	}
	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return errors.Wrap(err, "copying file to zip-archive")
	}

	return nil
}
