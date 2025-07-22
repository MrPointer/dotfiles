package dotfilesmanager

type DotfilesDataInitializer interface {
	Initialize(data DotfilesData) error
}

type DotfilesApplier interface {
	Apply() error
}

type DotfilesInstaller interface {
	Install() error
}

type DotfilesManager interface {
	DotfilesDataInitializer
	DotfilesApplier
	DotfilesInstaller
}
