package command

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"testing"

	"github.com/SeerUK/assert"
	"github.com/SeerUK/foldup/pkg/archive"
	"github.com/SeerUK/foldup/pkg/foldup"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

func TestBackupCommand(t *testing.T) {
	log.SetOutput(&bytes.Buffer{})

	t.Run("should return the backup command", func(t *testing.T) {
		factory := foldup.NewCLIFactory()
		backupCmd := BackupCommand(factory)

		assert.Equal(t, "backup", backupCmd.Name)
	})

	t.Run("should prepare the input definition", func(t *testing.T) {
		def := console.NewDefinition()

		factory := foldup.NewCLIFactory()
		backupCmd := BackupCommand(factory)
		backupCmd.Configure(def)

		args := def.Arguments()
		opts := def.Options()

		assert.Equal(t, 1, len(args))
		assert.Equal(t, 2, len(opts))

		assert.Equal(t, "DIRNAME", args[0].Name)
		assert.Equal(t, []string{"b", "bucket"}, opts[0].Names)
		assert.Equal(t, []string{"s", "schedule"}, opts[1].Names)
	})

	t.Run("should error if the storage gateway can't be created", func(t *testing.T) {
		def := console.NewDefinition()

		factory := &backupTestFactory{
			createGCSGatewayError: errors.New("oops"),
		}

		backupCmd := BackupCommand(factory)
		backupCmd.Configure(def)

		setArgValue(def.Arguments(), "DIRNAME", "testdata")
		setOptValue(def.Options(), "bucket", "test-bucket")

		input, output := createInputAndOutput(&bytes.Buffer{})

		result := backupCmd.Execute(input, output)

		assert.NotOK(t, result)
	})

	t.Run("should error if the archiving target doesn't exist", func(t *testing.T) {
		def := console.NewDefinition()

		factory := &backupTestFactory{}
		backupCmd := BackupCommand(factory)
		backupCmd.Configure(def)

		setArgValue(def.Arguments(), "DIRNAME", "i-do-not-exist")
		setOptValue(def.Options(), "bucket", "test-bucket")

		input, output := createInputAndOutput(&bytes.Buffer{})

		result := backupCmd.Execute(input, output)

		assert.NotOK(t, result)
	})

	t.Run("should error if archiving fails", func(t *testing.T) {
		defer revertStubs()

		archiveDirsf = func(ds []string, nf string, fn archive.FormatName) ([]string, error) {
			return []string{}, errors.New("oops")
		}

		def := console.NewDefinition()

		factory := &backupTestFactory{}
		backupCmd := BackupCommand(factory)
		backupCmd.Configure(def)

		setArgValue(def.Arguments(), "DIRNAME", "testdata")
		setOptValue(def.Options(), "bucket", "test-bucket")

		input, output := createInputAndOutput(&bytes.Buffer{})

		result := backupCmd.Execute(input, output)

		assert.NotOK(t, result)
	})

	t.Run("should error if opening the archive to upload fails", func(t *testing.T) {
		defer revertStubs()

		osOpen = func(name string) (*os.File, error) {
			return nil, errors.New("oops")
		}

		def := console.NewDefinition()

		factory := &backupTestFactory{}
		backupCmd := BackupCommand(factory)
		backupCmd.Configure(def)

		setArgValue(def.Arguments(), "DIRNAME", "testdata")
		setOptValue(def.Options(), "bucket", "test-bucket")

		input, output := createInputAndOutput(&bytes.Buffer{})

		result := backupCmd.Execute(input, output)

		assert.NotOK(t, result)
	})

	t.Run("should error if uploading the archive fails", func(t *testing.T) {
		def := console.NewDefinition()

		gateway := &backupTestStorageGateway{
			storeError: errors.New("oops"),
		}

		factory := &backupTestFactory{
			createGCSGatewayGateway: gateway,
		}

		backupCmd := BackupCommand(factory)
		backupCmd.Configure(def)

		setArgValue(def.Arguments(), "DIRNAME", "testdata")
		setOptValue(def.Options(), "bucket", "test-bucket")

		input, output := createInputAndOutput(&bytes.Buffer{})

		result := backupCmd.Execute(input, output)

		assert.NotOK(t, result)
	})

	t.Run("should error if removing the archive fails", func(t *testing.T) {
		defer revertStubs()

		osRemove = func(name string) error {
			return errors.New("oops")
		}

		def := console.NewDefinition()

		factory := &backupTestFactory{}
		backupCmd := BackupCommand(factory)
		backupCmd.Configure(def)

		setArgValue(def.Arguments(), "DIRNAME", "testdata")
		setOptValue(def.Options(), "bucket", "test-bucket")

		input, output := createInputAndOutput(&bytes.Buffer{})

		result := backupCmd.Execute(input, output)

		assert.NotOK(t, result)
	})

	t.Run("should perform a backup", func(t *testing.T) {
		def := console.NewDefinition()

		factory := &backupTestFactory{}
		backupCmd := BackupCommand(factory)
		backupCmd.Configure(def)

		setArgValue(def.Arguments(), "DIRNAME", "testdata")
		setOptValue(def.Options(), "bucket", "test-bucket")

		input, output := createInputAndOutput(&bytes.Buffer{})

		result := backupCmd.Execute(input, output)

		assert.OK(t, result)
	})

	t.Run("should be able to schedule a backup", func(t *testing.T) {
		def := console.NewDefinition()

		scheduleFunc = func(done <-chan int, expr string, fn func() error) error {
			return fn()
		}

		factory := &backupTestFactory{}
		backupCmd := BackupCommand(factory)
		backupCmd.Configure(def)

		setArgValue(def.Arguments(), "DIRNAME", "testdata")
		setOptValue(def.Options(), "bucket", "test-bucket")
		setOptValue(def.Options(), "schedule", "* * * * * * *")

		input, output := createInputAndOutput(&bytes.Buffer{})

		result := backupCmd.Execute(input, output)

		assert.OK(t, result)
	})

	log.SetOutput(os.Stdout)
}

func createInputAndOutput(writer io.Writer) (*console.Input, *console.Output) {
	return &console.Input{}, console.NewOutput(writer)
}

func setArgValue(args []parameters.Argument, name string, value string) error {
	for _, arg := range args {
		if arg.Name == name {
			return arg.Value.Set(value)
		}
	}

	return nil
}

func setOptValue(opts []parameters.Option, name string, value string) error {
	for _, opt := range opts {
		for _, n := range opt.Names {
			if n == name {
				return opt.Value.Set(value)
			}
		}
	}

	return nil
}
