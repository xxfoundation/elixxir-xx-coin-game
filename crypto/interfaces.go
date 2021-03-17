package crypto

type Crypto interface {
	RandomGeneration(message string) []byte
	Weight(digest []byte) uint
}

// todo: this will be read as config
var Salt = []byte{113, 2, 213, 85, 51, 97, 1, 238, 213, 61, 19, 249, 179, 241, 184, 198, 228, 73, 125, 64, 163, 94, 117, 141, 172, 173, 41, 249, 98, 89, 53, 53}
