package overlayxfs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOverlayxfs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Overlay+Xfs Driver Suite")
}
