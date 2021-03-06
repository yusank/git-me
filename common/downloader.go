/*
 * MIT License
 *
 * Copyright (c) 2018 Yusan Kurban <yusankurban@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/04/01        Yusan Kurban
 */

package common

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/fatih/color"

	"github.com/yusank/git-me/utils"
)

func (data *FormatData) calculateTotalSize() {
	var size int64
	for _, urlData := range data.URLs {
		size += urlData.Size
	}
	data.Size = size
}

// urlSave save url file
func (data FormatData) urlSave(
	urlData URLData, refer, fileName, output string, bar *pb.ProgressBar,
) {
	filePath := utils.FilePath(fileName, urlData.Ext, output, false)
	fileSize, exists := utils.FileSize(filePath)
	// TODO: Live video URLs will not return the size
	if fileSize == urlData.Size {
		fmt.Printf("%s: file already exists, skipping\n", filePath)
		bar.Add64(fileSize)
		return
	}
	if exists && fileSize != urlData.Size {
		// files with the same name but different size
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s: file already exists, overwriting? [y/n]", filePath)
		overwriting, _ := reader.ReadString('\n')
		overwriting = strings.Replace(overwriting, "\n", "", -1)
		if overwriting != "y" {
			return
		}
	}
	tempFilePath := filePath + ".download"
	tempFileSize, _ := utils.FileSize(tempFilePath)
	headers := map[string]string{
		"Referer": refer,
	}
	var file *os.File
	if tempFileSize > 0 {
		// range start from zero
		headers["Range"] = fmt.Sprintf("bytes=%d-", tempFileSize)
		file, _ = os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY, 0644)
		bar.Add64(tempFileSize)
	} else {
		file, _ = os.Create(tempFilePath)
	}

	// close and rename temp file at the end of this function
	// must be done here to avoid the following request error to cause the file can't close properly
	defer func() {
		file.Close()
		// must close the file before rename or it will cause `The process cannot access the file because it is being used by another process.` error.
		err := os.Rename(tempFilePath, filePath)
		if err != nil {
			log.Fatal(err)
		}
	}()

	res, _ := utils.HttpGet(urlData.URL, headers)
	if res.StatusCode >= 400 {
		red := color.New(color.FgRed)
		log.Print(urlData.URL)
		log.Fatal(red.Sprintf("HTTP error: %d", res.StatusCode))
	}
	defer res.Body.Close()
	writer := io.MultiWriter(file, bar)
	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	_, copyErr := io.Copy(writer, res.Body)
	if copyErr != nil {
		log.Fatal(fmt.Sprintf("Error while downloading: %s, %s", urlData.URL, copyErr))
	}
}

func printStream(k string, data FormatData) {
	blue := color.New(color.FgHiYellow)
	cyan := color.New(color.FgCyan)
	blue.Println(fmt.Sprintf("     [%s]  -------------------", k))
	if data.Quality != "" {
		cyan.Printf("     Quality:         ")
		fmt.Println(data.Quality)
	}
	cyan.Printf("     Size:            ")
	if data.Size == 0 {
		data.calculateTotalSize()
	}
	fmt.Printf("%.2f MiB (%d Bytes)\n", float64(data.Size)/(1024*1024), data.Size)
	cyan.Printf("     # download with: ")
	fmt.Println("git-me -f " + k + " [URL]")
	fmt.Println()
}

func (vid VideoData) printInfo(format string) {
	cyan := color.New(color.FgCyan)
	fmt.Println()
	cyan.Printf(" Site:      ")
	fmt.Println(vid.Site)
	cyan.Printf(" Title:     ")
	fmt.Println(vid.Title)
	cyan.Printf(" Type:      ")
	fmt.Println(vid.Type)
	if InfoOnly {
		cyan.Printf(" Streams:   ")
		fmt.Println("# All available quality")
		for k, data := range vid.Formats {
			printStream(k, data)
		}
	} else {
		cyan.Printf(" Stream:   ")
		fmt.Println()
		printStream(format, vid.Formats[format])
	}
}

// ParseVideo download urls
func (vid VideoData) Download(refer string) {
	var format, title string
	if Format == "" {
		// default is best quality
		format = "default"
	} else {
		format = Format
	}

	title = utils.FileName(vid.Title)
	data := vid.Formats[format]
	ok := len(vid.Formats) != 0
	if !ok {
		log.Fatal("No format named " + format)
	}
	if data.Size == 0 {
		data.calculateTotalSize()
	}
	vid.printInfo(format)
	if InfoOnly {
		return
	}

	bar := pb.New64(data.Size).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.ShowFinalTime = true
	bar.SetMaxWidth(1000)
	bar.Start()
	go MonitorSchedule(bar, refer)
	if len(data.URLs) == 1 {
		// only one fragment
		data.urlSave(data.URLs[0], refer, title, vid.OutputDir, bar)
		bar.Finish()
		return
	}
	wgp := utils.NewWaitGroupPool(16)
	// multiple fragments
	var parts []string
	for index, url := range data.URLs {
		partFileName := fmt.Sprintf("%s[%d]", title, index)
		partFilePath := utils.FilePath(partFileName, url.Ext, vid.OutputDir, false)
		parts = append(parts, partFilePath)

		wgp.Add()
		go func(url URLData, refer, fileName string, bar *pb.ProgressBar) {
			defer wgp.Done()
			data.urlSave(url, refer, fileName, vid.OutputDir, bar)
		}(url, refer, partFileName, bar)

	}
	wgp.Wait()
	bar.Finish()

	if vid.Type != "video" {
		return
	}
	// merge
	mergeFileName := title + ".txt" // merge list file should be in the current directory
	filePath := utils.FilePath(title, "mp4", vid.OutputDir, false)
	fmt.Printf("正将散列文件合并至 %s\n", filePath)
	var cmd *exec.Cmd
	if strings.Contains(vid.Site, "youtube") {
		// merge audio and video
		cmds := []string{
			"-y",
		}
		for _, part := range parts {
			cmds = append(cmds, "-i", part)
		}
		cmds = append(
			cmds, "-c:vid", "copy", "-c:a", "aac", "-strict", "experimental",
			filePath,
		)
		cmd = exec.Command("ffmpeg", cmds...)
	} else {
		// write ffmpeg input file list
		mergeFile, _ := os.Create(mergeFileName)
		for _, part := range parts {
			mergeFile.Write([]byte(fmt.Sprintf("file '%s'\n", part)))
		}
		mergeFile.Close()

		cmd = exec.Command(
			"ffmpeg", "-y", "-f", "concat", "-safe", "-1",
			"-i", mergeFileName, "-c", "copy", "-bsf:a", "aac_adtstoasc", filePath,
		)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + "\n" + stderr.String())
	}
	// remove parts
	os.Remove(mergeFileName)
	for _, part := range parts {
		os.Remove(part)
	}

	fmt.Printf("合并完成 \n")
	return
}

func MonitorSchedule(pb *pb.ProgressBar, uri string) {
	for {
		time.Sleep(1 * time.Second)
		cur := pb.Get()
		process := float64(cur) / float64(pb.Total) * 100

		process = utils.RoundSpec(process, 2)
		upload := UploadInfo{
			URL:      uri,
			Schedule: process,
			Status:   TaskStatusFinish,
		}

		ProcessChan <- upload

		if pb.IsFinished() {
			break
		}
	}
}
