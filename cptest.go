package main

import (
	"flag"
	"log"
	"os"
)

var path *string
var size *int

func init(){
	path = flag.String("path", "", "destination path")
	size = flag.Int("size", 1024, "file size")
}

func main(){
	flag.Parse()

	if !chechParam(path, size) {
		return
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