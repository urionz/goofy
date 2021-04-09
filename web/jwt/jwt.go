package jwt

import (
	"time"

	"github.com/kataras/iris/v12/middleware/jwt"
)

func NewSigner(signatureAlg Alg, signatureKey interface{}, maxAge time.Duration) Signer {
	return jwt.NewSigner(signatureAlg, signatureKey, maxAge)
}

func NewVerifier(signatureAlg Alg, signatureKey interface{}, validators ...TokenValidator) Verifier {
	return jwt.NewVerifier(signatureAlg, signatureKey, validators...)
}
