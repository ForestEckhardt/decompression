package decompression

import (
	"compress/gzip"
	"fmt"
	"io"
)

type ArchiveFormat string

const (
	Tar   ArchiveFormat = "tar"
	GZTar ArchiveFormat = "gzip tar"
)

type Decompressor struct {
	Reader decompressTar
}

func NewDecompressor(inputReader io.Reader, archiveFormat ArchiveFormat) (Decompressor, error) {
	switch archiveFormat {
	case Tar:
		return Decompressor{Reader: NewArchiveReader(inputReader)}, nil
	case GZTar:
		gzr, err := gzip.NewReader(inputReader)
		if err != nil {
			return Decompressor{}, fmt.Errorf("failed to create gzip reader: %s", err.Error())
		}
		return Decompressor{Reader: NewArchiveReader(gzr)}, nil
	}
	return Decompressor{}, nil
}

func (d Decompressor) Decompress(destination string) error {
	err := d.Reader.UnTar(destination)
	if err != nil {
		return err
	}
	return nil
}
