package main

import (
	"context"
	"errors"
	"log"
	"github.com/sota-rakuma/cat/myfile"
	"os"
	"runtime/trace"
	"strconv"
	"sync"
)


type file_key struct{}
type mu_key struct{}

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatal(err)
		return
	}
	defer trace.Stop()
	_main()
}

func _main() {
	if len(os.Args) <= 1 {
		log.Fatalln(errors.New("not enough arguments"))
		return
	}
	ctx, task := trace.NewTask(context.Background(), "reading files")
	defer task.End()

	setContext(&ctx, "")

	var wg sync.WaitGroup
	for i := 0; i < len(os.Args); i++ {
		wg.Add(1)
		num := i
		go func() {
			defer trace.StartRegion(ctx, "No." +  strconv.Itoa(num) + "goroutine").End()
			select {
			case <-ctx.Done():
				trace.Log(ctx, "No." +  strconv.Itoa(num) + "goroutine", "read canceled")
			default:
			}
			mf := myfile.NewFile(os.Args[num])
			mf.Read()
			outputFile(&ctx, mf)
			wg.Done()
		}()
	}
	wg.Wait()
}

func outputFile(ctx *context.Context, mf *myfile.MyFile) {
	mu := (*ctx).Value(mu_key{}).(sync.Mutex)
	f := (*ctx).Value(file_key{}).(*os.File)
	mu.Lock()
	defer mu.Unlock()
	f.Write(mf.Buff)
}
