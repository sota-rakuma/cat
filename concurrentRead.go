package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func setContext(ctx *context.Context, name string) {
	if ctx == nil {
		*ctx = context.Background()
	}

	f, err := getFile(name) // default
	if err != nil {
		log.Fatalln(err)
		return
	}
	var mu sync.Mutex
	*ctx = context.WithValue(*ctx, file_key{}, f)
	*ctx = context.WithValue(*ctx, mu_key{}, &mu)
	ctx_cp, cancel := context.WithCancel(*ctx)
	*ctx = ctx_cp
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT)
		select {
		case <-sig:
			cancel()
		default:
		}
	}()
}

func getFile(file string) (*os.File, error) {
	if file == "" {
		return os.Stdout, nil
	}
	if err := validateFile(file); err != nil {
		return nil, fmt.Errorf("cannot open or create file %v", err)
	}
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func validateFile(name string) error {
	_, err := os.Stat(name)
	return err
}