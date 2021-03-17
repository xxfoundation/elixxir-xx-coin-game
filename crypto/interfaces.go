package crypto

type Crypto interface {
	RandomGeneration(message string)[]byte
	Weight(digest []byte)uint
}
