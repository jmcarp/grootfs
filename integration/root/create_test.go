package root_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"code.cloudfoundry.org/grootfs/integration"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Create", func() {
	var (
		imagePath string
		rootUID   int
		rootGID   int
	)

	BeforeEach(func() {
		rootUID = 0
		rootGID = 0

		var err error
		imagePath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Chown(imagePath, rootUID, rootGID)).To(Succeed())
		Expect(os.Chmod(imagePath, 0755)).To(Succeed())

		grootFilePath := path.Join(imagePath, "foo")
		Expect(ioutil.WriteFile(grootFilePath, []byte("hello-world"), 0644)).To(Succeed())
		Expect(os.Chown(grootFilePath, int(GrootUID), int(GrootGID))).To(Succeed())
		rootFilePath := path.Join(imagePath, "bar")
		Expect(ioutil.WriteFile(rootFilePath, []byte("hello-world"), 0644)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(imagePath)).To(Succeed())
	})

	It("keeps the ownership and permissions", func() {
		bundle := integration.CreateBundle(GrootFSBin, StorePath, imagePath, "random-id")

		grootFi, err := os.Stat(path.Join(bundle.RootFSPath(), "foo"))
		Expect(err).NotTo(HaveOccurred())
		Expect(grootFi.Sys().(*syscall.Stat_t).Uid).To(Equal(uint32(GrootUID)))
		Expect(grootFi.Sys().(*syscall.Stat_t).Gid).To(Equal(uint32(GrootGID)))

		rootFi, err := os.Stat(path.Join(bundle.RootFSPath(), "bar"))
		Expect(err).NotTo(HaveOccurred())
		Expect(rootFi.Sys().(*syscall.Stat_t).Uid).To(Equal(uint32(rootUID)))
		Expect(rootFi.Sys().(*syscall.Stat_t).Gid).To(Equal(uint32(rootGID)))
	})

	Context("when mappings are provided", func() {
		// This test is in the root suite not because `grootfs` is run by root, but
		// because we need to write a file as root to test the translation.
		It("should translate the rootfs accordingly", func() {
			cmd := exec.Command(
				GrootFSBin, "--store", StorePath,
				"--log-level", "debug",
				"create", "--image", imagePath,
				"--uid-mapping", fmt.Sprintf("0:%d:1", GrootUID),
				"--uid-mapping", "1:100000:65000",
				"--gid-mapping", fmt.Sprintf("0:%d:1", GrootUID),
				"--gid-mapping", "1:100000:65000",
				"some-id",
			)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Credential: &syscall.Credential{
					Uid: GrootUID,
					Gid: GrootGID,
				},
			}
			sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(sess).Should(gexec.Exit(0))
			bundle := strings.TrimSpace(string(sess.Out.Contents()))

			grootFi, err := os.Stat(path.Join(bundle, "rootfs", "foo"))
			Expect(err).NotTo(HaveOccurred())
			Expect(grootFi.Sys().(*syscall.Stat_t).Uid).To(Equal(uint32(GrootUID + 99999)))
			Expect(grootFi.Sys().(*syscall.Stat_t).Gid).To(Equal(uint32(GrootGID + 99999)))

			rootFi, err := os.Stat(path.Join(bundle, "rootfs", "bar"))
			Expect(err).NotTo(HaveOccurred())
			Expect(rootFi.Sys().(*syscall.Stat_t).Uid).To(Equal(uint32(GrootUID)))
			Expect(rootFi.Sys().(*syscall.Stat_t).Gid).To(Equal(uint32(GrootGID)))
		})
	})
})
