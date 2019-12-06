package decompression_test

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ForestEckhardt/decompression"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testUnTar(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect func(interface{}, ...interface{}) Assertion
	)

	it.Before(func() {
		Expect = NewWithT(t).Expect
	})

	context("UnTar", func() {
		var (
			tempDir string
			reader  io.Reader
		)

		it.Before(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "decompression")
			Expect(err).NotTo(HaveOccurred())

			reader = bytes.NewReader(nil)
			buffer := bytes.NewBuffer(nil)
			tw := tar.NewWriter(buffer)

			Expect(tw.WriteHeader(&tar.Header{Name: "some-dir", Mode: 0755, Typeflag: tar.TypeDir})).To(Succeed())
			_, err = tw.Write(nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(tw.WriteHeader(&tar.Header{Name: filepath.Join("some-dir", "some-other-dir"), Mode: 0755, Typeflag: tar.TypeDir})).To(Succeed())
			_, err = tw.Write(nil)
			Expect(err).NotTo(HaveOccurred())

			nestedFile := filepath.Join("some-dir", "some-other-dir", "some-file")
			Expect(tw.WriteHeader(&tar.Header{Name: nestedFile, Mode: 0755, Size: int64(len(nestedFile))})).To(Succeed())
			_, err = tw.Write([]byte(nestedFile))
			Expect(err).NotTo(HaveOccurred())

			for _, file := range []string{"first", "second", "third"} {
				Expect(tw.WriteHeader(&tar.Header{Name: file, Mode: 0755, Size: int64(len(file))})).To(Succeed())
				_, err = tw.Write([]byte(file))
				Expect(err).NotTo(HaveOccurred())
			}

			reader = bytes.NewReader(buffer.Bytes())

			Expect(tw.Close()).To(Succeed())

		})

		it.After(func() {
			// Expect(os.RemoveAll(tempDir)).To(Succeed())
		})

		it("downloads the dependency and unpackages it into the path", func() {
			var err error
			err = decompression.UnTar(reader, tempDir)
			Expect(err).ToNot(HaveOccurred())

			files, err := filepath.Glob(fmt.Sprintf("%s/*", tempDir))
			Expect(err).NotTo(HaveOccurred())
			Expect(files).To(ConsistOf([]string{
				filepath.Join(tempDir, "first"),
				filepath.Join(tempDir, "second"),
				filepath.Join(tempDir, "third"),
				filepath.Join(tempDir, "some-dir"),
			}))

			// info, err := os.Stat(filepath.Join(tempDir, "first"))
			// Expect(err).NotTo(HaveOccurred())
			// Expect(info.Mode()).To(Equal(os.FileMode(0777)))

			Expect(filepath.Join(tempDir, "some-dir", "some-other-dir")).To(BeARegularFile())
			Expect(filepath.Join(tempDir, "some-dir", "some-other-dir", "some-file")).To(BeARegularFile())
		})
	})
}
