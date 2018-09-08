package grpcclient

import (
	"fmt"
	"testing"

	"github.com/wlcy/tron/explorer/core/utils"
)

func TestFoo(*testing.T) {
	base64ID := "JFXopSY7F+nbrEwIo8lp1SYzZ2KMrOROXfX6yw+T84c="
	fmt.Printf("%v", utils.HexEncode(utils.Base64Decode(base64ID)))
}
