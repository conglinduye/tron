package router

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

func TestA(t *testing.T) {
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	fmt.Printf("hashedBytes:%v", string(hashedBytes))
}
