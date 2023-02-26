package utils

import (
	"encoding/hex"
	"hash/fnv"
	"regexp"
)

const (
	invalidHashError = "'%s' is not a valid short path."
	invalidURLError  = "'%s' is not a valid URL."
)

var (
	short        = regexp.MustCompile(`[a-zA-Z0-9]{8}`)
	long         = regexp.MustCompile(`https?://(?:[-\w.]|%[\da-fA-F]{2})+`)
	sessionToken = "session_token"
)

func isShortCorrect(link string) bool {
	return short.FindStringIndex(link) != nil
}

func isLongCorrect(link string) bool {
	return long.FindStringIndex(link) != nil
}

func getHash(s []byte) (string, error) {
	hasher := fnv.New32a()
	_, err := hasher.Write(s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
