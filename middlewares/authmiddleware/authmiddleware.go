package authmiddleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type claimskey int
var claimsKey claimskey


func GenToken(userId int) (string,error){
	claims := jwt.MapClaims{}
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	var secretKey = []byte("secretpassword")
	tokenPtr := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return tokenPtr.SignedString(secretKey)
}

func VerifyToken(accessToken string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(accessToken, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("secretpassword"), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

func HandleClaims(r *http.Request) (*jwt.MapClaims) {
	claims, ok := r.Context().Value(claimsKey).(jwt.MapClaims)
	if(!ok){
		return nil
	}
	return &claims
}

func HandleClaimsAuthRoute(r *http.Request) (jwt.MapClaims, bool) {
	claims, ok := r.Context().Value(claimsKey).(jwt.MapClaims)
	return claims, ok
}

func Auth(next http.Handler) (http.Handler){
	return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie,err := r.Cookie("session_id")
		if(err != nil){
			switch err {
				case http.ErrNoCookie:
					fmt.Println("No cookie found")
					next.ServeHTTP(w,r)
					return
				default:
					fmt.Println(err)
					fmt.Println("error fectching cookie")
					return 
			}
		}

		if(cookie == nil){
			next.ServeHTTP(w,r)
			return
		}
		claims,err := VerifyToken(cookie.Value)
		if(err != nil){
			fmt.Println(err)
			return 
		}
		ctx := context.WithValue(r.Context(),claimsKey,claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w,r)
	}))
}





