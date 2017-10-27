package web

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WebPage wraps necessary data to make a POST or GET http request,
// and reads its reponse. It provides additional functionalities to
// assist parsing the results.
type WebPage struct {
	// request
	Url             string
	Method          string // "POST" or "GET" supported
	Timeout         time.Duration
	User            string // support for basic auth (not encrypted)
	Password        string
	FollowRedirects string // must set it to 'no' to disallow following redirects; otherwise, it use default settings
	// may set compression and other options here
	//	Header = map[string][]string{
	//		"Accept-Encoding": {"gzip, deflate"},
	//		"Accept-Language": {"en-us"},
	//		"Connection": {"keep-alive"},
	//	}
	Header   http.Header
	BodyType string    // e.g. "text/xml", "application/json", "application/x-www-form-urlencoded"
	Body     io.Reader // the data sent when using POST method

	//file upload
	File      io.ReadCloser     //file to upload.
	FieldName string            //field name of the file.
	FileName  string            //file name
	Fields    map[string]string //params.

	// response
	StatusCode         int    // response status code, client should check it to decide the next move
	Status             string // the status message returned
	Cookies            []*http.Cookie
	RespHeader         http.Header
	RespReader         io.ReadCloser //try to recieve the cons body as stream/reader.
	RespBody           []byte        // try to read even if stauts code is not 200, since some web server returns error message which might be useful
	InsecureSkipVerify bool
}

// DoRequest wraps the common http request paradigm to get the request result.
func (w *WebPage) DoRequest() error {
	// no need to check parameters, NewRequest would do it.
	req, err := http.NewRequest(w.Method, w.Url, w.Body)
	if err != nil {
		return err
	}
	if w.Header != nil {
		req.Header = w.Header
	}
	if w.User != "" && w.Password != "" {
		req.SetBasicAuth(w.User, w.Password)
	}
	// add body type
	if w.BodyType != "" {
		req.Header.Set("Content-Type", w.BodyType)
	}
	var resp *http.Response
	// do not follow redirects?
	if w.FollowRedirects == "no" {
		if w.Timeout > 0 {
			tr := &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					c, err := net.DialTimeout(netw, addr, w.Timeout)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(time.Now().Add(w.Timeout))
					return c, nil
				},
			}
			resp, err = tr.RoundTrip(req)
			if err != nil {
				return err
			}
		} else {
			tr := &http.Transport{}
			if w.InsecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
			resp, err = tr.RoundTrip(req)
			if err != nil {
				return err
			}
		}
	} else {
		var client *http.Client
		client = http.DefaultClient
		if w.Timeout > 0 {
			client = &http.Client{Timeout: w.Timeout}
		}
		if w.InsecureSkipVerify {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client.Transport = tr
		}
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	defer resp.Body.Close()
	w.Cookies = resp.Cookies()
	w.RespHeader = resp.Header
	w.StatusCode, w.Status = resp.StatusCode, resp.Status
	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	w.RespBody = page
	return nil
}

//upload multipart file.
func (w *WebPage) PostFile() error {
	if w.File == nil {
		return fmt.Errorf("invalid multipart file.")
	}
	if strings.TrimSpace(w.FieldName) == "" {
		return fmt.Errorf("FieldName can not be empty.")
	}
	if strings.TrimSpace(w.FileName) == "" {
		return fmt.Errorf("FileName can not be empty.")
	}
	rd, wt := io.Pipe()
	defer rd.Close()
	mpw := multipart.NewWriter(wt)

	errChan := make(chan error, 1)
	go func() {
		var part io.Writer
		var err error

		if w.Fields != nil {
			for k, v := range w.Fields {
				if err = mpw.WriteField(k, v); err != nil {
					errChan <- err
					return
				}
			}
		}
		if part, err = mpw.CreateFormFile(w.FieldName, w.FileName); err != nil {
			errChan <- err
			return
		}
		if _, err = io.Copy(part, w.File); err != nil {
			errChan <- err
			return
		}

		if err = mpw.Close(); err != nil {
			errChan <- err
			return
		}
		if err = wt.Close(); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()
	w.BodyType = mpw.FormDataContentType()
	w.Body = rd

	if err := w.DoRequestFile(); err != nil {
		return err
	}

	return <-errChan
}

func (ww *WebPage) PostMultiformSync(files map[string]string) error {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	var fw io.Writer
	for fileName, filePath := range files {
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		fw, err = w.CreateFormFile(fileName, fileName)
		if err != nil {
			return err
		}
		if _, err = io.Copy(fw, f); err != nil {
			return err
		}
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	ww.BodyType = w.FormDataContentType()
	// Don't forget to set the content type, this will contain the boundary.
	ww.Body = bytes.NewReader(b.Bytes())
	if err := ww.DoRequest(); err != nil {
		return nil
	}
	return nil
}

//upload multipart files.
// files: map fieldName to fileName, support one file per field.
func (w *WebPage) PostMultiform(files map[string]string, form map[string]string) error {
	// open files
	var flist []*os.File
	defer func() {
		for i, _ := range flist {
			if flist[i] != nil {
				flist[i].Close()
			}
		}
	}()
	for field, file := range files {
		if strings.TrimSpace(field) == "" {
			return fmt.Errorf("bad field name: %v", field)
		}
		if f, err := os.Open(file); err != nil {
			return err
		} else {
			flist = append(flist, f)
		}
	}
	// open pipe
	rd, wt := io.Pipe()
	defer rd.Close()
	mpw := multipart.NewWriter(wt)
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if err := mpw.Close(); err != nil {
				fmt.Printf(err.Error())
			}
			if err := wt.Close(); err != nil {
				fmt.Printf(err.Error())
			}
		}()
		var part io.Writer
		var err error
		// add form values
		if form != nil {
			for k, v := range form {
				if err = mpw.WriteField(k, v); err != nil {
					errChan <- err
					return
				}
			}
		}
		// add files
		if files != nil {
			idx := 0
			for field, file := range files {
				if part, err = mpw.CreateFormFile(field, filepath.Base(file)); err != nil {
					errChan <- err
					return
				}
				if _, err = io.Copy(part, flist[idx]); err != nil {
					errChan <- err
					return
				}
				idx++
			}
		}
		errChan <- nil
	}()

	w.BodyType = mpw.FormDataContentType()
	w.Body = rd

	if err := w.DoRequest(); err != nil {
		return err
	}

	return <-errChan
}

//upload multipart files.
// files: map field name to file names, support multiple files per field.
func (w *WebPage) PostMultiforms(files map[string]([]string), form map[string]string) error {
	// open files
	var flist []*os.File
	defer func() {
		for i, _ := range flist {
			if flist[i] != nil {
				flist[i].Close()
			}
		}
	}()
	for field, files := range files {
		if strings.TrimSpace(field) == "" {
			return fmt.Errorf("bad field name: %v", field)
		}
		for _, onefile := range files {
			if f, err := os.Open(onefile); err != nil {
				return err
			} else {
				flist = append(flist, f)
			}
		}
	}
	// open pipe
	rd, wt := io.Pipe()
	defer rd.Close()
	mpw := multipart.NewWriter(wt)
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			if err := mpw.Close(); err != nil {
				fmt.Printf(err.Error())
			}
			if err := wt.Close(); err != nil {
				fmt.Printf(err.Error())
			}
		}()
		var part io.Writer
		var err error
		// add form values
		if form != nil {
			for k, v := range form {
				if err = mpw.WriteField(k, v); err != nil {
					errChan <- err
					return
				}
			}
		}
		// add files
		if files != nil {
			idx := 0
			for field, files := range files {
				for _, onefile := range files {
					if part, err = mpw.CreateFormFile(field, filepath.Base(onefile)); err != nil {
						errChan <- err
						return
					}
					if _, err = io.Copy(part, flist[idx]); err != nil {
						errChan <- err
						return
					}
					idx++
				}
			}
		}
		errChan <- nil
	}()

	w.BodyType = mpw.FormDataContentType()
	w.Body = rd

	if err := w.DoRequest(); err != nil {
		return err
	}

	return <-errChan
}

// GetRespBody wraps the common http request paradigm to get the cons body,and you need close it by yourself.
func (w *WebPage) DoRequestFile() error {
	// no need to check parameters, NewRequest would do it.
	req, err := http.NewRequest(w.Method, w.Url, w.Body)
	if err != nil {
		return err
	}
	if w.Header != nil {
		req.Header = w.Header
	}
	if w.User != "" && w.Password != "" {
		req.SetBasicAuth(w.User, w.Password)
	}
	// add body type
	if w.BodyType != "" {
		req.Header.Set("Content-Type", w.BodyType)
	}
	var resp *http.Response
	// do not follow redirects?
	if w.FollowRedirects == "no" {
		if w.Timeout > 0 {
			tr := &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					c, err := net.DialTimeout(netw, addr, w.Timeout)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(time.Now().Add(w.Timeout))
					return c, nil
				},
			}
			resp, err = tr.RoundTrip(req)
			if err != nil {
				return err
			}
		} else {
			tr := &http.Transport{}
			resp, err = tr.RoundTrip(req)
			if err != nil {
				return err
			}
		}
	} else {
		var client *http.Client
		client = http.DefaultClient
		if w.Timeout > 0 {
			client = &http.Client{Timeout: w.Timeout}
		}
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	w.Cookies = resp.Cookies()
	w.RespHeader = resp.Header
	w.StatusCode, w.Status = resp.StatusCode, resp.Status
	w.RespReader = resp.Body
	return nil
}

// DoRequestWithProxy wraps the common http request paradigm to get the request result via a http proxy.
func (w *WebPage) DoRequestWithProxy(proxy_uri string) error {
	// no need to check parameters, NewRequest would do it.
	req, err := http.NewRequest(w.Method, w.Url, w.Body)
	if err != nil {
		return err
	}
	if w.Header != nil {
		req.Header = w.Header
	}
	if w.User != "" && w.Password != "" {
		req.SetBasicAuth(w.User, w.Password)
	}
	// add body type
	if w.BodyType != "" {
		req.Header.Set("Content-Type", w.BodyType)
	}
	proxy, err := url.Parse(proxy_uri)
	if err != nil {
		return err
	}
	var resp *http.Response
	if w.FollowRedirects == "no" {
		if w.Timeout > 0 {
			tr := &http.Transport{
				Proxy: http.ProxyURL(proxy),
				Dial: func(netw, addr string) (net.Conn, error) {
					c, err := net.DialTimeout(netw, addr, w.Timeout)
					if err != nil {
						return nil, err
					}
					if w.Timeout > 0 {
						c.SetDeadline(time.Now().Add(w.Timeout))
					}
					return c, nil
				},
			}
			resp, err = tr.RoundTrip(req)
			if err != nil {
				return err
			}
		} else {
			tr := &http.Transport{Proxy: http.ProxyURL(proxy)}
			resp, err = tr.RoundTrip(req)
			if err != nil {
				return err
			}
		}
	} else {
		timeout := http.DefaultClient.Timeout
		if w.Timeout > 0 {
			timeout = w.Timeout
		}
		httpClient := &http.Client{Timeout: timeout, Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
		resp, err = httpClient.Do(req)
		if err != nil {
			return err
		}
	}
	defer resp.Body.Close()
	w.Cookies = resp.Cookies()
	w.RespHeader = resp.Header
	w.StatusCode, w.Status = resp.StatusCode, resp.Status
	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	w.RespBody = page
	return nil
}

func (w *WebPage) JsonUnmarshal(entity interface{}) error {
	if w.RespBody == nil {
		fmt.Errorf("Webpage: can not unmarshal nil response body.")
	}
	return json.Unmarshal(w.RespBody, entity)
}

func (w *WebPage) AddHeaders(headers http.Header) {
	// copy BasicHeader from config
	h := http.Header{}
	for k, v := range headers {
		h.Add(k, v[0])
	}
	w.Header = h
}

// could overwrite existing headers
func (w *WebPage) SetHeader(key string, value string) {
	if w.Header == nil {
		w.Header = http.Header{}
	}
	w.Header.Set(key, value)
}

//get file name from cons header.
func (w *WebPage) GetMulitpartFileName() (string, error) {
	_, params, err := mime.ParseMediaType(w.RespHeader.Get("Content-Disposition"))
	if err != nil {
		return "", err
	}
	return params["filename"], nil
}

// check request succeed or not by http status code.
func (w *WebPage) CheckErr(prefix string) error {
	if w.StatusCode >= 300 {
		return fmt.Errorf(prefix+" request failed. status: %d,msg: %s", w.StatusCode, string(w.RespBody))
	}
	return nil
}

// check request succeed or not by http status code and cons.Code
func (w *WebPage) CheckErrPro(prefix string) error {
	if w.StatusCode >= 300 {
		return fmt.Errorf(prefix+" request failed. status: %d,msg: %s", w.StatusCode, string(w.RespBody))
	}
	resp := struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}{}
	if err := json.Unmarshal(w.RespBody, &resp); err != nil {
		return err
	}
	if resp.Code != 10200 {
		return fmt.Errorf(prefix+" request failed, response status: %d, msg: %s", resp.Code, resp.Message)
	}
	return nil
}

func (w *WebPage) SetTimeout(timeout time.Duration) {
	w.Timeout = timeout
}

func (w *WebPage) SetInsecureSkipVerify() {
	w.InsecureSkipVerify = true
}
