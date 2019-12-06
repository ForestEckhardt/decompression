package decompression_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitNode(t *testing.T) {
	suite := spec.New("node", spec.Report(report.Terminal{}))
	suite("Extractor", testExtractor)
	suite("UnTar", testUnTar)
	suite.Run(t)
}
