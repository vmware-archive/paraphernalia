package admin_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/pivotal-cf/paraphanerlia/operate/admin"
	"github.com/tedsuo/ifrit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Runner", func() {
	var (
		runner  ifrit.Runner
		process ifrit.Process

		port        string
		optionFuncs []admin.OptionFunc
	)

	BeforeEach(func() {
		port = strconv.Itoa(60061 + GinkgoParallelNode())
	})

	JustBeforeEach(func() {
		var err error
		runner, err = admin.Runner(port, optionFuncs...)
		Expect(err).NotTo(HaveOccurred())

		process = ifrit.Invoke(runner)
	})

	AfterEach(func() {
		process.Signal(os.Interrupt)

		err := <-process.Wait()
		Expect(err).NotTo(HaveOccurred())
	})

	It("starts the debug server", func() {
		response, err := http.Get("http://localhost:" + port + "/debug/pprof/cmdline")
		Expect(err).NotTo(HaveOccurred())

		body, err := ioutil.ReadAll(response.Body)
		Expect(err).NotTo(HaveOccurred())

		// This is the binary name for this packages test suite.
		Expect(body).To(ContainSubstring("admin.test"))
	})

	Describe("enabling the information endpoint", func() {
		BeforeEach(func() {
			optionFuncs = []admin.OptionFunc{
				admin.WithInfo(admin.ServiceInfo{
					Name:        "service-name",
					Description: "it's a thing which does a thing",
					Team:        "team name",
				}),
			}
		})

		It("let's the developers of a service tell people what it is", func() {
			response, err := http.Get("http://localhost:" + port + "/info")
			Expect(err).NotTo(HaveOccurred())

			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(body).To(ContainSubstring("service-name"))
			Expect(body).To(ContainSubstring("it's a thing which does a thing"))
			Expect(body).To(ContainSubstring("team name"))
		})
	})
})
