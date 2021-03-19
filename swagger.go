package gin_swagger

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

//go:embed ui/*
var ui embed.FS

var (
	fileList = []string{"index.html", "doc.json", "favicon-16x16.png", "favicon-32x32.png", "swagger-ui.css", "swagger-ui.js", "swagger-ui-bundle.js", "swagger-ui-standalone-preset.js"}
	t        = template.New("swagger_index.html")
	index, _ = t.Parse(swaggerIndexTemplate)
)

type swaggerUIBundle struct {
	URL         string
	DeepLinking bool
}

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
	return func(ctx *gin.Context) {
		tmp := strings.Split(ctx.Request.RequestURI, "/")
		path := tmp[len(tmp)-1]

		if !containsString(fileList, path) {
			ctx.String(404, "")
			return
		}

		if strings.HasSuffix(path, ".html") {
			ctx.Header("Content-Type", "text/html; charset=utf-8")
		} else if strings.HasSuffix(path, ".css") {
			ctx.Header("Content-Type", "text/css; charset=utf-8")
		} else if strings.HasSuffix(path, ".js") {
			ctx.Header("Content-Type", "application/javascript")
		} else if strings.HasSuffix(path, ".json") {
			ctx.Header("Content-Type", "application/json")
		}

		switch path {
		case "index.html":
			_ = index.Execute(ctx.Writer, &swaggerUIBundle{
				URL:         config.URL,
				DeepLinking: config.DeepLinking,
			})
		case "doc.json":
			doc, err := swag.ReadDoc()
			if err != nil {
				panic(err)
			}
			_, _ = ctx.Writer.Write(stringToBytes(doc))
		default:
			f, _ := ui.ReadFile(fmt.Sprintf("ui/%s", path))
			_, _ = ctx.Writer.Write(f)
		}
	}
}

func DisablingWrapHandler(envName string) gin.HandlerFunc {
	eFlag := os.Getenv(envName)
	if eFlag != "" {
		return func(ctx *gin.Context) {
			ctx.String(404, "")
		}
	}
	return WrapHandler()
}

func DisablingCustomWrapHandler(config *Config, envName string) gin.HandlerFunc {
	eFlag := os.Getenv(envName)
	if eFlag != "" {
		return func(ctx *gin.Context) {
			ctx.String(404, "")
		}
	}
	return CustomWrapHandler(config)
}
