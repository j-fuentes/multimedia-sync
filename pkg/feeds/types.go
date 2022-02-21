package feeds

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/j-fuentes/multimedia-sync/pkg/feeds/youtube"
	"gopkg.in/yaml.v2"
)

const (
	// KindYoutubePlaylist is the feed kind of a Youtube playlist
	KindYoutubePlaylist = "youtube_playlist"
)

// Feed describes the properties of a feed.
type Feed struct {
	// ID is the unique id of a feed
	ID string `yaml:"id"`
	// Name is the name of the feed
	Name string `yaml:"name"`
	// Kind is the kind of feed
	Kind string `yaml:"kind"`
	// Config is general config for the feed
	Config *Config `yaml:"config"`

	// YoutubePlaylist is the config for a YoutubePlaylist feed.
	YoutubePlaylist *youtube.YoutubePlaylistFeed `yaml:"youtube_playlist"`
}

// FeedImplementation is the interface that all the feed implementations need to satisfy.
type FeedImplementation interface {
	SyncToDirectory(dir string) (int, error)
}

// Config is the general configuration for the feed.
type Config struct {
	// OnlyAudio downloads only the audio from video. It does nothing with audio feeds.
	// TODO: at the moment this does nothing, it is still to be implemented.
	OnlyAudio bool `yaml:"only_audio,omitempty"`
}

// Sync syncs the feed to its corresponding directory inside rootDir.
func (f *Feed) Sync(rootDir string) (int, error) {
	var impl FeedImplementation
	switch f.Kind {
	case KindYoutubePlaylist:
		impl = f.YoutubePlaylist
	default:
		return 0, fmt.Errorf("unknown kind %q", f.Kind)
	}

	dir := path.Join(rootDir, f.ID)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return 0, err
	}

	return impl.SyncToDirectory(dir)
}

// Validate returns the validation errors in the feed description if any
func (f *Feed) Validate() error {
	var result error

	isAlpha := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString

	if f.ID == "" {
		result = multierror.Append(result, errors.New("id cannot be empty"))
	} else {
		if !isAlpha(f.ID) {
			result = multierror.Append(result, errors.New("id must be an alphanumeric value"))
		}
	}
	if f.Name == "" {
		result = multierror.Append(result, errors.New("name cannot be empty"))
	}

	switch f.Kind {
	case KindYoutubePlaylist:
	default:
		result = multierror.Append(result, fmt.Errorf("invalid kind %q", f.Kind))

	}

	return result
}

// LoadFeeds loads yaml bytes into a slice of feeds. It performs validations.
func LoadFeeds(bb []byte) ([]*Feed, error) {
	var feeds []*Feed

	err := yaml.Unmarshal(bb, &feeds)
	if err != nil {
		return nil, err
	}

	var resultErr error
	for idx, feed := range feeds {
		if valErrs := feed.Validate(); valErrs != nil {
			resultErr = multierror.Append(resultErr, fmt.Errorf("feed #%d is not valid: %w", idx, valErrs))
		}
	}

	if resultErr != nil {
		feeds = nil
	}

	return feeds, resultErr
}
