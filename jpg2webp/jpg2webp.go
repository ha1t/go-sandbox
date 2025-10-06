package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		for _, file := range os.Args[1:] {
			convertWebP(file)
		}
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("対象となるディレクトリのパスを入力してください。yでカレントディレクトリに対して実行します: ")
	dir, _ := reader.ReadString('\n')
	dir = strings.TrimSpace(dir)

	var path string
	if dir == "y" {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("カレントディレクトリの取得に失敗しました:", err)
			return
		}
		path = currentDir
	} else {
		path = dir
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("ディレクトリが存在しません:", err)
		return
	}

	if fileInfo.IsDir() {
		fmt.Println("ディレクトリが存在します")
		convertWebP(path)
	} else {
		convertWebP(path) // ファイルが直接指定された場合
	}

	fmt.Print("プログラムは正常に終了しました。なにか入力してください: ")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func convertWebP(dir string) {
	var files []string
	fileInfo, err := os.Stat(dir)
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}

	if fileInfo.IsDir() {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (strings.HasSuffix(strings.ToLower(path), ".jpg") ||
				strings.HasSuffix(strings.ToLower(path), ".gif") ||
				strings.HasSuffix(strings.ToLower(path), ".png")) {
				files = append(files, path)
			}
			return nil
		})
	} else {
		files = append(files, dir) // ファイルが直接指定された場合
	}

	for _, file := range files {
		ext := filepath.Ext(file)
		generateFileName := file[0:len(file)-len(ext)] + ".webp"
		fmt.Println("<Source File Name:", file, ">")

		// cwebp
		var cwebpCmd *exec.Cmd
		if strings.ToLower(ext) == ".jpg" {
			cwebpCmd = exec.Command("cwebp", "-psnr", "42", "-qrange", "40", "95", "-metadata", "all", "-sharp_yuv", "-o", generateFileName, "-progress", "-short", file)
		} else {
			cwebpCmd = exec.Command("cwebp", "-q", "75", "-metadata", "all", "-sharp_yuv", "-o", generateFileName, "-progress", "-short", file)
		}
		cwebpErr := cwebpCmd.Run()
		if cwebpErr != nil {
			fmt.Println("cwebpの実行に失敗しました:", cwebpErr)
			continue
		}

		// cwebp -metadata で exif を拾いきれないのでexiftoolを使い、magickでautorotateする
		// https://stackoverflow.com/questions/76681968/cwebp-metadata-all-does-not-keep-metadata

		// exiftool
		exiftoolCmd := exec.Command("exiftool", "-overwrite_original", "-TagsFromFile", file, generateFileName)
		exiftoolErr := exiftoolCmd.Run()
		if exiftoolErr != nil {
			fmt.Println("exiftoolの実行に失敗しました:", exiftoolErr)
			continue
		}

		// magick
		magickCmd := exec.Command("magick", generateFileName, "-auto-orient", generateFileName)
		magickErr := magickCmd.Run()
		if magickErr != nil {
			fmt.Println("magickの実行に失敗しました:", magickErr)
			continue
		}
	}
}
