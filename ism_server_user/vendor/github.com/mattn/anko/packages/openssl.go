//go:build !appengine
// +build !appengine

package packages

import (
	"reflect"

	"github.com/forgoer/openssl"
	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["openssl"] = map[string]reflect.Value{
		"AesECBDecrypt": reflect.ValueOf(openssl.AesECBDecrypt),
		"AesECBEncrypt": reflect.ValueOf(openssl.AesECBEncrypt),

		"AesCBCEncrypt": reflect.ValueOf(openssl.AesCBCEncrypt),
		"AesCBCDecrypt": reflect.ValueOf(openssl.AesCBCDecrypt),

		"DesECBEncrypt": reflect.ValueOf(openssl.DesECBEncrypt),
		"DesECBDecrypt": reflect.ValueOf(openssl.DesECBDecrypt),

		"DesCBCEncrypt": reflect.ValueOf(openssl.DesCBCEncrypt),
		"DesCBCDecrypt": reflect.ValueOf(openssl.DesCBCDecrypt),

		"Des3ECBEncrypt": reflect.ValueOf(openssl.Des3ECBEncrypt),
		"Des3ECBDecrypt": reflect.ValueOf(openssl.Des3ECBDecrypt),

		"Des3CBCEncrypt": reflect.ValueOf(openssl.Des3CBCEncrypt),
		"Des3CBCDecrypt": reflect.ValueOf(openssl.Des3CBCDecrypt),

		"RSAGenerateKey":       reflect.ValueOf(openssl.RSAGenerateKey),
		"RSAGeneratePublicKey": reflect.ValueOf(openssl.RSAGeneratePublicKey),

		"RSAEncrypt": reflect.ValueOf(openssl.RSAEncrypt),
		"RSADecrypt": reflect.ValueOf(openssl.RSADecrypt),

		"RSASign":   reflect.ValueOf(openssl.RSASign),
		"RSAVerify": reflect.ValueOf(openssl.RSAVerify),

		"PKCS7_PADDING": reflect.ValueOf(openssl.PKCS7_PADDING),
		"PKCS5_PADDING": reflect.ValueOf(openssl.PKCS5_PADDING),
		"ZEROS_PADDING": reflect.ValueOf(openssl.ZEROS_PADDING),
	}
	env.PackageTypes["openssl"] = map[string]reflect.Type{}
}
