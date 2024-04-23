package main

import (
	"context"
	"fmt"
	"go-magicinfo/internal/api"
	"go-magicinfo/internal/auth"
	"go-magicinfo/internal/client"
	"go-magicinfo/internal/config"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	c := &http.Client{}

	t, err := auth.RestAPIToken(ctx, c, cfg.OldUrl, cfg.User, cfg.Password)
	if err != nil {
		os.Exit(1)
	}

	slog.Info("RestAPIToken received", "msg", t)

	fmt.Printf("Mac %s\n", cfg.MacAddresses)

	mc := client.New(c, t)

	// Semaphore
	sem := make(chan struct{}, cfg.MaxConn)
	var wg sync.WaitGroup

	for _, mac := range cfg.MacAddresses {
		// Acquire semaphore
		sem <- struct{}{}
		wg.Add(1)

		go func(ctx context.Context, mc *client.Client, mac, oldUrl, newUrl string) {
			defer wg.Done()
			fmt.Printf("MAC: %s\n", mac)
			api.ChangeUrl(ctx, mc, mac, oldUrl, newUrl)
			<-sem
		}(ctx, mc, mac, cfg.OldUrl, cfg.NewUrl)
	}

	wg.Wait()
	slog.Info("Done!")
}
