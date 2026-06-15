package packages

import (
	"crypto/cipher"
	"reflect"

	"github.com/mattn/anko/env"
)

func init() {
	env.Packages["crypto/cipher"] = map[string]reflect.Value{
		"NewGCM":              reflect.ValueOf(cipher.NewGCM),
		"NewGCMWithNonceSize": reflect.ValueOf(cipher.NewGCMWithNonceSize),
		"NewGCMWithTagSize":   reflect.ValueOf(cipher.NewGCMWithTagSize),
		"NewCBCDecrypter":     reflect.ValueOf(cipher.NewCBCDecrypter),
		"NewCBCEncrypter":     reflect.ValueOf(cipher.NewCBCEncrypter),
		"NewCFBDecrypter":     reflect.ValueOf(cipher.NewCFBDecrypter),
		"NewCFBEncrypter":     reflect.ValueOf(cipher.NewCFBEncrypter),
		"NewCTR":              reflect.ValueOf(cipher.NewCTR),
		"NewOFB":              reflect.ValueOf(cipher.NewOFB),
	}

	env.PackageTypes["crypto/cipher"] = map[string]reflect.Type{}
}
