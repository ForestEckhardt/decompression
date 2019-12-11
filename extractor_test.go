package decompression_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ForestEckhardt/decompression"
	"github.com/ForestEckhardt/decompression/fakes"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testExtractor(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		reader *bytes.Reader
	)

	it.Before(func() {
		var err error
		reader = bytes.NewReader(nil)
		Expect(err).ToNot(HaveOccurred())
	})

	context("NewExtractor", func() {
		context("when the archive format is tar", func() {
			it("returns a TarExtractor object", func() {
				extractor, err := decompression.NewDecompressor(reader, decompression.Tar)
				Expect(err).ToNot(HaveOccurred())
				Expect(extractor).To(Equal(decompression.Decompressor{Reader: decompression.NewArchiveReader(reader)}))
			})
		})

		context("when the archive format is gzip tar", func() {
			var gzipReader *bytes.Reader

			it.Before(func() {
				var err error
				buffer := bytes.NewBuffer(nil)
				gzipWriter := gzip.NewWriter(buffer)
				tw := tar.NewWriter(gzipWriter)

				Expect(tw.WriteHeader(&tar.Header{Name: "some-file", Mode: 0755, Size: int64(len("some-file"))})).To(Succeed())
				_, err = tw.Write([]byte("some-file"))
				Expect(err).NotTo(HaveOccurred())

				Expect(gzipWriter.Close()).To(Succeed())
				Expect(tw.Close()).To(Succeed())

				gzipReader = bytes.NewReader(buffer.Bytes())
			})
			it("returns a GZTarExtractor object", func() {
				extractor, err := decompression.NewDecompressor(gzipReader, decompression.GZTar)
				Expect(err).ToNot(HaveOccurred())

				bufFinal := bytes.NewBuffer(nil)
				_, err = io.Copy(bufFinal, extractor.Reader.GetReader())
				Expect(err).ToNot(HaveOccurred())

				_, err = gzipReader.Seek(0, 0)
				Expect(err).ToNot(HaveOccurred())

				gzipResults, err := gzip.NewReader(gzipReader)
				Expect(err).ToNot(HaveOccurred())

				bufCompare := bytes.NewBuffer(nil)

				_, err = io.Copy(bufCompare, gzipResults)
				Expect(err).ToNot(HaveOccurred())

				Expect(bufFinal.Bytes()).To(Equal(bufCompare.Bytes()))
			})

			context("failure case", func() {
				it("returns an error", func() {
					_, err := decompression.NewDecompressor(bytes.NewBuffer([]byte(`something`)), decompression.GZTar)
					Expect(err).To(MatchError(ContainSubstring("failed to create gzip reader")))
				})
			})
		})
	})

	context("Decompress", func() {
		var (
			extractor     decompression.Decompressor
			decompressTar *fakes.DecompressTar
			tempDir       string
		)

		it.Before(func() {
			var err error
			decompressTar = &fakes.DecompressTar{}
			extractor = decompression.Decompressor{Reader: decompressTar}
			tempDir, err = ioutil.TempDir("", "decompression")
			Expect(err).ToNot(HaveOccurred())
		})

		it.After(func() {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		})

		it("extracts files from the reader", func() {
			err := extractor.Decompress(tempDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(decompressTar.UnTarCall.Receives.Destination).To(Equal(tempDir))
		})

		context("failure case", func() {
			context("when the untar fails", func() {
				it.Before(func() {
					decompressTar.UnTarCall.Returns.Error = errors.New("failed to untar")
				})
				it("throws an error", func() {
					err := extractor.Decompress(tempDir)
					Expect(err).To(MatchError(ContainSubstring("failed to untar")))
				})
			})
		})
	})
}
