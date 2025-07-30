package domain

import (
	"log"
	"testing"
)

func TestHackSafeCode(t *testing.T) {
	code := GenSafeCode(8)
	secret := GenSafeSecret(8)

	hacked := HackSafeCode(code, secret)
	hacked = HackSafeCode(code, hacked)
	hacked = HackSafeCode(code, hacked)
	hacked = HackSafeCode(code, hacked)

	log.Printf("code: %s, secret: %s, hacked: %s", code, secret, hacked)
}
