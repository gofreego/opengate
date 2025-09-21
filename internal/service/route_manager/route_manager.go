package routemanager

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gofreego/opengate/internal/models"
	"github.com/gofreego/opengate/pkg/utils"
)

type Manager interface {
	GetRoutes() []*models.ServiceRoute
	GetRouteByName(name string) *models.ServiceRoute
	GetRouteByRequest(req *http.Request) *models.ServiceRoute
	AddRoute(route *models.ServiceRoute)
	ReplaceRoutes(routes []*models.ServiceRoute)
	GetReverseProxy(targetURL string) *httputil.ReverseProxy
}

type manager struct {
	routes      []*models.ServiceRoute            // Original slice for compatibility
	trie        *utils.Trie[*models.ServiceRoute] // Trie for efficient path matching
	nameIndex   map[string]*models.ServiceRoute   // Hash map for name-based lookups
	cachedProxy map[string]*httputil.ReverseProxy // Cache of reverse proxies by target URL
}

func New() Manager {
	m := &manager{
		routes:    make([]*models.ServiceRoute, 0),
		trie:      utils.NewTrie[*models.ServiceRoute](),
		nameIndex: make(map[string]*models.ServiceRoute),
	}
	return m
}

func (m *manager) GetRoutes() []*models.ServiceRoute {
	// Return a copy to prevent external modification
	result := make([]*models.ServiceRoute, len(m.routes))
	copy(result, m.routes)
	return result
}

// Optimized O(1) lookup using hash map
func (m *manager) GetRouteByName(name string) *models.ServiceRoute {
	return m.nameIndex[name]
}

// Optimized O(m) lookup using trie where m is path length
func (m *manager) GetRouteByRequest(req *http.Request) *models.ServiceRoute {
	requestPath := req.URL.Path

	// First try trie-based lookup for exact and prefix matches
	if route := m.trie.FindLongestMatch(requestPath); route != nil {
		return route
	}

	return nil
}

func (m *manager) AddRoute(route *models.ServiceRoute) {
	// Add to routes slice
	m.routes = append(m.routes, route)

	// Add to name index
	m.nameIndex[route.Name] = route

	// Add to trie
	m.trie.Insert(route.PathPrefix, route)
}

func (m *manager) ReplaceRoutes(routes []*models.ServiceRoute) {
	trie := buildTrie(routes)
	nameIndex := make(map[string]*models.ServiceRoute)
	for _, route := range routes {
		nameIndex[route.Name] = route
	}

	m.trie = trie
	m.routes = routes
	m.nameIndex = nameIndex
}

func (m *manager) GetReverseProxy(targetURL string) *httputil.ReverseProxy {
	// Check cache first
	if proxy, exists := m.cachedProxy[targetURL]; exists {
		return proxy
	}

	// Create new reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: targetURL})

	// Initialize cache map if nil
	if m.cachedProxy == nil {
		m.cachedProxy = make(map[string]*httputil.ReverseProxy)
	}

	// Cache the proxy
	m.cachedProxy[targetURL] = proxy
	return proxy
}

// Helper method to rebuild trie
func buildTrie(routes []*models.ServiceRoute) *utils.Trie[*models.ServiceRoute] {
	trie := utils.NewTrie[*models.ServiceRoute]()
	for _, route := range routes {
		trie.Insert(route.PathPrefix, route)
	}
	return trie
}
