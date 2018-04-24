package main

import (
	"flag"
	"log"
	"os"
	"math/rand"
	"time"
	"path"
	"strconv"
	"io/ioutil"
	"os/signal"
)

const (
	blockGenData = 1024*1024
)

var (
	dir *string
	size *int
	delta *int
	maxRepeat *int
)

func init(){
	dir = flag.String("dir", "", "destination dir")
	size = flag.Int("size", 512, "file size")
	delta = flag.Int("delta", 100, "maximum number of files")
	maxRepeat = flag.Int("max", 1024, "maximum number of files")

	rand.Seed(time.Now().Unix())
}

func main(){
	flag.Parse()

	if !checkParam(dir, size, delta, maxRepeat) {
		return
	}

	log.Printf("destination directory: %s, file size: %d, number of files: %d\n", *dir, *size, *maxRepeat)

	data := make(chan []byte)
	stop := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go genData(data, stop, *size)

	fullTime := time.Now()
	curTime := time.Now()
	j := 0
	i := 0
	STOP:
	for i = 0; i < *maxRepeat; i++ {
		select {
		case <-c:
			close(stop)
			break STOP
		default:
		}

		if j == *delta {
			curr := time.Now().Sub(curTime)
			log.Printf("processing time %d of files: %s \t file per seconds: %f\n", *delta, curr.String(), float64(j) / curr.Seconds())
			curTime = time.Now()
			j = 0
		}
		j++
		name := randomName(*dir)
		time.Sleep(time.Millisecond )
		if err := writeToFile(name, <-data); err != nil {
			log.Fatalf("error writing to file %s %+v", name, err)
		}

	}
	log.Println("------------------------------------------------------------------------------")
	total := time.Now().Sub(fullTime)
	log.Printf("total processing time %d files: %s \t file per seconds: %f\n", i, total.String(), float64(i) / total.Seconds())

}

func checkParam(s *string, i, delta, max *int) bool {
	if *s == "" {
		log.Printf("destination dir is empty")
		return false
	}
	st, err := os.Stat(*s)
	if err == nil && !st.IsDir() {
		log.Printf("destination is not a dirrectory")
		return false
	}
	if err == nil {
		log.Printf("dirrectory %s is exist", *s)
		log.Println("removing dirrectory...")
		if err = os.RemoveAll(*s); err != nil {
			log.Printf("error removing directory %s %+v", *s, err)
			return false
		}
		log.Println("removing dirrectory... done")
		err = os.ErrNotExist
	}

	if os.IsNotExist(err) {
		if err = os.MkdirAll(*s, 0777); err != nil {
			return false
		}
	}

	if *i <=0 {
		log.Printf("invalid size")
		return false
	}

	if *delta <=0 {
		log.Printf("invalid number of delta")
		return false
	}
	if *max <=0 {
		log.Printf("invalid number of files")
		return false
	}
	return true
}

func genData(data chan []byte, stop <-chan struct{}, size int){
	for {
		randSlice := make([]byte, size)
		if size <= blockGenData {
			makeRandomSlice(randSlice)
		} else {
			blockData := make([]byte, blockGenData)
			i := 0
			for i < len(randSlice){
				makeRandomSlice(blockData)
				copy(randSlice[i:], blockData)
				i += blockGenData
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

func randomName(fulldir string) string {
	return path.Join(fulldir, "test-" + strconv.FormatInt(rand.Int63(), 16))
}

func writeToFile(name string, d []byte) error {
	return ioutil.WriteFile(name, d, 0666)
}