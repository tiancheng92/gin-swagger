package gin_swagger

const swaggerIndexTemplate = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui.css" />
    <link rel="stylesheet" type="text/css" href="index.css" />
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>
    <script src="./swagger-ui-bundle.js" charset="UTF-8"> </script>
    <script src="./swagger-ui-standalone-preset.js" charset="UTF-8"> </script>
    <script src="./swagger-initializer.js" charset="UTF-8"> </script>
    <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        {{- if eq (len .URLList) 1 }}
        {{- $url := index .URLList 0 }}
		url: {{ $url.Url }},
		{{- else }}
		urls: [
			{{- range .URLList }}
		  { url:"{{ .Url }}", name: "{{ .Name }}" },
			{{- end }}
		],
		{{- end }}
		dom_id: '#swagger-ui',
		validatorUrl: null,
		oauth2RedirectUrl: {{ .Oauth2RedirectURL }},
		docExpansion: "{{ .DocExpansion }}",
		deepLinking: {{ .DeepLinking }},
		defaultModelsExpandDepth: {{ .DefaultModelsExpandDepth }},
		filter: {{ .ShowFilterTag }},
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
      window.ui = ui;
    };
  </script>
  </body>
</html>`
