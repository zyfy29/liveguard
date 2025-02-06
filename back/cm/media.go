package cm

import (
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const maxRetry = 6

func parsePlaylist(playlistUrl string) (*m3u8.MediaPlaylist, error) {
	resp, err := http.Get(playlistUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return nil, err
	}
	if listType != m3u8.MEDIA {
		return nil, fmt.Errorf("expected media playlist, got different type: %v", listType)
	}
	mediaPlaylist, ok := playlist.(*m3u8.MediaPlaylist)
	if !ok {
		return nil, fmt.Errorf("failed to parse %v to m3u8.MediaPlaylist", playlist)
	}
	return mediaPlaylist, nil
}

// DownloadPlaylistAudio download the video of given playlist to ./data/live-{uuid}.mp3
func DownloadPlaylistAudio(playlistUrl string) (string, error) {
	tmpTsDir, err := os.MkdirTemp("", "live-*.d")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpTsDir)

	mediaPlaylist, err := parsePlaylist(playlistUrl)
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(playlistUrl)
	if err != nil {
		return "", err
	}
	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	// download ts files to custom tmp dir
	var tsFiles []string
	retries := make(map[int]int)
	for i, segment := range mediaPlaylist.Segments {
		if segment == nil {
			continue
		}
		tsFile := filepath.Join(tmpTsDir, filepath.Base(segment.URI))
		if err := downloadFile(lo.Must(url.JoinPath(baseURL, segment.URI)), tsFile); err != nil {
			for retries[i] < maxRetry || err != nil {
				log.Printf("Retrying download for segment %d", i)
				retries[i]++
				time.Sleep(time.Duration(retries[i]) * 3 * time.Second)
				err = downloadFile(lo.Must(url.JoinPath(baseURL, segment.URI)), tsFile)
			}
			if err != nil {
				return "", errors.Wrapf(err, "max retries exceeded for segment %d", i)
			}
		}
		tsFiles = append(tsFiles, tsFile)
	}

	tmpMP3, err := mergeTSToMP3(tsFiles)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpMP3)

	targetFile := GetRandomDataFilePathWithNameAndExt("live", "mp3")
	if err := copyFile(tmpMP3, targetFile); err != nil {
		return "", err
	}
	return targetFile, nil
}

// mergeTSToMP3 merge tsFiles to a tmp mp3 file
func mergeTSToMP3(tsFiles []string) (string, error) {
	listFile, err := os.CreateTemp("", "*.txt")
	if err != nil {
		return "", err
	}
	outputFile, err := os.CreateTemp("", "*.mp3")
	if err != nil {
		return "", err
	}

	defer func() {
		listFile.Close()
		os.Remove(listFile.Name())
		outputFile.Close()
	}()

	for _, file := range tsFiles {
		_, err := listFile.WriteString(fmt.Sprintf("file '%s'\n", file))
		if err != nil {
			return "", err
		}
	}

	// ffmpegコマンドを実行して、tsファイルをmp3に変換する
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile.Name(),
		"-q:a", "0", // 最高音質
		"-map", "a", // 音声のみ
		outputFile.Name(),
		"-y", // 既存の出力ファイルがあれば上書き
	)

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outputFile.Name(), nil
}

// MergeMP4Files mp4 files -> mp4 file (tmp)
// unused
func MergeMP4Files(files []string) (string, error) {
	tmpFile, err := os.CreateTemp("", "*.mp4")
	if err != nil {
		return "", err
	}
	outputFile := tmpFile.Name()

	listFile, err := os.CreateTemp("", "*.txt")
	if err != nil {
		return "", err
	}

	for _, file := range files {
		_, err := listFile.WriteString(fmt.Sprintf("file '%s'\n", file))
		if err != nil {
			return "", err
		}
	}

	defer func() {
		listFile.Close()
		os.Remove(listFile.Name())
		tmpFile.Close()
	}()

	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile.Name(),
		"-c", "copy",
		outputFile,
		"-y", // Overwrite the output file if it exists
	)

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

// GetPlaylistDuration get video duration (seconds, ceiled) for playListUrl
func GetPlaylistDuration(playlistUrl string) (int, error) {
	mediaPlaylist, err := parsePlaylist(playlistUrl)
	if err != nil {
		return 0, err
	}
	var totalDuration float64

	for _, segment := range mediaPlaylist.Segments {
		if segment != nil {
			totalDuration += segment.Duration
		}
	}
	return int(math.Ceil(totalDuration)), nil
}
