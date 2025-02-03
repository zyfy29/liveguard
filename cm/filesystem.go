package cm

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func copyFile(src, dest string) error {
	// コピー先のディレクトリが存在しない場合は作成
	destDir := filepath.Dir(dest)
	// os.MkdirAllでは、存在しない階層をすべて作ってくれる
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// コピー元ファイルを開く
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// コピー先ファイルを作成（存在する場合は上書き）
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// ファイルをコピー
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// コピー先ファイルのパーミッションを元ファイルと同じにする
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}
	if err := os.Chmod(dest, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}
	return nil
}

func getAbsDataDir() string {
	return filepath.Join(GetProjectRoot(), viper.GetString("app.data_dir"))
}

func GetDataFilePath(parts ...string) string {
	res := filepath.Join(append([]string{getAbsDataDir()}, parts...)...)
	_ = os.MkdirAll(filepath.Dir(res), os.ModePerm)
	return res
}

func GetRandomDataFilePathWithNameAndExt(name, ext string) string {
	fileName := uuid.NewString()
	if name != "" {
		fileName = name + "-" + fileName
	}
	return GetDataFilePath(
		fmt.Sprintf("%s.%s", fileName, ext),
	)
}

func GetProjectRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		_, err := os.ReadFile(filepath.Join(currentDir, "go.mod"))
		if os.IsNotExist(err) {
			if currentDir == filepath.Dir(currentDir) {
				return ""
			}
			currentDir = filepath.Dir(currentDir)
			continue
		} else if err != nil {
			return ""
		}
		break
	}
	return currentDir
}

// downloadFile 将url的文件下载到localPath，localPath是文件名，不考虑母路径不存在的情况
func downloadFile(url, localPath string) error {
	log.Printf("Downloading %s to %s", url, localPath)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Add the User-Agent header
	req.Header.Set("User-Agent", "SNH48 ENGINE")

	// Perform the request
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
