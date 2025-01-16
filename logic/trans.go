package logic

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"translate/model"
	"translate/util"
)

var (
	seed = rand.New(rand.NewSource(time.Now().Unix()))
)

const (
	TIMEOUT = 10 //second

)

func Start() {
	//proxy :=os.Setenv("PROXY", proxy)
	//}
	//if localDeeplxService != "" {
	//	os.Setenv("LocalDeeplxService", localDeeplxService)
	//}
	path := "/videos"
	if !isExist(path) {
		path = "videos"
		os.MkdirAll(path, os.ModePerm)
	}
	files := GetFiles(path)
	for _, file := range files {
		if strings.HasSuffix(file, ".srt") {
			Trans(file)
		}
	}
	log.Println("本轮翻译完成")
}

func Trans(fp string) {
	r := seed.Intn(2000)
	srt := fp
	//log.Fatalf("%v根据文件名:%s\t替换的字幕名:%s\n", p.GetPattern(), fp, srt)
	tmpname := strings.Join([]string{strings.Replace(srt, ".srt", "", 1), strconv.Itoa(r), ".srt"}, "")
	before := util.ReadInSlice(srt)
	fmt.Println(before)
	after, _ := os.OpenFile(tmpname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer func() {
		if err := recover(); err != nil {
			v := fmt.Sprintf("捕获到错误:%v\n", err)
			if strings.Contains(v, "index out of range") {
				fmt.Println("捕获到 index out of range 类型错误,忽略并继续执行重命名操作")
				{
					origin := strings.Join([]string{strings.Replace(srt, ".srt", "", 1), "_origin", ".srt"}, "")
					err1 := os.Rename(srt, origin)
					err2 := os.Rename(tmpname, srt)
					if err1 != nil || err2 != nil {
						fmt.Printf("字幕文件重命名出现错误:%v\n", err)
					}
				}
				return
			} else {
				log.Fatalf("捕获到其他错误:%v\n", v)
			}
		}
	}()
	for i := 0; i < len(before); i += 4 {
		if i+3 > len(before) {
			continue
		}
		log.Printf("翻译之前序号\"%s\"时间\"%s\"正文\"%s\"空行\"%s\"\n", before[i], before[i+1], before[i+2], before[i+3])
		log.SetPrefix(before[i])
		after.WriteString(before[i])
		after.WriteString(before[i+1])
		src := before[i+2]
		afterSrc := src
		one := new(model.TranslateHistory)
		one.Src = afterSrc
		var dst string
		if has, err := one.FindBySrc(); has {
			dst = one.Dst
			fmt.Printf("在缓存中找到:%s\n", dst)
		} else if err != nil {
			fmt.Printf("未在缓存中找到:%s\n", src)
			dst = Translate(afterSrc)
			one.Dst = dst
			randomNumber := seed.Intn(401) + 100
			time.Sleep(time.Duration(randomNumber) * time.Millisecond) // 暂停 100 毫秒
		} else {
			fmt.Printf("未在缓存中找到:%s\n", src)
			dst = Translate(afterSrc)
			one.Dst = dst
			randomNumber := seed.Intn(401) + 100
			time.Sleep(time.Duration(randomNumber) * time.Millisecond) // 暂停 100 毫秒
		}
		_, err := one.InsertOne()
		if err != nil {
			log.Printf("插入数据库失败:%v\n", err)
			return
		}
		log.Printf("翻译之后序号\"%s\"时间\"%s\"正文\"%s\"空行\"%s\"\n", before[i], before[i+1], before[i+2], before[i+3])
		log.Printf("原文\"%s\"\t译文\"%s\"\n", src, dst)
		after.WriteString(src)
		after.WriteString(dst)
		after.WriteString(before[i+3])
		after.WriteString(before[i+3])
		after.Sync()
	}
	after.Close()
	origin := strings.Join([]string{strings.Replace(srt, ".srt", "", 1), "_origin", ".srt"}, "")
	err1 := os.Rename(srt, origin)
	err2 := os.Rename(tmpname, srt)
	if err1 != nil || err2 != nil {
		fmt.Printf("字幕文件重命名出现错误:%v:%v\n", err1, err2)
	}
	log.Println("单次翻译完成")
}

/*
执行命令过程中可以循环打印消息
*/
func execCommand(c *exec.Cmd) (e error) {
	log.Printf("开始执行命令:%v\n", c.String())
	stdout, err := c.StdoutPipe()
	c.Stderr = c.Stdout
	if err != nil {
		log.Printf("连接Stdout产生错误:%v\n", err)
		return err
	}
	if err = c.Start(); err != nil {
		log.Printf("启动cmd命令产生错误:%v\n", err)
		return err
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		t := string(tmp)
		t = strings.Replace(t, "\u0000", "", -1)
		fmt.Printf("\r%v", t)
		if err != nil {
			break
		}
	}
	if err = c.Wait(); err != nil {
		log.Printf("命令执行中产生错误:%v\n", err)
		return err
	}
	return nil
}
func GetFiles(currentDir string) (filePaths []string) {
	err := filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查是否是文件
		if !info.IsDir() {
			filePaths = append(filePaths, path) // 将文件的绝对路径添加到切片
		}
		return nil
	})

	if err != nil {
		fmt.Println("遍历目录失败:", err)
		return
	}

	// 打印所有文件的绝对路径
	for _, filePath := range filePaths {
		fmt.Println(filePath)
	}
	return filePaths
}
func Translate(src string) string {
	//trans -brief ja:zh "私の手の動きに合わせて|そう"
	var dst string
	if src == "" {
		return dst
	}
	//fmt.Println("富强|民主|文明|和谐|自由|平等|公正|法治|爱国|敬业|诚信|友善")
	once := new(sync.Once)
	wg := new(sync.WaitGroup)
	defer wg.Wait()
	ack := make(chan string, 1)
	wg.Add(1)
	//go TransByDeeplx(src, p.GetProxy(), once, wg, ack)
	if proxy := os.Getenv("PROXY"); proxy == "" {
		go TransByDeeplx(src, once, wg, ack)
	} else {
		go TransByGoogle(src, proxy, once, wg, ack)
		go TransByBing(src, proxy, once, wg, ack)
		go TransByDeeplx(src, once, wg, ack)
	}
	select {
	case dst = <-ack:
		//constant.Info(fmt.Sprintf("收到翻译结果:%v\n", dst))
	case <-time.After(TIMEOUT * time.Second): // 设置超时时间为5秒
		fmt.Printf("翻译超时,重试\n此时的src = %v\n", src)
		Translate(src)
	}
	if dst == "" {
		fmt.Printf("翻译结果为空,重试\n此时的src = %v\n", src)
		return src
	}
	dst = strings.Replace(dst, "\r\n", "", -1)
	dst = strings.Replace(dst, "\n", "", -1)
	return dst
}

func isExist(dirPath string) bool {
	//dirPath := "/path/to/your/directory" // 替换为你想要检查的目录路径
	_, err := os.Stat(dirPath)
	if err == nil {
		fmt.Printf("目录 %s 存在\n", dirPath)
		return true
	} else if os.IsNotExist(err) {
		fmt.Printf("目录 %s 不存在\n", dirPath)
		return false
	} else {
		fmt.Printf("访问目录 %s 时出错: %v\n", dirPath, err)
		return false
	}
}
