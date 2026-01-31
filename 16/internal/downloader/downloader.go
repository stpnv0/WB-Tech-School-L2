package downloader

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"wb-wget/internal/parser"
	"wb-wget/internal/urlutil"
)

type Downloader struct {
	baseURL    string
	outputDir  string
	maxDepth   uint
	timeout    time.Duration
	numWorkers int

	visited      map[string]struct{}
	visitedMutex sync.Mutex

	httpClient *http.Client
}

type Task struct {
	URL   string
	Depth uint
}

type DownloadResult struct {
	URL         string
	LocalPath   string
	ContentType string
	Depth       uint
	Links       []parser.LinkInfo
	Err         error
}

func NewDownloader(baseURL, outputDir string, maxDepth uint, timeout time.Duration, numWorkers int) *Downloader {
	client := &http.Client{
		Timeout: timeout,
	}

	return &Downloader{
		baseURL:      baseURL,
		outputDir:    outputDir,
		maxDepth:     maxDepth,
		timeout:      timeout,
		numWorkers:   numWorkers,
		visited:      make(map[string]struct{}),
		visitedMutex: sync.Mutex{},
		httpClient:   client,
	}
}

func (d *Downloader) Run(startURL string) error {
	tasks := make(chan Task, d.numWorkers*2)
	var wg sync.WaitGroup

	enqueue := func(urlStr string, depth uint) {
		normalizedURL, err := urlutil.Normalize(urlStr)
		if err != nil {
			return
		}

		d.visitedMutex.Lock()
		if _, ok := d.visited[normalizedURL]; ok {
			d.visitedMutex.Unlock()
			return
		}
		d.visited[normalizedURL] = struct{}{}
		d.visitedMutex.Unlock()

		wg.Add(1)
		tasks <- Task{URL: urlStr, Depth: depth}
	}

	enqueue(startURL, 0)

	go func() {
		wg.Wait()
		close(tasks)
	}()

	wg2 := &sync.WaitGroup{}
	for i := 0; i < d.numWorkers; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for task := range tasks {
				result := d.downloadPage(task.URL, task.Depth)

				if result.Err != nil {
					slog.Warn("Failed to download page", "url", result.URL, "error", result.Err)
					wg.Done()
					continue
				}

				if isHTML(result.ContentType) && result.Depth < d.maxDepth {
					for _, link := range result.Links {
						resolvedURL, err := d.resolveURL(result.URL, link.URL)
						if err != nil {
							continue
						}

						if !urlutil.IsSameDomain(startURL, resolvedURL) {
							continue
						}

						if link.LinkType == "page" {
							go enqueue(resolvedURL, result.Depth+1)
						} else {
							go enqueue(resolvedURL, result.Depth)
						}
					}
				}

				wg.Done()
			}
		}()
	}

	wg2.Wait()
	return nil
}

func (d *Downloader) downloadPage(urlStr string, depth uint) *DownloadResult {
	res := &DownloadResult{
		URL:   urlStr,
		Depth: depth,
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		res.Err = fmt.Errorf("parse URL error: %w", err)
		return res
	}

	if parsedURL.Host != "" && !urlutil.IsSameDomain(d.baseURL, urlStr) {
		res.Err = fmt.Errorf("external domain: %s", parsedURL.Host)
		return res
	}

	resp, err := d.httpClient.Get(urlStr)
	if err != nil {
		res.Err = fmt.Errorf("http get error: %w", err)
		return res
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		res.Err = fmt.Errorf("bad status code: %d", resp.StatusCode)
		return res
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		res.Err = fmt.Errorf("read body error: %w", err)
		return res
	}

	var dataToSave = data
	contentType := resp.Header.Get("Content-Type")

	if isHTML(contentType) {
		parseRes, err := parser.Parse(data, nil)
		if err != nil {
			res.Err = fmt.Errorf("parse error: %w", err)
			return res
		}

		res.Links = parseRes.Links

		urlToLocalMap := d.buildURLToLocalMap(parseRes.Links, urlStr)
		parseRes, err = parser.Parse(data, urlToLocalMap)
		if err != nil {
			res.Err = fmt.Errorf("parse with replace error: %w", err)
			return res
		}

		dataToSave = parseRes.ModifiedHTML
	}

	localPath, err := d.saveFile(parsedURL, dataToSave, contentType)
	if err != nil {
		res.Err = fmt.Errorf("save file error: %w", err)
		return res
	}

	res.LocalPath = localPath
	res.ContentType = contentType

	return res
}

func isHTML(key string) bool {
	return strings.HasPrefix(key, "text/html")
}

func (d *Downloader) saveFile(u *url.URL, data []byte, contentType string) (string, error) {
	localPath := d.URLToFilePath(u, contentType)

	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return "", err
	}

	return localPath, nil
}

func (d *Downloader) URLToFilePath(u *url.URL, contentType string) string {
	path := u.Path
	if path == "" || path == "/" {
		return filepath.Join(d.outputDir, u.Host, "index.html")
	}
	if strings.HasSuffix(path, "/") {
		return filepath.Join(d.outputDir, u.Host, path, "index.html")
	}
	if !strings.Contains(strings.ToLower(path), ".") {
		ext := getExtension(contentType)
		path += ext
	}

	return filepath.Join(d.outputDir, u.Host, path)
}

func getExtension(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "text/html"):
		return ".html"
	case strings.HasPrefix(contentType, "text/css"):
		return ".css"
	case strings.HasPrefix(contentType, "application/javascript"), strings.HasPrefix(contentType, "text/javascript"):
		return ".js"
	case strings.HasPrefix(contentType, "image/jpeg"):
		return ".jpg"
	case strings.HasPrefix(contentType, "image/png"):
		return ".png"
	case strings.HasPrefix(contentType, "image/svg+xml"):
		return ".svg"
	case strings.HasPrefix(contentType, "image/gif"):
		return ".gif"
	case strings.HasPrefix(contentType, "image/webp"):
		return ".webp"
	default:
		return ".html"
	}
}

func (d *Downloader) resolveURL(base, href string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	hrefURL, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	return baseURL.ResolveReference(hrefURL).String(), nil
}

func (d *Downloader) buildURLToLocalMap(links []parser.LinkInfo, currentPageURL string) map[string]string {
	urlToLocal := make(map[string]string)

	for _, link := range links {
		if link.LinkType == "page" || link.LinkType == "css" ||
			link.LinkType == "js" || link.LinkType == "image" {
			resolvedURL, err := d.resolveURL(currentPageURL, link.URL)
			if err != nil {
				continue
			}

			if !urlutil.IsSameDomain(d.baseURL, resolvedURL) {
				continue
			}

			parsedURL, err := url.Parse(resolvedURL)
			if err != nil {
				continue
			}

			var contentType string
			switch link.LinkType {
			case "css":
				contentType = "text/css"
			case "js":
				contentType = "application/javascript"
			case "image":
				contentType = "image/jpeg"
			default:
				contentType = "text/html"
			}

			localPath := d.URLToFilePath(parsedURL, contentType)
			urlToLocal[link.URL] = localPath
		}
	}

	return urlToLocal
}
