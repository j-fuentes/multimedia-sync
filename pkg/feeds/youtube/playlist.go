package youtube

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	youtube "github.com/kkdai/youtube/v2"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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

	tmpPath := path + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}

	defer file.Close()
	_, err = io.Copy(file, stream)
	if err != nil {
		return fmt.Errorf("cannot download stream: %w", err)
	}

	defer os.Remove(tmpPath)
	return addMetadata(tmpPath, path, video)
}

func addMetadata(in string, out string, video *youtube.Video) error {
	args := []ffmpeg.KwArgs{
		{"c": "copy"},
	}
	// metadata keys: https://wiki.multimedia.cx/index.php/FFmpeg_Metadata
	metadata := map[string]string{
		"title": video.Title,
		"author": video.Author,
		"description": video.Description,
		"show": video.Author,
	}
	for k, v := range metadata {
		args = append(args, ffmpeg.KwArgs{"metadata": fmt.Sprintf("%s=%s", k, v)})
	}
	// -codec=copy: so it does not transcode and just edits metadata
	return ffmpeg.Input(in).Output(out, args...).OverWriteOutput().ErrorToStdOut().Run()
}
