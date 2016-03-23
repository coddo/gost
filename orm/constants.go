package orm

const (
	// NormalUserAccountType represents a ordinary application user
	AccountTypeNormalUser = iota
	// AdministratorAccountType represents an administrator
	AccountTypeAdministrator = iota
)

const (
	// StatusAccountActivated represents an account that is active in the system
	AccountStatusActivated = true
	// StatusAccountDeactivated represents an account that is inactive in the system
	AccountStatusDeactivated = false
)
