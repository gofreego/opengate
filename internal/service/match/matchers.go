package match

import (
	"api-gateway/internal/models/dao"
	"context"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/logger"
)

type Matcher func(c *gin.Context) bool

type matchers struct {
	id       string
	matchers []Matcher
}

func NewMatchers(id string, ms []Matcher) matchers {
	return matchers{
		id:       id,
		matchers: ms,
	}
}

func (m *matchers) isMatching(ctx *gin.Context) bool {
	for _, match := range m.matchers {
		if !match(ctx) {
			return false
		}
	}
	return true
}

func (s *MatchService) getMatchers(ctx context.Context, cfg *dao.RouteConfig) []Matcher {

	var matchers []Matcher

	if cfg.Match.Regex != "" {
		regex, err := regexp.Compile(cfg.Match.Regex)
		if err != nil {
			logger.Error(ctx, "error while compiling regex for ID - %s : Err: %s ", cfg.ID, err.Error())
			return nil
		}
		regexMatcher := func(ctx *gin.Context) bool {
			match := regex.MatchString(ctx.Request.URL.String())
			return match
		}
		matchers = append(matchers, regexMatcher)
	}

	if cfg.Match.Prefix != "" {
		prefixMatcher := func(ctx *gin.Context) bool {
			return strings.HasPrefix(ctx.Request.URL.String(), cfg.Match.Prefix)
		}
		matchers = append(matchers, prefixMatcher)
	}

	if len(cfg.Match.Methods) != 0 {
		methodMatcher := func(c *gin.Context) bool {
			methodAllowed := false
			for _, m := range cfg.Match.Methods {
				if c.Request.Method == m {
					methodAllowed = true
					break
				}
			}
			return methodAllowed
		}
		matchers = append(matchers, methodMatcher)
	}

	return matchers
}
