package youtube

import (
	"fmt"
	"strings"

	youtube "github.com/kkdai/youtube/v2"
)

func NewYoutubePlaylistFeed(url string) *YoutubePlaylistFeed {
	return &YoutubePlaylistFeed{
		url: url,
	}
}

type YoutubePlaylistFeed struct {
	url string
}

// SyncToDirectory makes sure all the videos in the playlist are in the directory. Returns the number of new videos downloaded and error.
func (p *YoutubePlaylistFeed) SyncToDirectory(dir string) (int, error) {
	id, err := getPlaylistIDFromURL(p.url)
	if err != nil {
		return 0, err
	}
	client := youtube.Client{}

	playlist, err := client.GetPlaylist(id)
	if err != nil {
		return 0, err
	}

	numVideos := 0
	// TODO: get list of existing videos and do not download them
	for _, video := range playlist.Videos {
		fmt.Println(video.Title)
		numVideos++
	}


	return numVideos, nil
}

func getPlaylistIDFromURL(url string) (string, error) {
	result := strings.Trim(url, "https://www.youtube.com/playlist?list=")
	if result == "" {
		return "", fmt.Errorf("cannot extract ID from %q", url)
	}

	return result, nil
}
