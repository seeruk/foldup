package command

import (
	"testing"

	"github.com/SeerUK/assert"
	"github.com/SeerUK/foldup/pkg/foldup"
	"github.com/eidolon/console"
)

func TestBackupCommand(t *testing.T) {
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

	t.Run("should ", func(t *testing.T) {

	})
}
