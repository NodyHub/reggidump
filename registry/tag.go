package registry

import (
	"fmt"
	"strings"
)

type Tag struct {
	Name     string
	Manifest *ManifestV1
}

func (t *Tag) Starter() string {

	// compose the starter command
	starter := t.Entrypoint()
	if len(starter) > 0 {
		starter = starter + " "
	}
	starter = starter + t.Command()

	return starter
}

func (t *Tag) Command() string {
	if t.Manifest == nil {
		return ""
	}
	if len(t.Manifest.History) == 0 {
		return ""
	}

	// compose the command
	cmd := t.Manifest.History[0].HistoryEntry.Config.Cmd
	command := strings.Join(cmd, " ")

	return command
}

func (t *Tag) Entrypoint() string {
	if t.Manifest == nil {
		return ""
	}
	if len(t.Manifest.History) == 0 {
		return ""
	}

	// compose the entrypoint
	entry := t.Manifest.History[0].HistoryEntry.Config.Entrypoint
	entrypoint := strings.Join(entry, " ")

	return entrypoint
}

func (t *Tag) GeneratePatchScript(script string) string {

	// get existing entrypoint and command
	entryPoint := t.Entrypoint()
	command := t.Command()

	// compose the new starter command
	switch {

	case len(entryPoint) > 0 && len(command) == 0: // only entrypoint
		return fmt.Sprintf(ENTRYPOINT_TEMPLATE, script, entryPoint)

	case len(entryPoint) == 0 && len(command) > 0: // only command
		return fmt.Sprintf(COMMAND_TEMPLATE, script, command)

	case len(entryPoint) > 0 && len(command) > 0: // both entrypoint and command
		return fmt.Sprintf(ENTRYPOINT_COMMAND_TEMPLATE, script, entryPoint, command)

	default: // nothing found
		return fmt.Sprintf(EMPTY_STARTER, script)
	}
}

var ENTRYPOINT_TEMPLATE = `#!/bin/sh
set -e

# --------------------
# Patch script
# --------------------
%s
# --------------------

# Old entrypoint
exec %s "$@"
`

var COMMAND_TEMPLATE = `#!/bin/sh
set -e

# --------------------
# Patch script
# --------------------
%s
# --------------------

# Old command
exec %s
`

var ENTRYPOINT_COMMAND_TEMPLATE = `#!/bin/sh
set -e

# --------------------
# Patch script
# --------------------
%s
# --------------------

# Old entrypoint and command
exec %s %s
`

var EMPTY_STARTER = `#!/bin/sh
set -e

# --------------------
# Patch script
# --------------------
%s
# --------------------

# no further commands found
`
