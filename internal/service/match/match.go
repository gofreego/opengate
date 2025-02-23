package match

import (
	"api-gateway/internal/models/dao"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/customerrors"
)

type MatchService struct {
	//routes with host configured
	allHostMatchers map[string]map[string]matchers
	// routes where host is not cofigured id -> matchers
	otherMatchers map[string]matchers
}

func NewMatchService() *MatchService {
	return &MatchService{
		allHostMatchers: make(map[string]map[string]matchers),
		otherMatchers:   make(map[string]matchers),
	}
}

func (s *MatchService) UpdateRoutesMatch(ctx context.Context, matches ...*dao.RouteConfig) error {
	for _, m := range matches {
		// updating for host matcher
		if m.Match.Host != "" {
			if hostMatchers, found := s.allHostMatchers[m.Match.Host]; found {
				hostMatchers[m.ID] = NewMatchers(m.ID, s.getMatchers(m))
			} else {
				s.allHostMatchers[m.Match.Host] = make(map[string]matchers)
				s.allHostMatchers[m.Match.Host][m.ID] = NewMatchers(m.ID, s.getMatchers(m))
			}
			delete(s.otherMatchers, m.ID)
			continue
		}

		// updating for others where host is configured
		s.otherMatchers[m.ID] = NewMatchers(m.ID, s.getMatchers(m))
		//deleting from host matchers
		for host, hostMatchers := range s.allHostMatchers {
			if _, found := hostMatchers[m.ID]; found {
				delete(hostMatchers, m.ID)
				if len(hostMatchers) == 0 {
					delete(s.allHostMatchers, host)
				}
			}
		}
	}
	return nil
}

func (s *MatchService) GetMatchID(ctx *gin.Context) (string, error) {

	//matching with host matchers
	hostMatchers := s.allHostMatchers[ctx.Request.Host]
	for _, match := range hostMatchers {
		if match.isMatching(ctx) {
			return match.id, nil
		}
	}

	// matching with other matchers
	for id, match := range s.otherMatchers {
		if match.isMatching(ctx) {
			return id, nil
		}
	}

	return "", customerrors.New(http.StatusBadGateway, "bad gateway")
}
