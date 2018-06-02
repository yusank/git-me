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

type URLData struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
}

// VideoData data struct of video info
type VideoData struct {
	Site      string `json:"site"`
	Title     string `json:"title"`
	OutputDir string `json:"-"`
	// [URLData: {URL, Size, Ext}, ...]
	// Some video files have multiple fragments
	// and support for downloading multiple image files at once

	Formats map[string]FormatData
	Type    string
}

// FormatData data struct of every format
type FormatData struct {
	// [URLData: {URL, Size, Ext}, ...]
	// Some video files have multiple fragments
	// and support for downloading multiple image files at once
	URLs    []URLData `json:"urls"`
	Quality string    `json:"quality"`
	Size    int64     `json:"size"` // total size of all urls
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type InnerTask struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	URL    string `json:"url"`
	Status int    `json:"status"`
	Sort   int    `json:"sort"`
	Tp     int    `json:"tp"`
}

type UploadInfo struct {
	URL      string
	Status   int
	Schedule float64
}
