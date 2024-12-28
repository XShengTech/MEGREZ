package request

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Client 		==> 客户端实例
type Client struct {
	Request *Request
	Cookies []*http.Cookie
	Result  Result
}

// Request 		==> 请求体
type Request struct {
	Url           string
	Method        string
	Data          io.Reader
	ContentType   string
	Authorization string
	UserAgent     string
	Header        map[string]string
	Timeout       time.Duration
	// The proxy type is determined by the URL scheme. "http",
	// "https", and "socks5" are supported. If the scheme is empty,
	//
	// If Proxy is nil or nil *URL, no proxy is used.
	ProxyUrl     url.URL
	Redirect     bool
	MaxRedirects int
}

// Result 		==> 结果集
type Result struct {
	Header   http.Header
	Location *url.URL
	Body     []byte
	Status   int
}

// NewRequest 		==> 新建请求
func NewRequest() *Client {
	return &Client{
		Request: &Request{
			Method:       "GET",
			UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			Header:       make(map[string]string),
			Redirect:     false,
			MaxRedirects: 10,
		},
		Result: Result{},
	}
}

// Do 		==> 执行请求
func (c *Client) Do() *Client {
	//HTTP请求构造
	request, _ := http.NewRequest(c.Request.Method, c.Request.Url, c.Request.Data)
	request.Header.Set("Content-Type", c.Request.ContentType)
	if c.Request.Authorization != "" {
		request.Header.Set("Authorization", c.Request.Authorization)
	}
	if c.Request.UserAgent != "" {
		request.Header.Set("User-Agent", c.Request.UserAgent)
	}
	if len(c.Cookies) != 0 {
		for _, cookie := range c.Cookies {
			request.AddCookie(cookie)
		}
	}
	// 支持自定义Header
	for k, v := range c.Request.Header {
		request.Header.Set(k, v)
	}

	var client *http.Client
	if c.Request.ProxyUrl == (url.URL{}) {
		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if !c.Request.Redirect {
					return http.ErrUseLastResponse
				}
				if len(via) >= c.Request.MaxRedirects {
					return http.ErrUseLastResponse
				}
				return nil
			},
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(&c.Request.ProxyUrl)},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if !c.Request.Redirect {
					return http.ErrUseLastResponse
				}
				if len(via) >= c.Request.MaxRedirects {
					return http.ErrUseLastResponse
				}
				return nil
			},
		}
	}
	if c.Request.Timeout != 0 {
		client.Timeout = c.Request.Timeout
	}
	res, err := client.Do(request)
	if err != nil {
		// fmt.Println(err)
		return c
	}
	if len(res.Cookies()) > 1 {
		c.Cookies = append(c.Cookies, res.Cookies()...)
	}
	defer res.Body.Close()
	c.Result.Status = res.StatusCode
	c.Result.Body, _ = io.ReadAll(res.Body)
	c.Result.Header = res.Header
	return c
}

// Get 		==> 定义请求方式
func (c *Client) Get() *Client {
	c.Request.Method = "GET"
	return c
}

// Post 		==> 定义请求方式
func (c *Client) Post() *Client {
	c.Request.Method = "POST"
	return c
}

// Put 		==> 定义请求方式
func (c *Client) Put() *Client {
	c.Request.Method = "PUT"
	return c
}

// Patch 		==> 定义请求方式
func (c *Client) Patch() *Client {
	c.Request.Method = "PATCH"
	return c
}

// Delete 		==> 定义请求方式
func (c *Client) Delete() *Client {
	c.Request.Method = "DELETE"
	return c
}

// SetUrl 		==> 定义请求目标
func (c *Client) SetUrl(url ...any) *Client {
	c.Request.Url = fmt.Sprintf(url[0].(string), url[1:]...)
	return c
}

// SetMethod 		==> 定义请求方法
func (c *Client) SetMethod(method string) *Client {
	c.Request.Method = method
	return c
}

// SetContentType 		==> 定义内容类型
func (c *Client) SetContentType(contentType string) *Client {
	c.Request.ContentType = contentType
	return c
}

// SetUserAgent 		==> 定义用户代理
func (c *Client) SetUserAgent(userAgent string) *Client {
	c.Request.UserAgent = userAgent
	return c
}

// SetBody 		==> 定义请求内容
func (c *Client) SetBody(body io.Reader) *Client {
	c.Request.Data = body
	return c
}

// SerHeaders 		==> 定义请求头列表
func (c *Client) SetHeaders(headers map[string]string) *Client {
	c.Request.Header = headers
	return c
}

// SetHeader 		==> 定义请求头
func (c *Client) SetHeader(key, value string) *Client {
	c.Request.Header[key] = value
	return c
}

// SetAuthorization 		==> 定义身份验证
func (c *Client) SetAuthorization(credentials string) *Client {
	c.Request.Authorization = credentials
	return c
}

// SetTimeOut 		==> 设置会话超时上限
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.Request.Timeout = timeout
	return c
}

// SetCookie 		==> 设置Cookie
func (c *Client) SetCookie(cookie *http.Cookie) *Client {
	c.Cookies = append(c.Cookies, cookie)
	return c
}

// SetCookies 		==> 设置Cookies
func (c *Client) SetCookies(cookies string) *Client {
	cookielist := strings.Split(cookies, "; ")
	for _, cookie := range cookielist {
		cookiekv := strings.Split(cookie, "=")
		c.SetCookie(&http.Cookie{
			Name:  cookiekv[0],
			Value: strings.Join(cookiekv[1:], "="),
		})
	}
	return c
}

// SetProxy 		==> 设置代理
func (c *Client) SetProxy(proxyUrl url.URL) *Client {
	c.Request.ProxyUrl = proxyUrl
	return c
}

// SetRedirect 		==> 设置重定向
func (c *Client) SetRedirect(redirect bool) *Client {
	c.Request.Redirect = redirect
	return c
}

// SetMaxRedirects 		==> 设置最大重定向次数
func (c *Client) SetMaxRedirects(maxRedirects int) *Client {
	c.Request.MaxRedirects = maxRedirects
	return c
}

// GetStatusCode 		==> 获取请求状态码
func (c *Client) GetStatusCode() int {
	return c.Result.Status
}

// GetBody 		==> 获取返回内容
func (c *Client) GetBody() []byte {
	return c.Result.Body
}

// GetBody 		==> 获取返回内容
func (c *Client) GetBodyString() string {
	return string(c.Result.Body)
}

// GetHeaders		==> 获取返回头字典
func (c *Client) GetHeaders() http.Header {
	return c.Result.Header
}

// GetHeader 		==> 获取返回头
func (c *Client) GetHeader(key string) string {
	return c.Result.Header.Get(key)
}

// SaveToFile 		==> 写出结果到文件
func (c *Client) SaveToFile(filepath string) (err error) {
	// Write the body to file
	err = os.WriteFile(filepath, c.GetBody(), 0777)
	if err != nil {
		return err
	}
	return nil
}
