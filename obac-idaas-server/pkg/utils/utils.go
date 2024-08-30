package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(n int) (string, error) {

	b, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func RemoveDuplicates(slice *[]string) {
	encountered := map[string]bool{}
	// Iterate over the slice and remove duplicates in place
	for i := 0; i < len(*slice); i++ {
		if encountered[(*slice)[i]] {
			// Duplicate element found, remove it by shifting remaining elements left
			(*slice)[i] = (*slice)[len(*slice)-1]
			*slice = (*slice)[:len(*slice)-1]
			i-- // Revisit this index since it now contains a new element
		} else {
			// Unique element found, mark it as encountered
			encountered[(*slice)[i]] = true
		}
	}
}

func EmptyChannel(ch *chan bool) {
	for len(*ch) > 0 {
		<-*ch
	}
}
