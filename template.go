package main

import (
	"bytes"
	"crypto/rc4"
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

	encryptedMnemonic, err := rc4Encrypt(wallet.Mnemonic().String(), password)
	if err != nil {
		return nil, err
	}

	return &walletTemplate{
		Mnemonic:  encryptedMnemonic,
		Address:   address,
		AddressQR: addressQRBase64,
	}, nil
}

func rc4Encrypt(data, key string) (string, error) {
	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plainText := []byte(data)
	cipherText := make([]byte, len(plainText))
	cipher.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
