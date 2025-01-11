package crypto

type Crypto interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) bool
}
