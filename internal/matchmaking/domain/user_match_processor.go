package matchdomain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	domain_error "github.com/xfrr/randomtalk/internal/shared/domain-error"
)

var _ MatchmakingProcessor = (*UserMatchProcessor)(nil)

// NotificationsChannel defines push-based notification behavior.
type NotificationsChannel interface {
	Notify(ctx context.Context, userID string, match *Match) error
}

// UserMatchProcessor is the concrete matchmaking service.
type UserMatchProcessor struct {
	notifications   NotificationsChannel
	matchRepository MatchRepository
	userStore       UserStore
	matcher         StableMatchFinder
	logger          *zerolog.Logger
}

// UserMatchMakerOption defines a functional option to configure the UserMatchMaker.
type UserMatchMakerOption func(*UserMatchProcessor)

// WithLogger overrides the default zerolog.Logger.
func WithLogger(logger *zerolog.Logger) UserMatchMakerOption {
	return func(s *UserMatchProcessor) {
		s.logger = logger
	}
}

// NewUserMatchProcessor initializes a new UserMatchProcessor.
func NewUserMatchProcessor(
	matchRepo MatchRepository,
	userStore UserStore,
	matcher StableMatchFinder,
	notifications NotificationsChannel,
	opts ...UserMatchMakerOption,
) (*UserMatchProcessor, error) {
	svc := &UserMatchProcessor{
		notifications:   notifications,
		matchRepository: matchRepo,
		matcher:         matcher,
		logger:          &zerolog.Logger{},
		userStore:       userStore,
	}

	for _, opt := range opts {
		opt(svc)
	}

	if err := svc.ensureDependencies(); err != nil {
		return nil, err
	}

	svc.setupLogger()
	return svc, nil
}

// ProcessMatchRequest attempts to match a set of users and enqueues them if no match is found.
func (svc *UserMatchProcessor) ProcessMatchRequest(ctx context.Context, user User) error {
	svc.logger.Debug().
		Str("user_id", user.ID()).
		Msg("processing match user request")

	// try to attempt a match immediately
	err := svc.attemptMatch(ctx, &user)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoActiveUsers):
			if err = svc.userStore.AddUser(ctx, user); err != nil {
				return fmt.Errorf("ading user to user store: %w", err)
			}

			svc.logger.Debug().
				Str("user_id", user.ID()).
				Msg("no active users, user added to store for later matching")
			return nil
		default:
			return fmt.Errorf("failed to attempt match: %w", err)
		}
	}
	return nil
}

func (svc *UserMatchProcessor) attemptMatch(ctx context.Context, candidates ...*User) error {
	activeUsers, err := svc.userStore.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all active users: %w", err)
	}

	idxCol := svc.matcher.FindStableMatches(candidates, activeUsers)
	if len(idxCol) == 0 {
		return ErrNoActiveUsers
	}

	// process the match
	for idxsa, idxsb := range idxCol {
		if idxsb == -1 {
			// no match found for this user
			continue
		}

		candidate := candidates[idxsa]
		matchedUser := activeUsers[idxsb]

		// remove the matched user from the store
		if err = svc.userStore.RemoveUsers(ctx, matchedUser.ID()); err != nil {
			return fmt.Errorf("failed to remove matched user: %w", err)
		}

		if err = svc.processMatch(ctx, candidate, matchedUser); err != nil {
			return fmt.Errorf("failed to process match: %w", err)
		}
	}
	return nil
}

func (svc *UserMatchProcessor) processMatch(ctx context.Context, candidate, matchedUser *User) error {
	match, createErr := svc.createAndPersistMatch(ctx, *candidate, *matchedUser)
	if createErr != nil {
		return fmt.Errorf("failed to create match: %w", createErr)
	}

	if notifyErr := svc.notifyMatch(ctx, match, candidate.ID(), matchedUser.ID()); notifyErr != nil {
		return fmt.Errorf("notification error: %w", notifyErr)
	}

	svc.logger.Info().
		Str("match_id", match.ID()).
		Strs("user_ids", []string{candidate.ID(), matchedUser.ID()}).
		Msg("match created and notified")
	return nil
}

func (svc *UserMatchProcessor) createAndPersistMatch(ctx context.Context, user1, user2 User) (*Match, error) {
	matchID := uuid.New().String()
	match, err := NewMatch(MatchID(matchID), user1, user2)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate match: %w", err)
	}

	if saveErr := svc.matchRepository.Save(ctx, match); saveErr != nil {
		return nil, fmt.Errorf("failed to save match: %w", saveErr)
	}

	return match, nil
}

func (svc *UserMatchProcessor) notifyMatch(ctx context.Context, match *Match, userIDs ...string) error {
	for _, uid := range userIDs {
		if err := svc.notifications.Notify(ctx, uid, match); err != nil {
			svc.logger.Error().
				Err(err).
				Str("user_id", uid).
				Str("match_id", match.ID()).
				Msg("failed to notify user about match")
			return err
		}
	}
	return nil
}

func (svc *UserMatchProcessor) ensureDependencies() error {
	if svc.notifications == nil {
		return domain_error.New("missing notifications channel")
	}
	if svc.matchRepository == nil {
		return domain_error.New("missing match repository")
	}
	if svc.matcher == nil {
		return domain_error.New("missing user store")
	}
	return nil
}

func (svc *UserMatchProcessor) setupLogger() {
	*svc.logger = svc.logger.With().Str("component", "randomtalk.matchmaking.user_matchmaker").Logger()
}
