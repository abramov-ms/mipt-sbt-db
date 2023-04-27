package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"playwithredis/playwithredis"
	"time"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "playwithredis",
		Short: "Tool to save and load IMDB reviews using Redis",
		Args:  cobra.NoArgs,
		RunE:  runCli,
	}

	flagLoadMovies   []int
	flagClusterAddrs []string
	flagImportJSON   string
	flagClear        bool
)

func importJSON(ctx context.Context, c *playwithredis.Client, filename string) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	movies := make([]playwithredis.Movie, 0)
	if err := json.Unmarshal(bytes, &movies); err != nil {
		return err
	}

	client, err := playwithredis.NewClient(ctx, flagClusterAddrs)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	begin := time.Now()
	for _, m := range movies {
		if err := client.StoreMovie(ctx, m); err != nil {
			return err
		}
	}
	duration := time.Since(begin)
	log.Printf(
		"Imported %d movies in %v (%v per movie)\n",
		len(movies),
		duration,
		duration/time.Duration(len(movies)),
	)

	return nil
}

func runCli(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := playwithredis.NewClient(ctx, flagClusterAddrs)
	if err != nil {
		return err
	}

	if flagImportJSON != "" {
		if err := importJSON(ctx, client, flagImportJSON); err != nil {
			return err
		}
	}

	if flagClear {
		if err := client.Clear(ctx); err != nil {
			return err
		}
	}

	for _, id := range flagLoadMovies {
		begin := time.Now()
		movie, err := client.LoadMovie(ctx, playwithredis.ID(id))
		if err != nil {
			return err
		}
		end := time.Now()
		movie.Dump()
		log.Printf("Loaded movie in %v", end.Sub(begin))
	}

	return nil
}

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&flagClusterAddrs, "addrs", nil, "Addresses of Redis instances")
	rootCmd.PersistentFlags().StringVar(&flagImportJSON, "import", "", "Import JSON file into Redis")
	rootCmd.PersistentFlags().BoolVar(&flagClear, "clear", false, "Clear Redis data")
	rootCmd.PersistentFlags().IntSliceVar(&flagLoadMovies, "load", []int{}, "Load movies from Redis")
}

func execute() error {
	return rootCmd.Execute()
}
