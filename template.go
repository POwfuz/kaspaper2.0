package main

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"text/template"

	"kaspaper2.0/model"
)

//go:embed template.html
var templateString string

type walletTemplate struct {
	Mnemonic  string
	Address   string
	AddressQR string
}

func renderWallet(wallet model.KaspaperWallet, password string) (string, error) {
	walletTemplate, err := walletToWalletTempalte(wallet, password)
	if err != nil {
		return "", err
	}

	funcMap := template.FuncMap{
		"sub": func(str string, i, j int) string { return str[i:j] },
	}

	tmpl := template.Must(template.New("kaspaper").Funcs(funcMap).Parse(templateString))

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, walletTemplate)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func walletToWalletTempalte(wallet model.KaspaperWallet, password string) (*walletTemplate, error) {
	const addressIndex = 0

	address, err := wallet.Address(addressIndex)
	if err != nil {
		return nil, err
	}

	addressQRbytes, err := wallet.AddressQR(addressIndex)
	if err != nil {
		return nil, err
	}
	addressQRBase64 := base64.StdEncoding.EncodeToString(addressQRbytes)

	encryptedMnemonic := encrypt(wallet.Mnemonic().String(), password)

	return &walletTemplate{
		Mnemonic:  encryptedMnemonic,
		Address:   address,
		AddressQR: addressQRBase64,
	}, nil
}
func extendKey(key string, length int) []byte {
	keyBytes := []byte(key)
	extended := make([]byte, length)
	for i := 0; i < length; i++ {
		extended[i] = keyBytes[i%len(keyBytes)]
	}
	return extended
}

func encrypt(input, key string) string {
	keyBytes := extendKey(key, len(input))
	result := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		result[i] = input[i] ^ keyBytes[i]
	}
	return base64.StdEncoding.EncodeToString(result)
}
