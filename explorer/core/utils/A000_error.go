package utils

import "fmt"

// error define
var (
	ErrorDecrypt = fmt.Errorf("Decrypt Failed")

	ErrorInvalidSign = fmt.Errorf("Invalid Transaction Signature, should be 65 length bytes")

	ErrorCreateGrpClient = fmt.Errorf("Create GRPC Client failed")

	ErrorNotImplement = fmt.Errorf("Not implement")
)
