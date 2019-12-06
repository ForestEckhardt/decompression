package decompression_test

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ForestEckhardt/decompression"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testUnTar(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
	)

	context("UnTar", func() {
		var (
			tempDir string
			reader  decompression.ArchiveReader
		)

		it.Before(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "decompression")
			Expect(err).NotTo(HaveOccurred())

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

			reader = decompression.NewArchiveReader(bytes.NewReader(buffer.Bytes()))

			Expect(tw.Close()).To(Succeed())

		})

		it.After(func() {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		})

		it("downloads the dependency and unpackages it into the path", func() {
			var err error
			err = reader.UnTar(tempDir)
			Expect(err).ToNot(HaveOccurred())

			files, err := filepath.Glob(fmt.Sprintf("%s/*", tempDir))
			Expect(err).NotTo(HaveOccurred())
			Expect(files).To(ConsistOf([]string{
				filepath.Join(tempDir, "first"),
				filepath.Join(tempDir, "second"),
				filepath.Join(tempDir, "third"),
				filepath.Join(tempDir, "some-dir"),
			}))

			info, err := os.Stat(filepath.Join(tempDir, "first"))
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Mode()).To(Equal(os.FileMode(0755)))

			Expect(filepath.Join(tempDir, "some-dir", "some-other-dir")).To(BeADirectory())
			Expect(filepath.Join(tempDir, "some-dir", "some-other-dir", "some-file")).To(BeARegularFile())
		})

		context("failure cases", func() {
			context("when it fails to read the tar response", func() {
				it("returns an error", func() {
					err := decompression.NewArchiveReader(bytes.NewBuffer([]byte(`something`))).UnTar(tempDir)
					Expect(err).To(MatchError(ContainSubstring("failed to read tar response")))

				})
			})

			context("when it is unable to create an archived directory", func() {
				it.Before(func() {
					Expect(os.Chmod(tempDir, 0000)).To(Succeed())
				})

				it.After(func() {
					Expect(os.Chmod(tempDir, os.ModePerm)).To(Succeed())
				})

				it("returns an error", func() {
					err := reader.UnTar(tempDir)
					Expect(err).To(MatchError(ContainSubstring("failed to create archived directory")))
				})
			})

			context("when it is unable to create an archived file", func() {
				it.Before(func() {
					Expect(os.MkdirAll(filepath.Join(tempDir, "some-dir", "some-other-dir"), os.ModePerm)).To(Succeed())
					Expect(os.Chmod(filepath.Join(tempDir, "some-dir", "some-other-dir"), 0000)).To(Succeed())
				})

				it.After(func() {
					Expect(os.Chmod(filepath.Join(tempDir, "some-dir", "some-other-dir"), os.ModePerm)).To(Succeed())
				})

				it("returns an error", func() {
					err := reader.UnTar(tempDir)
					Expect(err).To(MatchError(ContainSubstring("failed to create archived file")))
				})
			})
		})
	})
}
