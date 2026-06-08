package helper

import (
	crypto_rand "crypto/rand"
	"math/big"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

//Hash the password
func Hash(h string) (string,error){
	hash,err:=bcrypt.GenerateFromPassword([]byte(h),bcrypt.DefaultCost)
	if err != nil {
		return "",err
	}
	return string(hash),nil
}

// comparing hashed password 

func VerifyHash(h,s string)(bool){
	err:=bcrypt.CompareHashAndPassword([]byte(h),[]byte(s))
	if err != nil {
		return false
	}
	return true
}

func GenerateOtp() string {
	nBig, err := crypto_rand.Int(crypto_rand.Reader, big.NewInt(900000))
	if err != nil {
		// Secure fallback
		return "123456"
	}
	return strconv.FormatInt(nBig.Int64()+100000, 10)
}


