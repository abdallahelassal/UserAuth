package jwt

import (
	"errors"
	"time"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/golang-jwt/jwt/v4"
	
)





func CreateAccessToken(user *domain.User, secret string , expiry int)(accesToken string,err error){
	exp := time.Now().Add(time.Hour * time.Duration(expiry))

	claims := &JwtCustomClaims{
		UserID: user.ID.String(),
		UserName: user.UserName ,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	t , err := token.SignedString([]byte(secret))
	if err != nil {
		return "",err
	}
	return t ,nil
}



func RefreshToken(user *domain.User,secret string, expiry int)(refreshToken string, err error){
	claims := &JwtCustomRefreshToken{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiry))),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	t , err := token.SignedString([]byte(secret))
	if err != nil {
		return "",err
	}
	return t,nil
}

func IsAuthorized(tokenString string, secret string )(bool,error){

	token , err := jwt.Parse(tokenString,func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil , errors.New("Unexpected signing method")
		}
		return []byte(secret) , nil 
	})
	if err != nil {
	
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, domain.ErrTokenExpired
		}
		return false, domain.ErrInvalidToken
	}
	
	if !token.Valid {
	return false, domain.ErrInvalidToken
	}

	return  true , nil 

}

func ExtractIDFromToken(tokenString string,secret string)(string,error){
	token , err := jwt.ParseWithClaims(tokenString,&JwtCustomClaims{},func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return "" , domain.ErrInvalidToken
		}
		return []byte(secret) , nil
	})
	if err != nil {
		return "", err 
	}
	claims , ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return "" , domain.ErrInvalidToken
	}
	return claims.UserID , nil
}
