package store

// Store ...
type Store interface {
	User() UserRepository
	UserAccount() UserAccountRepository
	UserPackage() UserPackageRepository
}
