package merge

import (
	"fmt"
	"github.com/h2non/filetype"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Merge() {
	files := getVideo("C:\\Users\\zen\\Github\\trans-web\\videos")
	log.Println(files)
	for _, file := range files {
		srt := strings.Replace(file, filepath.Ext(file), ".srt", -1)
		log.Printf("xxx%v\n", srt)
		if isExist(srt) {
			base := filepath.Base(srt)
			dir := filepath.Dir(srt)
			space := strings.Replace(base, " ", "", -1)
			log.Printf("dir:%v\tbase:%v\n", dir, base)
			err := os.Rename(filepath.Join(dir, base), filepath.Join(dir, space))
			if err != nil {
				log.Printf("err:%v\n", err)
			}

			output := strings.Replace(file, filepath.Ext(file), "_with_subtit2le.mp4", -1)
			vf := strings.Join([]string{"subtitles='", filepath.Join(dir, space), "'"}, "")
			cmd := exec.Command("ffmpeg", "-i", file, "-c:v", "h264_nvenc", "-vf", vf, output)
			log.Printf("ffmpeg %v\n", cmd)
			b, err := cmd.CombinedOutput()
			log.Printf("ffmpeg %s has err:%v\n", string(b), err)
		}
	}
}
func getVideo(root string) []string {
	dirPath := root // 请将这里替换为你想要遍历的目录的实际路径
	filePaths := []string{}
	err := filepath.WalkDir(dirPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			f, _ := os.Open(absPath)
			// We only have to pass the file header = first 261 bytes
			head := make([]byte, 261)
			f.Read(head)
			if filetype.IsVideo(head) {
				filePaths = append(filePaths, absPath)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}
	}
	return filePaths
}
func isExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}
