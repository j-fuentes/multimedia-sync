package cmd

import (
	"log"
	"os"

	"github.com/j-fuentes/multimedia-sync/pkg/feeds"
	"github.com/spf13/cobra"
)

var rootDir string

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronizes the feeds.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hasError := false
		for _, feedsFile := range args {
			bb, err := os.ReadFile(feedsFile)
			if err != nil {
				log.Fatalf("cannot open file: %+v", err)
			}

			ff, err := feeds.LoadFeeds(bb)
			if err != nil {
				log.Fatal(err)
			}
			rootDir = expandPath(rootDir)
			for idx, feed := range ff {
				numDownloads, err := feed.Sync(rootDir)
				if err != nil {
					log.Printf("%s feed#%d: %q(%s) kind %s FAILED to sync: %+v", feedsFile, idx, feed.Name, feed.ID, feed.Kind, err)
					hasError = true
				} else {
					log.Printf("%s feed#%d: %q(%s) kind %s SUCCEDED to sync: %d clips were synced", feedsFile, idx, feed.Name, feed.ID, feed.Kind, numDownloads)
				}
			}
		}
		if hasError {
			log.Fatalf("Some feeds failed to sync, check log")
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.PersistentFlags().StringVar(&rootDir, "root-dir", "./downloads", "Root directory for the sync.")
}
