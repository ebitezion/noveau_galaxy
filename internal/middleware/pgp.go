package middleware

// pgp.go

import "github.com/ProtonMail/gopenpgp/v2/helper"

func EncryptMessageArmored(publicKey, message string) (string, error) {
	return helper.EncryptMessageArmored(publicKey, message)
}

func DecryptMessageArmored(privateKey string, passphrase []byte, armoredMessage string) (string, error) {
	return helper.DecryptMessageArmored(privateKey, passphrase, armoredMessage)
}

func EncryptBinaryMessageArmored(publicKey string, message []byte) (string, error) {
	return helper.EncryptBinaryMessageArmored(publicKey, message)
}

func DecryptBinaryMessageArmored(privateKey string, passphrase []byte, armoredMessage string) ([]byte, error) {
	return helper.DecryptBinaryMessageArmored(privateKey, passphrase, armoredMessage)
}
