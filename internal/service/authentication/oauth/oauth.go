package authentication

type Config struct {
}

type OAuthAuthenticationService struct {
	// cfg Config
}

func NewOAuthAuthenticationService(config Config) *OAuthAuthenticationService {
	panic("OAuthAuthenticationService unimplemented")
	// return &OAuthAuthenticationService{cfg: config}
}

func (s *OAuthAuthenticationService) Authenticate() error {
	return nil
}
