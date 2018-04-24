package main

import (
	"flag"
	"log"
	"os"
	"fmt"
	"math/rand"
	"time"
)

const blockGenData = 1024*1024

var path *string
var size *int

func init(){
	path = flag.String("path", "", "destination path")
	size = flag.Int("size", 1024, "file size")

	rand.Seed(time.Now().Unix())
}

func main(){
	flag.Parse()

	if !chechParam(path, size) {
		return
	}

	fmt.Printf("path %s, size %d\n", path, size)

	var (
		data chan []byte
		stop chan struct{}
	)

	go genData(data, stop, *size)

	for i := 0; i < 10; i++ {
		d := <-data
		fmt.Printf("%s\n", d)
	}

}

func chechParam(s *string, i *int) bool {
	if *s == "" {
		log.Printf("destination path is empty")
		return false
	}
	if st, err := os.Stat(*s); os.IsNotExist(err) || !st.IsDir() {
		log.Printf("destination is not exist or not a dirrectory")
		return false
	}
	if *i <=0 {
		log.Printf("invalid size")
		return false
	}
	return true
}

func genData(data chan<- []byte, stop <-chan struct{}, size int){
	randSlice := make([]byte, size)
	blockData := make([]byte, blockGenData)

	for {
		if size <= blockGenData {
			makeRandomSlice(randSlice)
		} else {
			for i := 0; i < len(randSlice)%blockGenData; {
				makeRandomSlice(blockData)
				copy(randSlice[i:], blockData)
				i += size
			}
		}
		select {
		case <-stop:
			return
		default:
			data <- randSlice
		}
	}
}

func makeRandomSlice(p []byte) {
	rand.Read(p)
}