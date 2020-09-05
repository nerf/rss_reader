package rss_reader_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRssReader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RssReader Suite")
}

var _ = Describe("RssReader", func() {})
