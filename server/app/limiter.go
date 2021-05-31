package app

import (
	"context"
)

func (s *Service) AllowRead(ctx context.Context, username string) bool {
	lim := s.factory.TryGetRateLimiter(ctx, username)

	if lim == nil {
		return false
	}

	return lim.AllowRead()
}

func (s *Service) AllowWrite(ctx context.Context, username string) bool {
	lim := s.factory.TryGetRateLimiter(ctx, username)

	if lim == nil {
		return false
	}

	return lim.AllowWrite()
}
