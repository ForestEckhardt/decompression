package decompression

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func UnTar(compressionReader io.Reader, destination string) error {
	tr := tar.NewReader(compressionReader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
			// return fmt.Errorf("failed to read tar response: %s", err.Error())
		}

		path := filepath.Join(destination, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(path, hdr.FileInfo().Mode())
			if err != nil {
				panic(err)
				// return fmt.Errorf("failed to create archived directory: %s", err.Error())
			}

		case tar.TypeReg:
			fmt.Println(hdr.FileInfo().Mode())
			file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, hdr.FileInfo().Mode())
			if err != nil {
				panic(err)
				// return fmt.Errorf("failed to create file %s", err.Error())
			}
			defer file.Close()

			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}
			// err := writeStreamingFile(tr, path, )
			// if err != nil {
			// 	return fmt.Errorf("failed to stream file from archive: %s", err.Error())
			// }
		}
	}

	return nil
}
