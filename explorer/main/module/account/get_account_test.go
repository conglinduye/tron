package account

import (
	"fmt"
	"testing"

	"github.com/wlcy/tron/explorer/core/utils"
)

func TestGetAcc(*testing.T) {
	accs, err := GetRawAccount([]string{
		"TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9",
		"TAahLbGTZk6YuCycii72datPQEtyC5x231",
		"TV9QitxEJ3pdiAUAfJ2QuPxLKp9qTTR3og",
	})

	fmt.Printf("%v\n%v\n", err, utils.ToJSONStr(accs))
}
