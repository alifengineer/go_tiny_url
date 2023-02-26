package utils

import (
	"encoding/hex"
	"hash/fnv"
	"regexp"
)

const (
	InvalidHashError = "'%s' is not a valid short path."
	InvalidURLError  = "'%s' is not a valid URL."
	InvalidUUIDError = "'%s' is not a valid UUID."
)

var (
	short        = regexp.MustCompile(`[a-zA-Z0-9]{8}`)
	long         = regexp.MustCompile(`https?://(?:[-\w.]|%[\da-fA-F]{2})+`)
	sessionToken = "session_token"
)

func IsShortCorrect(link string) bool {
	return short.FindStringIndex(link) != nil
}

func IsLongCorrect(link string) bool {
	return long.FindStringIndex(link) != nil
}

func GetHash(s []byte) (string, error) {
	hasher := fnv.New32a()
	_, err := hasher.Write(s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
