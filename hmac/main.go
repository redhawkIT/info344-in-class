package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

const usage = `
usage:
	hmac sign|verify <key> <value>
`

func main() {
	if len(os.Args) < 4 ||
		(os.Args[1] != "sign" && os.Args[1] != "verify") {
		fmt.Println(usage)
		os.Exit(1)
	}

	cmd := os.Args[1]
	key := os.Args[2]
	value := os.Args[3]

	switch cmd {
	case "sign":
		v := []byte(value)
		h := hmac.New(sha256.New, []byte(key))
		h.Write(v)
		sig := h.Sum(nil)

		buf := make([]byte, len(v)+len(sig))
		copy(buf, v)
		//copy sig into last part of buf
		copy(buf[len(v):], sig)
		fmt.Println(base64.URLEncoding.EncodeToString(buf))

	case "verify":
		buf, err := base64.URLEncoding.DecodeString(value)
		if err != nil {
			fmt.Printf("error decoding: %v\n", err)
			os.Exit(1)
		}
		v := buf[:len(buf)-sha256.Size]
		sig := buf[len(buf)-sha256.Size:]

		h := hmac.New(sha256.New, []byte(key))
		h.Write(v)
		sig2 := h.Sum(nil)
		if hmac.Equal(sig, sig2) {
			fmt.Println("signature is valid!")
		} else {
			fmt.Println("DANGER DANGER INVALID!!")
		}
	}

}
