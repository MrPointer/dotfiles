package shell

// ShellSource specifies where to find the shell binary.
type ShellSource string

const (
	ShellSourceAuto   ShellSource = "auto"
	ShellSourceBrew   ShellSource = "brew"
	ShellSourceSystem ShellSource = "system"
)

// String returns the string representation of the ShellSource.
func (s ShellSource) String() string {
	return string(s)
}
