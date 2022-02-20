package feeds

// SyncFeeds synchronizes multiple feeds to a base directory.
func SyncFeeds(baseDir string, feeds []*Feed) error {
	for _, feed := range feeds {
		err := SyncFeed(baseDir, feed)
		if err != nil {
			return err
		}
	}
	return nil
}

// SyncFeed synchronizes one feed to a base directory.
func SyncFeed(baseDir string, feed *Feed) error {
	// TODO: implement proper conversion between config and feeds implementation


	return nil
}
