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

		port string
	)

	BeforeEach(func() {
		port = strconv.Itoa(60061 + GinkgoParallelNode())

		var err error
		runner, err = admin.Runner(port)
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
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
})
