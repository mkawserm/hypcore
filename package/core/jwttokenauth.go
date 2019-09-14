package core

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"github.com/mkawserm/hypcore/package/models"
	"github.com/pascaldekloe/jwt"
	"strings"
	"time"
)

func JWTTokenAuth(ctx *HContext) *graphql.Field {
	return &graphql.Field{
		Type: graphql.String,
		Args: graphql.FieldConfigArgument{
			"username": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			glog.Infoln("JWT Token Auth")

			username := params.Args["username"].(string)
			password := params.Args["password"].(string)

			user := &models.User{Pk: username}

			if ctx.GetObject(user) {
				glog.Infoln(user.Password)
				if user.IsPasswordValid(password) {
					claims := jwt.Claims{Set: make(map[string]interface{})}

					claims.Subject = username
					claims.Set["group"] = user.GetGroup()
					claims.Set["uid"] = username

					now := time.Now().Round(time.Second)
					claims.Issued = jwt.NewNumericTime(now)

					if user.IsSuperUser() {
						claims.Expires = jwt.NewNumericTime(now.Add(time.Duration(ctx.AuthTokenSuperGroupTimeout)))
					} else if user.IsServiceUser() {
						claims.Expires = jwt.NewNumericTime(now.Add(time.Duration(ctx.AuthTokenServiceGroupTimeout)))
					} else if user.IsNormalUser() {
						claims.Expires = jwt.NewNumericTime(now.Add(time.Duration(ctx.AuthTokenNormalGroupTimeout)))
					} else {
						claims.Expires = jwt.NewNumericTime(now.Add(time.Duration(ctx.AuthTokenDefaultTimeout)))
					}

					if strings.HasPrefix(ctx.AuthAlgorithm, "HS") {
						token, err := claims.HMACSign(ctx.AuthAlgorithm, []byte(ctx.AuthSecretKey))
						if err == nil {
							return string(token), nil
						} else {
							glog.Errorf("HS Token generation error: %s\n", err.Error())
						}
					} else if strings.HasPrefix(ctx.AuthAlgorithm, "EdDSA") {
						var privateKey ed25519.PrivateKey

						key, err1 := LoadPrivateKey([]byte(ctx.AuthPrivateKey))
						if err1 != nil {
							glog.Errorf("EdDSA Token generation error: %s\n", err1.Error())
							return "", nil
						}

						switch t := key.(type) {
						case ed25519.PrivateKey:
							privateKey = t
						default:
							glog.Errorf("EdDSA invalid PrivateKey\n")
							return "", nil
						}

						token, err := claims.EdDSASign(privateKey)
						if err == nil {
							return string(token), nil
						} else {
							glog.Errorf("EdDSA Token generation error: %s\n", err.Error())
						}
					} else if strings.HasPrefix(ctx.AuthAlgorithm, "ES") {
						var privateKey *ecdsa.PrivateKey
						key, err1 := LoadPrivateKey([]byte(ctx.AuthPrivateKey))
						if err1 != nil {
							glog.Errorf("ES Token generation error: %s\n", err1.Error())
							return "", nil
						}

						switch t := key.(type) {
						case *ecdsa.PrivateKey:
							privateKey = t
						default:
							glog.Errorf("ECDSA invalid PrivateKey\n")
							return "", nil
						}

						token, err := claims.ECDSASign(ctx.AuthAlgorithm, privateKey)
						if err == nil {
							return string(token), nil
						} else {
							glog.Errorf("ES Token generation error: %s\n", err.Error())
						}
					} else if strings.HasPrefix(ctx.AuthAlgorithm, "RS") {
						var privateKey *rsa.PrivateKey

						key, err1 := LoadPrivateKey([]byte(ctx.AuthPrivateKey))
						if err1 != nil {
							glog.Errorf("RS Token generation error: %s\n", err1.Error())
							return "", nil
						}

						switch t := key.(type) {
						case *rsa.PrivateKey:
							privateKey = t
						default:
							glog.Errorf("RS invalid PrivateKey\n")
							return "", nil
						}

						token, err := claims.RSASign(ctx.AuthAlgorithm, privateKey)
						if err == nil {
							return string(token), nil
						} else {
							glog.Errorf("ES Token generation error: %s\n", err.Error())
						}
					} else if strings.HasPrefix(ctx.AuthAlgorithm, "PS") {
						var privateKey *rsa.PrivateKey

						key, err1 := LoadPrivateKey([]byte(ctx.AuthPrivateKey))
						if err1 != nil {
							glog.Errorf("PS Token generation error: %s\n", err1.Error())
							return "", nil
						}

						switch t := key.(type) {
						case *rsa.PrivateKey:
							privateKey = t
						default:
							glog.Errorf("PS invalid PrivateKey\n")
							return "", nil
						}

						token, err := claims.RSASign(ctx.AuthAlgorithm, privateKey)
						if err == nil {
							return string(token), nil
						} else {
							glog.Errorf("PS Token generation error: %s\n", err.Error())
						}
					}
				} else {
					glog.Errorf("Invalid password\n")
				}
			} else {
				glog.Errorf("User '{%s}' not found\n", username)
			}

			return "", nil
		},
		Description: "Generate JWT token",
	}
}
