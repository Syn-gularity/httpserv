package auth

import(
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error){
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}
	//pw_ret := fmt.Sprintf("%x",pw)
	return string(pw), nil
}

func CheckPasswordHash(password, hash string) error{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil{
		return err
	}
	return nil
}