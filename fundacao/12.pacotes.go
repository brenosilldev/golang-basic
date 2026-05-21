package fundacao

import (
	"fmt"

	"github.com/google/uuid"
)

func PacoteTest() {

	uuidcreate, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(uuidcreate)
}
