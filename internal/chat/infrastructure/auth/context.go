package auth

import "context"

type contextKey string

type SessionContext struct {
	UserID string
}

func (a SessionContext) GetUserID() string {
	return a.UserID
}

func ContextWithSession(ctx context.Context, session SessionContext) context.Context {
	return context.WithValue(ctx, contextKey("session"), session)
}

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	session := SessionContext{
		UserID: userID,
	}
	return ContextWithSession(ctx, session)
}

func SessionFromContext(ctx context.Context) (SessionContext, bool) {
	session, ok := ctx.Value(contextKey("session")).(SessionContext)
	return session, ok
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	session, ok := SessionFromContext(ctx)
	if !ok {
		return "", false
	}
	return session.UserID, true
}
