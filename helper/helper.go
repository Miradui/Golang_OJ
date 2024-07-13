package helper

import (
	"Project/define"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

type UserClaims struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	IsAdmin  int    `json:"is_admin"`
	jwt.RegisteredClaims
}

var myKey = []byte("oj_key")

func GenerateToken(identity, name string, isAdmin int) (string, error) {
	UserClaims := &UserClaims{
		Name:             name,
		Identity:         identity,
		RegisteredClaims: jwt.RegisteredClaims{},
		IsAdmin:          isAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func AnalyseToken(tokenString string) (*UserClaims, error) {
	UserClaims := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, UserClaims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims.Valid {
		return UserClaims, nil
	}
	return nil, fmt.Errorf("analyse Token Error:%v", err)
}

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func SendCodeByEmail(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Get <3099349352@qq.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码"
	e.HTML = []byte("验证码: <b>" + code + "</b>")

	return e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "3099349352@qq.com",
		"xggjrdhiurdwdedg", "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})

}

func GetUUID() string {
	return uuid.NewV4().String()
}

func GenerateVerificationCode() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(r.Intn(10))
	}
	return s
}

func SaveCode(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, os.ModePerm)
	if err != nil {
		return "", err
	}

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(code)
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return path, nil
}

func CheckGoCodeValid(path string) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	code := string(b)
	for i := 0; i < len(code)-6; i++ {
		if code[i:i+6] == "import" {
			var flag byte
			for i = i + 7; i < len(code); i++ {
				if code[i] == ' ' {
					continue
				}
				flag = code[i]
				break
			}
			if flag == '(' {
				for i = i + 1; i < len(code); i++ {
					if code[i] == ')' {
						break
					}
					if code[i] == '"' {
						t := ""
						for i = i + 1; i < len(code); i++ {
							if code[i] == '"' {
								break
							}
							t += string(code[i])
						}
						if _, ok := define.ValidGolangPackageMap[t]; !ok {
							return false, nil
						}
					}
				}
			} else if flag == '"' {
				t := ""
				for i = i + 1; i < len(code); i++ {
					if code[i] == '"' {
						break
					}
					t += string(code[i])
				}
				if _, ok := define.ValidGolangPackageMap[t]; !ok {
					return false, nil
				}
			}
		}
	}
	return true, nil
}
