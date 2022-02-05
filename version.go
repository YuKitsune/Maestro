package maestro

var (
	// Version defines current version of Maestro
	// This is populated via ldflags
	Version string
)

func init() {
	// If version, commit, or build time are not set, make that clear.
	const unknown = "unknown"
	if Version == "" {
		Version = unknown
	}
}
