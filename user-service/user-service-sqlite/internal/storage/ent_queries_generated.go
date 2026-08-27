package storage

import (
	"context"
	"fmt"

	"github.com/openchami/user-service-pr-sqlite/internal/storage/ent"
	"github.com/openchami/user-service-pr-sqlite/internal/storage/ent/label"
	entresource "github.com/openchami/user-service-pr-sqlite/internal/storage/ent/resource"

	v1 "github.com/openchami/user-service-pr-sqlite/apis/storage-annotations.openchami.org/v1"
)

// ensureEntClient verifies the ent client has been initialized
func ensureEntClient() {
	if entClient == nil {
		panic("ent client not initialized: call SetEntClient in main.go before using storage")
	}
}

// QueryResources returns a generic query builder for a given kind
func QueryResources(ctx context.Context, kind string) *ent.ResourceQuery {
	ensureEntClient()
	return entClient.Resource.Query().
		Where(entresource.KindEQ(kind))
}

// QueryResourcesByLabels queries resources by kind and exact-match labels
func QueryResourcesByLabels(ctx context.Context, kind string, labels map[string]string) (*ent.ResourceQuery, error) {
	ensureEntClient()
	q := entClient.Resource.Query().Where(entresource.KindEQ(kind))
	for k, v := range labels {
		q = q.Where(entresource.HasLabelsWith(
			label.KeyEQ(k),
			label.ValueEQ(v),
		))
	}
	return q, nil
}

// Queryusers returns a query builder for users
func Queryusers(ctx context.Context) *ent.ResourceQuery {
	return QueryResources(ctx, "User")
}

// GetUserByUID loads a single User by UID
func GetUserByUID(ctx context.Context, uid string) (*v1.User, error) {
	ensureEntClient()
	r, err := entClient.Resource.Query().
		Where(entresource.UIDEQ(uid), entresource.KindEQ("User")).
		WithLabels().
		WithAnnotations().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load User %s: %w", uid, err)
	}
	v, err := FromEntResource(ctx, r)
	if err != nil {
		return nil, err
	}
	return v.(*v1.User), nil
}

// ListusersByLabels returns users matching all provided labels
func ListusersByLabels(ctx context.Context, labels map[string]string) ([]*v1.User, error) {
	q, err := QueryResourcesByLabels(ctx, "User", labels)
	if err != nil {
		return nil, err
	}
	rs, err := q.WithLabels().WithAnnotations().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.User, 0, len(rs))
	for _, r := range rs {
		v, err := FromEntResource(ctx, r)
		if err != nil {
			continue
		}
		out = append(out, v.(*v1.User))
	}
	return out, nil
}
