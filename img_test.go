package main

import (
	"fmt"
	"hdbdown/pool"
	"os"
	"strings"

	"testing"
)

func TestAdd(t *testing.T) {
	src := "big.jpg"

	dst := strings.Replace(src, ".", "_small.", 1)

	fmt.Println("src=", src, " dst=", dst)

	fIn, _ := os.Open(src)
	defer fIn.Close()
	fOut, _ := os.Create(dst)

	defer fOut.Close()
	// err := clip(fIn, fOut, 0, 0, 150, 150, 100)
	// if err != nil {
	// panic(err)
	// }
	img, fm, err := pool.Scale(fIn, 690, 390, 0)
	if err != nil {
		panic(err)
	}

	err = pool.Cutter(img, fOut, fm, 690, 390, 100)
	if err != nil {
		panic(err)
	}
}
