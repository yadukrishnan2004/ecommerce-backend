package helper

import "golang.org/x/crypto/bcrypt"


//Hash the password
func Hash(h string) (string,error){
	hash,err:=bcrypt.GenerateFromPassword([]byte(h),bcrypt.DefaultCost)
	if err != nil {
		return "",err
	}
	return string(hash),nil
}