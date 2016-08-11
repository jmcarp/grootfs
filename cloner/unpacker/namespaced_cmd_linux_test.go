package unpacker_test

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"code.cloudfoundry.org/grootfs/cloner"
	"code.cloudfoundry.org/grootfs/cloner/unpacker"
	"code.cloudfoundry.org/grootfs/cloner/unpacker/unpackerfakes"
	"code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/cloudfoundry/gunk/command_runner/fake_command_runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NamespacedCmdUnpacker", func() {
	var (
		fakeIDMapper      *unpackerfakes.FakeIDMapper
		fakeCommandRunner *fake_command_runner.FakeCommandRunner
		tarUnpacker       *unpacker.NamespacedCmdUnpacker

		logger     lager.Logger
		bundlePath string
		rootFSPath string

		commandError error
	)

	BeforeEach(func() {
		var err error

		fakeIDMapper = new(unpackerfakes.FakeIDMapper)
		fakeCommandRunner = fake_command_runner.New()
		tarUnpacker = unpacker.NewNamespacedCmdUnpacker(
			fakeCommandRunner, fakeIDMapper, "ginkgo-unpack",
		)

		logger = lagertest.NewTestLogger("test-store")

		bundlePath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		rootFSPath = filepath.Join(bundlePath, "rootfs")

		commandError = nil
	})

	JustBeforeEach(func() {
		fakeCommandRunner.WhenRunning(fake_command_runner.CommandSpec{
			Path: os.Args[0],
		}, func(cmd *exec.Cmd) error {
			cmd.Process = &os.Process{
				Pid: 12, // don't panic
			}
			return commandError
		})
	})

	AfterEach(func() {
		Expect(os.RemoveAll(bundlePath)).To(Succeed())
	})

	It("passes the rootfs path to the provided command", func() {
		Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
			RootFSPath: rootFSPath,
		})).To(Succeed())

		commands := fakeCommandRunner.StartedCommands()
		Expect(commands).To(HaveLen(1))
		expectedPath := os.Args[0]
		Expect(commands[0].Path).To(Equal(expectedPath))
		Expect(commands[0].Args).To(Equal([]string{
			expectedPath, "ginkgo-unpack", rootFSPath,
		}))
	})

	It("uses the provided stream", func() {
		streamR, streamW, err := os.Pipe()
		Expect(err).NotTo(HaveOccurred())

		Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
			Stream:     streamR,
			RootFSPath: rootFSPath,
		})).To(Succeed())

		commands := fakeCommandRunner.StartedCommands()
		Expect(commands).To(HaveLen(1))

		_, err = streamW.WriteString("hello-world")
		Expect(err).NotTo(HaveOccurred())
		Expect(streamW.Close()).To(Succeed())

		contents, err := ioutil.ReadAll(commands[0].Stdin)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(Equal("hello-world"))
	})

	It("starts the provided command in a user namespace", func() {
		Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
			UIDMappings: []groot.IDMappingSpec{
				groot.IDMappingSpec{HostID: 1000, NamespaceID: 2000, Size: 10},
			},
			RootFSPath: rootFSPath,
		})).To(Succeed())

		commands := fakeCommandRunner.StartedCommands()
		Expect(commands).To(HaveLen(1))
		Expect(commands[0].SysProcAttr.Cloneflags).To(Equal(uintptr(syscall.CLONE_NEWUSER)))
	})

	Context("when no mappings are provided", func() {
		It("starts the provided command in the same namespaces", func() {
			Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
				RootFSPath: rootFSPath,
			})).To(Succeed())

			commands := fakeCommandRunner.StartedCommands()
			Expect(commands).To(HaveLen(1))
			Expect(commands[0].SysProcAttr).To(BeNil())
		})
	})

	It("signals the namespaced command to continue using the contol pipe", func(done Done) {
		Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
			RootFSPath: rootFSPath,
		})).To(Succeed())

		commands := fakeCommandRunner.StartedCommands()
		Expect(commands).To(HaveLen(1))
		buffer := make([]byte, 1)
		_, err := commands[0].ExtraFiles[0].Read(buffer)
		Expect(err).NotTo(HaveOccurred())

		close(done)
	}, 1.0)

	Describe("UIDMappings", func() {
		It("uses the provided uid mapping", func() {
			Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
				RootFSPath: rootFSPath,
				UIDMappings: []groot.IDMappingSpec{
					groot.IDMappingSpec{HostID: 1000, NamespaceID: 2000, Size: 10},
				},
			})).To(Succeed())

			Expect(fakeIDMapper.MapUIDsCallCount()).To(Equal(1))
			_, _, mappings := fakeIDMapper.MapUIDsArgsForCall(0)

			Expect(mappings).To(Equal([]groot.IDMappingSpec{
				groot.IDMappingSpec{HostID: 1000, NamespaceID: 2000, Size: 10},
			}))
		})

		Context("when mapping fails", func() {
			BeforeEach(func() {
				fakeIDMapper.MapUIDsReturns(errors.New("Boom!"))
			})

			It("returns an error", func() {
				Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
					RootFSPath: rootFSPath,
					UIDMappings: []groot.IDMappingSpec{
						groot.IDMappingSpec{HostID: 1000, NamespaceID: 2000, Size: 10},
					},
				})).To(MatchError(ContainSubstring("Boom!")))
			})

			It("closes the control pipe", func() {
				Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
					RootFSPath: rootFSPath,
					UIDMappings: []groot.IDMappingSpec{
						groot.IDMappingSpec{HostID: 1000, NamespaceID: 2000, Size: 10},
					},
				})).NotTo(Succeed())

				commands := fakeCommandRunner.StartedCommands()
				Expect(commands).To(HaveLen(1))
				buffer := make([]byte, 1)
				_, err := commands[0].ExtraFiles[0].Read(buffer)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GIDMappings", func() {
		It("uses the provided gid mapping", func() {
			Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
				RootFSPath: rootFSPath,
				GIDMappings: []groot.IDMappingSpec{
					groot.IDMappingSpec{HostID: 1000, NamespaceID: 2000, Size: 10},
				},
			})).To(Succeed())

			Expect(fakeIDMapper.MapGIDsCallCount()).To(Equal(1))
			_, _, mappings := fakeIDMapper.MapGIDsArgsForCall(0)

			Expect(mappings).To(Equal([]groot.IDMappingSpec{
				groot.IDMappingSpec{HostID: 1000, NamespaceID: 2000, Size: 10},
			}))
		})

		Context("when mapping fails", func() {
			BeforeEach(func() {
				fakeIDMapper.MapGIDsReturns(errors.New("Boom!"))
			})

			It("returns an error", func() {
				Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
					RootFSPath: rootFSPath,
					GIDMappings: []groot.IDMappingSpec{
						groot.IDMappingSpec{HostID: 1000, NamespaceID: 2000, Size: 10},
					},
				})).To(MatchError(ContainSubstring("Boom!")))
			})

			It("closes the control pipe", func() {
				Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
					RootFSPath: rootFSPath,
					GIDMappings: []groot.IDMappingSpec{
						groot.IDMappingSpec{HostID: 1000, NamespaceID: 2000, Size: 10},
					},
				})).NotTo(Succeed())

				commands := fakeCommandRunner.StartedCommands()
				Expect(commands).To(HaveLen(1))
				buffer := make([]byte, 1)
				_, err := commands[0].ExtraFiles[0].Read(buffer)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when it fails to start unpacking", func() {
		BeforeEach(func() {
			commandError = errors.New("failed to start unpack")
		})

		It("returns an error", func() {
			Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
				RootFSPath: rootFSPath,
			})).To(
				MatchError(ContainSubstring("failed to start unpack")),
			)
		})
	})

	Context("when it fails to unpack", func() {
		BeforeEach(func() {
			fakeCommandRunner.WhenWaitingFor(fake_command_runner.CommandSpec{
				Path: os.Args[0],
			}, func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte("hello-world"))
				return errors.New("exit status 1")
			})
		})

		It("returns an error", func() {
			Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
				RootFSPath: rootFSPath,
			})).NotTo(Succeed())
		})

		It("returns the command output", func() {
			Expect(tarUnpacker.Unpack(logger, cloner.UnpackSpec{
				RootFSPath: rootFSPath,
			})).To(
				MatchError(ContainSubstring("hello-world")),
			)
		})
	})
})
