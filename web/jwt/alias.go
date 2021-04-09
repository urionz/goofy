package jwt

import "github.com/kataras/iris/v12/middleware/jwt"

type (
	Alg            = jwt.Alg
	Signer         = *jwt.Signer
	TokenValidator = jwt.TokenValidator
	Verifier       = *jwt.Verifier
)
