package feeds

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"
)

const (
	// KindYoutubePlaylist is the feed kind of a Youtube playlist
	KindYoutubePlaylist = "YoutubePlaylist"
)

// Feed describes the properties of a feed
type Feed struct {
	// ID is the unique id of a feed
	ID string `yaml:"id"`
	// Name is the name of the feed
	Name string `yaml:"name"`
	// Kind is the kind of feed
	Kind string `yaml:"kind"`
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
