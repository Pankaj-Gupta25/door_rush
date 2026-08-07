package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Param("proxyPath") // This assumes we use a wildcard param like /*proxyPath
		}

		// Custom error handler for the proxy
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable: " + err.Error()})
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// SimpleProxyHandler proxies requests preserving the original path structure or appending it to the target
func SimpleProxyHandler(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			// Update the headers to allow for SSL redirection
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.Host = remote.Host

			// Optional: Trim a prefix if the target service doesn't expect it
			// req.URL.Path = strings.TrimPrefix(req.URL.Path, "/prefix")
		}

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			// Check if the error is "connection refused"
			c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable: " + err.Error()})
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
