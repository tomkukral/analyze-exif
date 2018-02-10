package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

var dir string

func init() {
	flag.Parse()

	if len(flag.Args()) >= 1 {
		dir = flag.Arg(0)
	} else {
		fmt.Println("Missing directory argument")
		os.Exit(2)
	}
}

func worker(jobs <-chan string, data *[]exifData, wg *sync.WaitGroup) {
	for filename := range jobs {

		func() {
			exifData := GetExif(filename)
			*data = append(*data, exifData)
			defer wg.Done()
		}()
	}

}

const (
	workers int = 100
)

func main() {
	var wg sync.WaitGroup

	jobs := make(chan string, workers)
	data := make([]exifData, 0)

	files := FindFiles(dir)
	fmt.Printf("Discovered %d file(s)\n", len(files))
	wg.Add(len(files))

	// start workers
	for i := 0; i < workers; i++ {
		go worker(jobs, &data, &wg)
	}

	for _, filename := range files {
		jobs <- filename
	}

	wg.Wait()

	fmt.Printf("Got %d results\n", len(data))
	fmt.Printf("%+v\n", data)
}
