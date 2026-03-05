package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	authv1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/auth/v1"
	todov1 "github.com/97vaibhav/todo-crud-wire-tutorial-grpc/gen/todo/v1"
	"github.com/97vaibhav/todo-crud-wire-tutorial-grpc/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed api/todo/v1/todo.swagger.json
var todoSwagger []byte

//go:embed api/auth/v1/auth.swagger.json
var authSwagger []byte

func init() {
	merged, err := mergeSwaggerSpecs(todoSwagger, authSwagger)
	if err != nil {
		log.Fatalf("merge swagger specs: %v", err)
	}
	mergedSwaggerJSON = merged
}

var mergedSwaggerJSON []byte

// mergeSwaggerSpecs combines two OpenAPI 2.0 (Swagger) JSON specs into one.
func mergeSwaggerSpecs(a, b []byte) ([]byte, error) {
	var specA, specB map[string]interface{}
	if err := json.Unmarshal(a, &specA); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &specB); err != nil {
		return nil, err
	}
	pathsA, _ := specA["paths"].(map[string]interface{})
	pathsB, _ := specB["paths"].(map[string]interface{})
	if pathsA == nil {
		pathsA = make(map[string]interface{})
	}
	for k, v := range pathsB {
		pathsA[k] = v
	}
	specA["paths"] = pathsA
	defsA, _ := specA["definitions"].(map[string]interface{})
	defsB, _ := specB["definitions"].(map[string]interface{})
	if defsA == nil {
		defsA = make(map[string]interface{})
	}
	for k, v := range defsB {
		defsA[k] = v
	}
	specA["definitions"] = defsA
	tagsA, _ := specA["tags"].([]interface{})
	tagsB, _ := specB["tags"].([]interface{})
	if tagsA == nil {
		tagsA = []interface{}{}
	}
	specA["tags"] = append(tagsA, tagsB...)
	specA["info"] = map[string]interface{}{
		"title":   "Todo + Auth API",
		"version": "1.0",
	}
	// Add Bearer JWT security so Swagger UI shows the "Authorize" button.
	// After login, click Authorize and enter: Bearer <your_access_token>
	specA["securityDefinitions"] = map[string]interface{}{
		"bearerAuth": map[string]interface{}{
			"type":        "apiKey",
			"name":        "Authorization",
			"in":          "header",
			"description": "Enter: Bearer <access_token> (from POST /api/v1/auth/login)",
		},
	}
	specA["security"] = []interface{}{
		map[string]interface{}{"bearerAuth": []interface{}{}},
	}
	return json.Marshal(specA)
}

const swaggerUIHTML = `<!DOCTYPE html>
<html>
<head>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: '/swagger.json',
      dom_id: '#swagger-ui',
    });
  </script>
</body>
</html>
`

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Forward HTTP "Authorization" header to gRPC metadata "authorization" so the
	// backend's auth interceptor can validate the JWT.
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if strings.EqualFold(key, "Authorization") {
				return "authorization", true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := todov1.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCBackendAddr, opts); err != nil {
		log.Fatalf("register todo service gateway: %v", err)
	}
	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCBackendAddr, opts); err != nil {
		log.Fatalf("register auth service gateway: %v", err)
	}

	root := http.NewServeMux()
	root.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(mergedSwaggerJSON)
	})
	root.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(swaggerUIHTML))
	})
	root.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs", http.StatusMovedPermanently)
	})
	root.Handle("/", mux)

	addr := ":" + cfg.GatewayPort
	log.Printf("REST API gateway listening on %s (backend gRPC: %s)", addr, cfg.GRPCBackendAddr)
	log.Printf("Swagger UI: http://localhost%s/docs  spec: http://localhost%s/swagger.json", addr, addr)
	if err := http.ListenAndServe(addr, root); err != nil && err != http.ErrServerClosed {
		log.Fatalf("gateway serve: %v", err)
	}
}
