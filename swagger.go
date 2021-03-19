package gin_swagger

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

//go:embed ui/*
var ui embed.FS

type Config struct {
	URL         string
	DeepLinking bool
}

func URL(url string) func(c *Config) {
	return func(c *Config) {
		c.URL = url
	}
}

func DeepLinking(deepLinking bool) func(c *Config) {
	return func(c *Config) {
		c.DeepLinking = deepLinking
	}
}

func WrapHandler(configs ...func(c *Config)) gin.HandlerFunc {
	defaultConfig := &Config{
		URL:         "doc.json",
		DeepLinking: true,
	}

	for _, c := range configs {
		c(defaultConfig)
	}

	return CustomWrapHandler(defaultConfig)
}

func CustomWrapHandler(config *Config) gin.HandlerFunc {
	t := template.New("swagger_index.html")
	index, _ := t.Parse(swaggerIndexTemplate)

	var regular = regexp.MustCompile(`(.*)(index\.html|doc\.json|favicon-16x16\.png|favicon-32x32\.png|swagger-ui\.css|swagger-ui\.js|swagger-ui-bundle\.js|swagger-ui-standalone-preset\.js)[|.]*`)
	return func(c *gin.Context) {
		type swaggerUIBundle struct {
			URL         string
			DeepLinking bool
		}

		var matches []string
		if matches = regular.FindStringSubmatch(c.Request.RequestURI); len(matches) != 3 {
			c.Status(404)
			_, _ = c.Writer.Write([]byte("404 page not found"))
			return
		}
		path := matches[2]

		if strings.HasSuffix(path, ".html") {
			c.Header("Content-Type", "text/html; charset=utf-8")
		} else if strings.HasSuffix(path, ".css") {
			c.Header("Content-Type", "text/css; charset=utf-8")
		} else if strings.HasSuffix(path, ".js") {
			c.Header("Content-Type", "application/javascript")
		} else if strings.HasSuffix(path, ".json") {
			c.Header("Content-Type", "application/json")
		}

		switch path {
		case "index.html":
			_ = index.Execute(c.Writer, &swaggerUIBundle{
				URL:         config.URL,
				DeepLinking: config.DeepLinking,
			})
		case "doc.json":
			doc, err := swag.ReadDoc()
			if err != nil {
				panic(err)
			}
			_, _ = c.Writer.Write([]byte(doc))
			return
		default:
			f, _ := ui.ReadFile(fmt.Sprintf("ui/%s", path))
			_, _ = c.Writer.Write(f)
		}
	}
}

func DisablingWrapHandler(envName string) gin.HandlerFunc {
	eFlag := os.Getenv(envName)
	if eFlag != "" {
		return func(c *gin.Context) {
			c.String(404, "")
		}
	}
	return WrapHandler()
}

func DisablingCustomWrapHandler(config *Config, envName string) gin.HandlerFunc {
	eFlag := os.Getenv(envName)
	if eFlag != "" {
		return func(c *gin.Context) {
			c.String(404, "")
		}
	}
	return CustomWrapHandler(config)
}
