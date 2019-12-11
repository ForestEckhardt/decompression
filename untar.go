package decompression

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//go:generate faux --interface decompressTar --output fakes/decompress_tar.go
type decompressTar interface {
	UnTar(destination string) error
	GetReader() io.Reader
}

type ArchiveReader struct {
	Reader io.Reader
}

func NewArchiveReader(inputReader io.Reader) ArchiveReader {
	return ArchiveReader{Reader: inputReader}
}

func (a ArchiveReader) UnTar(destination string) error {
	tr := tar.NewReader(a.Reader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar response: %s", err.Error())
		}

		path := filepath.Join(destination, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(path, hdr.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create archived directory: %s", err.Error())
			}

		case tar.TypeReg:
			file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, hdr.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("failed to create archived file %s", err.Error())
			}

			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}

			err = file.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a ArchiveReader) GetReader() io.Reader {
	return a.Reader
}
