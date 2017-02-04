// Package env creates and updates the env.json file.
package env

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// UpdateFileKeys updates the session keys in the env.json.
func UpdateFileKeys(src string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}

	var newFile string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		if len(newFile) > 0 {
			newFile += "\n"
		}

		if strings.Contains(scanner.Text(), `"AuthKey"`) {
			newFile += fmt.Sprintf(`    "AuthKey":"%v",`, EncodedKey(64))
		} else if strings.Contains(scanner.Text(), `"EncryptKey"`) {
			newFile += fmt.Sprintf(`    "EncryptKey":"%v",`, EncodedKey(32))
		} else if strings.Contains(scanner.Text(), `"CSRFKey"`) {
			newFile += fmt.Sprintf(`    "CSRFKey":"%v",`, EncodedKey(32))
		} else {
			newFile += scanner.Text()
		}
	}

	if err := scanner.Err(); err != nil {
		file.Close()
		return err
	}

	file.Close()
	return ioutil.WriteFile(src, []byte(newFile), 0644)
}

// EncodedKey returns a base64 encoded random key.
func EncodedKey(length int) string {
	k := make([]byte, length)
	rand.Read(k)
	return base64.StdEncoding.EncodeToString(k)
}
