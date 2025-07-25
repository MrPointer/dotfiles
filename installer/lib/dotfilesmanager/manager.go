package dotfilesmanager

// DotfilesDataInitializer initializes the dotfiles data.
type DotfilesDataInitializer interface {
	// Initialize initializes the dotfiles data.
	// It takes a DotfilesData object as input and returns an error if any.
	//
	// data: The DotfilesData object to initialize.
	Initialize(data DotfilesData) error
}

// DotfilesApplier applies the dotfiles.
type DotfilesApplier interface {
	// Apply applies the dotfiles.
	// It returns an error if any.
	Apply() error
}

// DotfilesInstaller installs the dotfiles.
type DotfilesInstaller interface {
	// Install installs the dotfiles.
	// It returns an error if any.
	Install() error
}

// DotfilesManager manages the dotfiles, by providing a unified interface for initializing, applying, and installing dotfiles.
type DotfilesManager interface {
	DotfilesDataInitializer
	DotfilesApplier
	DotfilesInstaller
}
