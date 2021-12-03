package gin_swagger

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

//go:embed ui/*
var ui embed.FS

var (
	fileList = []string{"index.html", "doc.json", "favicon-16x16.png", "favicon-32x32.png", "swagger-ui.css", "swagger-ui.js", "swagger-ui-bundle.js", "swagger-ui-standalone-preset.js", "oauth2-redirect.html"}
	t        = template.New("swagger_index.html")
	index, _ = t.Parse(swaggerIndexTemplate)
)

type swaggerUIBundle struct {
	URLList                  []SwaggerURL
	DeepLinking              bool
	DocExpansion             string
	DefaultModelsExpandDepth int
	Oauth2RedirectURL        template.JS
	Title                    string
	ShowFilterTag            bool
}

type SwaggerURL struct {
	Url  string
	Name string
}

type Config struct {
	URLList                  []SwaggerURL
	DeepLinking              bool
	DocExpansion             string
	DefaultModelsExpandDepth int
	Title                    string
	ShowFilterTag            bool
}

func URL(urls ...SwaggerURL) func(c *Config) {
	return func(c *Config) {
		c.URLList = urls
	}
}

func DocExpansion(docExpansion string) func(c *Config) {
	return func(c *Config) {
		c.DocExpansion = docExpansion
	}
}

func Title(title string) func(c *Config) {
	return func(c *Config) {
		c.Title = title
	}
}

func DeepLinking(deepLinking bool) func(c *Config) {
	return func(c *Config) {
		c.DeepLinking = deepLinking
	}
}

func ShowFilterTag(showFilterTag bool) func(c *Config) {
	return func(c *Config) {
		c.ShowFilterTag = showFilterTag
	}
}

func DefaultModelsExpandDepth(depth int) func(c *Config) {
	return func(c *Config) {
		c.DefaultModelsExpandDepth = depth
	}
}

func WrapHandler(configs ...func(c *Config)) gin.HandlerFunc {
	defaultConfig := &Config{
		URLList:                  []SwaggerURL{{Url: "doc.json", Name: "Default"}},
		DeepLinking:              true,
		DocExpansion:             "list",
		DefaultModelsExpandDepth: 1,
		Title:                    "Swagger UI",
		ShowFilterTag:            false,
	}

	for _, c := range configs {
		c(defaultConfig)
	}

	return CustomWrapHandler(defaultConfig)
}

func CustomWrapHandler(config *Config) gin.HandlerFunc {
	if config.Title == "" {
		config.Title = "Swagger UI"
	}

	return func(ctx *gin.Context) {
		tmp := strings.Split(ctx.Request.RequestURI, "/")
		path := tmp[len(tmp)-1]

		if !containsString(fileList, path) {
			ctx.String(http.StatusNotFound, "")
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
				URLList:                  config.URLList,
				DeepLinking:              config.DeepLinking,
				DocExpansion:             config.DocExpansion,
				DefaultModelsExpandDepth: config.DefaultModelsExpandDepth,
				Oauth2RedirectURL: template.JS(
					"`${window.location.protocol}//${window.location.host}$" +
						"{window.location.pathname.split('/').slice(0, window.location.pathname.split('/').length - 1).join('/')}" +
						"/oauth2-redirect.html`",
				),
				Title:         config.Title,
				ShowFilterTag: config.ShowFilterTag,
			})
		case "doc.json":
			doc, err := swag.ReadDoc()
			if err != nil {
				panic(err)
			}
			_, _ = ctx.Writer.Write(stringToBytes(doc))
		default:
			f, _ := ui.ReadFile(strings.Join([]string{"ui", path}, "/"))
			_, _ = ctx.Writer.Write(f)
		}
	}
}

func DisablingWrapHandler(envName string) gin.HandlerFunc {
	eFlag := os.Getenv(envName)
	if eFlag != "" {
		return func(ctx *gin.Context) {
			ctx.String(http.StatusNotFound, "")
		}
	}
	return WrapHandler()
}

func DisablingCustomWrapHandler(config *Config, envName string) gin.HandlerFunc {
	eFlag := os.Getenv(envName)
	if eFlag != "" {
		return func(ctx *gin.Context) {
			ctx.String(http.StatusNotFound, "")
		}
	}
	return CustomWrapHandler(config)
}
