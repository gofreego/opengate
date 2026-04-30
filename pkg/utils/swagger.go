package utils

import (
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func RegisterSwaggerHandler(mux *runtime.ServeMux, swaggerPath, swaggerDir, defaultJsonFile string) {
	// Serve Swagger specification files (JSON or YAML)
	swaggerFS := http.FileServer(http.Dir(swaggerDir))
	mux.HandlePath("GET", swaggerPath+"/**", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.StripPrefix(swaggerPath, swaggerFS).ServeHTTP(w, r)
	})

	// Serve a simple HTML page that references the Swagger UI from the official CDN
	mux.HandlePath("GET", swaggerPath, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
            <!DOCTYPE html>
            <html>
            <head>
                <title>OpenGate API - Swagger UI</title>
                <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.1/swagger-ui.css" />
            </head>
            <body>
                <div id="swagger-ui"></div>
                <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.1/swagger-ui-bundle.js"></script>
                <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.1/swagger-ui-standalone-preset.js"></script>
                <script>
                    window.onload = function() {
                        SwaggerUIBundle({
                            url: "`+swaggerPath+defaultJsonFile+`",
                            dom_id: '#swagger-ui',
                            presets: [
                                SwaggerUIBundle.presets.apis,
                                SwaggerUIStandalonePreset
                            ],
                            layout: "StandaloneLayout"
                        });
                    };
                </script>
            </body>
            </html>
        `)
	})
}
