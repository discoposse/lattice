package command_factory_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/cloudfoundry-incubator/lattice/ltc/task_runner/command_factory"
	"github.com/cloudfoundry-incubator/lattice/ltc/task_runner/fake_task_runner"
	"github.com/cloudfoundry-incubator/lattice/ltc/terminal"
	"github.com/cloudfoundry-incubator/lattice/ltc/terminal/colors"
	"github.com/cloudfoundry-incubator/lattice/ltc/test_helpers"
	"github.com/codegangsta/cli"
)

var _ = Describe("CommandFactory", func() {

	var (
		outputBuffer   *gbytes.Buffer
		terminalUI     terminal.UI
		fakeTaskRunner *fake_task_runner.FakeTaskRunner
	)

	BeforeEach(func() {
		outputBuffer = gbytes.NewBuffer()
		terminalUI = terminal.NewUI(nil, outputBuffer, nil)
		fakeTaskRunner = new(fake_task_runner.FakeTaskRunner)
	})

	Describe("SubmitTask", func() {

		var (
			submitTaskCommand cli.Command
			tmpDir            string
			tmpFile           *os.File
			err               error
		)

		BeforeEach(func() {
			commandFactory := command_factory.NewTaskRunnerCommandFactory(fakeTaskRunner, terminalUI)
			submitTaskCommand = commandFactory.MakeSubmitTaskCommand()
		})

		Context("when the json file exists", func() {
			BeforeEach(func() {
				tmpDir = os.TempDir()
				tmpFile, err = ioutil.TempFile(tmpDir, "tmp_json")

				Expect(err).ToNot(HaveOccurred())
			})

			It("submits a task from json", func() {
				jsonContents := []byte(`{"Value":"test value"}`)
				ioutil.WriteFile(tmpFile.Name(), jsonContents, 0700)
				args := []string{tmpFile.Name()}
				fakeTaskRunner.SubmitTaskReturns("some-task", nil)

				test_helpers.ExecuteCommandWithArgs(submitTaskCommand, args)

				Expect(outputBuffer).To(test_helpers.Say(colors.Green("Successfully submitted some-task")))
				Expect(fakeTaskRunner.SubmitTaskCallCount()).To(Equal(1))
				Expect(fakeTaskRunner.SubmitTaskArgsForCall(0)).To(Equal(jsonContents))
			})

			It("prints an error returned by the task_runner", func() {
				jsonContents := []byte(`{"Value":"test value"}`)
				ioutil.WriteFile(tmpFile.Name(), jsonContents, 0700)
				args := []string{tmpFile.Name()}
				fakeTaskRunner.SubmitTaskReturns("some-task", errors.New("taskypoo"))

				test_helpers.ExecuteCommandWithArgs(submitTaskCommand, args)

				Expect(fakeTaskRunner.SubmitTaskCallCount()).To(Equal(1))
				Expect(fakeTaskRunner.SubmitTaskArgsForCall(0)).To(Equal(jsonContents))

				Expect(outputBuffer).To(test_helpers.Say("Error submitting some-task: taskypoo"))
			})

		})

		It("is an error when no path is passed in", func() {
			test_helpers.ExecuteCommandWithArgs(submitTaskCommand, []string{})

			Expect(outputBuffer).To(test_helpers.Say("Path to JSON is required"))
			Expect(fakeTaskRunner.SubmitTaskCallCount()).To(BeZero())
		})

		Context("when the file cannot be read", func() {
			It("prints an error", func() {
				args := []string{filepath.Join(tmpDir, "file-no-existy")}

				test_helpers.ExecuteCommandWithArgs(submitTaskCommand, args)

				Expect(outputBuffer).To(test_helpers.Say(fmt.Sprintf("Error reading file: open %s: no such file or directory", filepath.Join(tmpDir, "file-no-existy"))))
				Expect(fakeTaskRunner.SubmitTaskCallCount()).To(Equal(0))
			})
		})

	})

})
