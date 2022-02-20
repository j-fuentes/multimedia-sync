package feeds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFeeds(t *testing.T) {
	tests := []struct {
		name    string
		input string
		want    []*Feed
		wantErr bool
	}{
		{
			"loads valid feeds",
			`
- id: abc1
  name: feed1
  kind: YoutubePlaylist
- id: abc2
  name: feed2
  kind: YoutubePlaylist
`,
			[]*Feed{
				{
					ID: "abc1",
					Name: "feed1",
					Kind: "YoutubePlaylist",
				},
				{
					ID: "abc2",
					Name: "feed2",
					Kind: "YoutubePlaylist",
				},
			},
			false,
		},
		{
			"detects invalid kind",
			`
- id: abc1
  name: feed1
  kind: YoutubePlaylist
- id: abc2
  name: feed2
  kind: Unexisting
`,
			nil,
			true,
		},
		{
			"detects wrong ID",
			`
- id: abc+1
  name: feed1
  kind: YoutubePlaylist
`,
			nil,
			true,
		},
		{
			"detects empty ID",
			`
- name: feed1
  kind: YoutubePlaylist
`,
			nil,
			true,
		},
		{
			"detects empty Name",
			`
- id: abc1
  kind: YoutubePlaylist
`,
			nil,
			true,
		},
		{
			"detects empty Kind",
			`
- id: abc1
  name: feed1
`,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadFeeds([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFeeds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, got, tt.want)
		})
	}
}
