package authmiddleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)


func GenToken(userId int) (string,error){
	claims := jwt.MapClaims{}
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	var secretKey = []byte("secretpassword")
	tokenPtr := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return tokenPtr.SignedString(secretKey)
}

func VerifyToken(access_token string)(jwt.MapClaims,error){
	claims := jwt.MapClaims{}
	token,err := jwt.ParseWithClaims(access_token,&claims,func(t *jwt.Token) (interface{}, error) {
		return []byte("secretpassword"),nil
	})

	if(err != nil){
		return nil,err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil,err
}

func Auth(next http.Handler) (http.Handler){
	// get token from session 
	// if token is not present send and error stating no token found.
	// var access_token string
	// claims,err := VerifyToken(access_token)
	// if(err != nil){
	// 	return nil,errors.New("invalid Token")
	// }
	return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie,err := r.Cookie("session_id")
		if(err != nil){
			fmt.Println(err)
		}
		if(cookie == nil){
			next.ServeHTTP(w,r)
			return
		}
		claims,err := VerifyToken(cookie.Value)
		if(err != nil){
			fmt.Println(err)
		}
		type claimskey int
		var claimsKey claimskey
		ctx := context.WithValue(r.Context(),claimsKey,claims)
		r = r.WithContext(ctx)
		fmt.Println(ctx.Value(claimsKey).(jwt.MapClaims))
		next.ServeHTTP(w,r)
	}))
	
}
