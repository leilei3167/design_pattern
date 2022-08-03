package main

import (
	"errors"
	"fmt"
	"github.com/leilei3167/copy_design_pattern/db"
)

func main() {
	if err := somErr(); err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			fmt.Println("it is Not FOund Err:", err)
		}
		var e db.ErrDB
		if errors.As(err, &e) {
			fmt.Println("it s ErrDB!")
		}
	}

}

func somErr() error {
	return fmt.Errorf("someErr1:%w", db.ErrRecordNotFound)
}
