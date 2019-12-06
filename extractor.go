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
	XZTar ArchiveFormat = "xz tar"
)

type Extractor interface {
	Extract(destination string) error
}

type ArchiveExtractor struct {
	Reader Decompressor
}

func NewExtractor(inputReader io.Reader, archiveFormat ArchiveFormat) (Extractor, error) {
	switch archiveFormat {
	case Tar:
		return ArchiveExtractor{Reader: NewArchiveReader(inputReader)}, nil
	case GZTar:
		gzr, err := gzip.NewReader(inputReader)
		if err != nil {
			return ArchiveExtractor{}, fmt.Errorf("failed to create gzip reader: %s", err.Error())
		}
		return ArchiveExtractor{Reader: NewArchiveReader(gzr)}, nil
	}
	return ArchiveExtractor{}, nil
}

func (a ArchiveExtractor) Extract(destination string) error {
	err := a.Reader.UnTar(destination)
	if err != nil {
		return err
	}
	return nil
}
