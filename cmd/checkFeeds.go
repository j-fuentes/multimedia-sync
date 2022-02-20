package cmd

import (
	"log"
	"os"

	"github.com/j-fuentes/multimedia-sync/pkg/feeds"
	"github.com/spf13/cobra"
)

// checkFeedsCmd represents the checkFeeds command
var checkFeedsCmd = &cobra.Command{
	Use:   "checkFeeds",
	Short: "Read as feeds file and check its format",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, feedsFile := range args {
			bb, err := os.ReadFile(feedsFile)
			if err != nil {
				log.Fatalf("cannot open file: %+v", err)
			}
			
			ff, err := feeds.LoadFeeds(bb)
			if err != nil {
				log.Fatal(err)
			}
			for idx, feed := range ff {
				log.Printf("%s feed#%d: %q(%s) kind %s", feedsFile, idx, feed.Name, feed.ID, feed.Kind)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(checkFeedsCmd)
}
