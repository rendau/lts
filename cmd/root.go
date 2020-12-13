package cmd

import (
	"errors"
	"fmt"
	"github.com/rendau/lts/internal/domain/usecases"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var rootCmd = &cobra.Command{
	Use:   "lts <url> <request_count> <worker_count>",
	Short: "Load testing tool",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("missing required arguments")
		}
		if _, err := strconv.Atoi(args[1]); err != nil {
			return errors.New("`Request count` argument is not valid")
		}
		if _, err := strconv.Atoi(args[2]); err != nil {
			return errors.New("`Worker count` argument is not valid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ucs := usecases.New()

		uri := args[0]
		requestCount, _ := strconv.Atoi(args[1])
		workerCount, _ := strconv.Atoi(args[2])

		ucs.Run(uri, requestCount, workerCount)

		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
