package youtube

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	youtube "github.com/kkdai/youtube/v2"
)

type YoutubePlaylistFeed struct {
	URL string `yaml:"url"`
}

// SyncToDirectory makes sure all the videos in the playlist are in the directory. Returns the number of new videos downloaded and error.
func (p *YoutubePlaylistFeed) SyncToDirectory(dir string) (int, error) {
	id, err := getPlaylistIDFromURL(p.URL)
	if err != nil {
		return 0, err
	}

	log.Printf("syncing Youtube playlist %q to %q", id, dir)

	client := youtube.Client{}

	playlist, err := client.GetPlaylist(id)
	if err != nil {
		return 0, err
	}

	numVideos := 0
	// TODO: get list of existing videos and do not download them
	for _, playlistEntry := range playlist.Videos {
		video, err := client.VideoFromPlaylistEntry(playlistEntry)
		if err != nil {
			return numVideos, err
		}
		path := path.Join(dir, video.ID+".mp4")
		downloadVideo(client, video, path)
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

func downloadVideo(client youtube.Client, video *youtube.Video, path string) error {
	log.Printf("Downloading %s by '%s'!\n", video.Title, video.Author)
	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		return fmt.Errorf("cannot get stream: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}

	defer file.Close()
	_, err = io.Copy(file, stream)
	if err != nil {
		return fmt.Errorf("cannot download stream: %w", err)
	}

	return nil
}
