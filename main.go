package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//noinspection GoUnusedConst
const (
	kb = 1024
	mb = kb * 1024
	gb = mb * 1024
)

func main() {
	os.Exit(run())
}

func run() int {
	var (
		numberOfConcurrent int
		blockSize          int
	)

	var (
		wg              = sync.WaitGroup{}
		rootCtx         = context.Background()
		mainCtx, cancel = context.WithCancel(rootCtx)
	)

	flag.IntVar(&numberOfConcurrent, "g", 100, "goroutine count")
	flag.IntVar(&blockSize, "b", 1024, "block size per write")
	flag.Parse()

	log.Printf("goroutine count: %d\tblock size: %d\n", numberOfConcurrent, blockSize)

	dir, err := ioutil.TempDir("", "diskio")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[Dir ] ", dir)

	wg.Add(numberOfConcurrent)
	for gc := 0; gc < numberOfConcurrent; gc++ {
		go func(index int, ctx context.Context) {
			defer wg.Done()

			var (
				fname = fmt.Sprintf("disk-io-data.%04d", index)
				fpath = filepath.Join(dir, fname)
			)

			// https://qiita.com/kamol/items/fae07e8533b36f553714
			file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
			if err != nil {
				log.Fatal(err)
			}

			//noinspection GoUnhandledErrorResult
			defer file.Close()

			for i := 0; true; i++ {
				data := make([]byte, blockSize)
				for j := 0; j < blockSize; j++ {
					data[j] = 1
				}

				_, _ = file.Write(data)
				_ = file.Sync()

				select {
				case <-ctx.Done():
					return
				default:
				}
			}
		}(gc, mainCtx)
	}

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		var (
			tick = time.Tick(15 * time.Second)
		)

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick:
				// https://stackoverflow.com/questions/32482673/how-to-get-directory-total-size
				var totalSize uint64 = 0

				files, _ := ioutil.ReadDir(dir)
				for _, f := range files {
					fpath := filepath.Join(dir, f.Name())
					stat, _ := os.Stat(fpath)
					totalSize += uint64(stat.Size())
				}

				log.Printf("[Size] %04d MB\n", int(totalSize/mb))
			}
		}
	}(mainCtx)

	fmt.Println("Press ENTER to stop...")

	s := ""
	_, _ = fmt.Scanln(&s)

	cancel()
	wg.Wait()

	fmt.Println("Remove temporary directory: ", dir)
	err = os.RemoveAll(dir)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("DONE")
	fmt.Println("Press ENTER to exit...")

	_, _ = fmt.Scanln(&s)

	return 0
}
