package storage

import (
	"context"
	"fmt"

	"github.com/openchami/named-resource-service/internal/storage/ent"
	"github.com/openchami/named-resource-service/internal/storage/ent/label"
	entresource "github.com/openchami/named-resource-service/internal/storage/ent/resource"

	v1 "github.com/openchami/named-resource-service/apis/named-resource-service.openchami.org/v1"
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

// Querynameds returns a query builder for nameds
func Querynameds(ctx context.Context) *ent.ResourceQuery {
	return QueryResources(ctx, "Named")
}

// GetNamedByUID loads a single Named by UID
func GetNamedByUID(ctx context.Context, uid string) (*v1.Named, error) {
	ensureEntClient()
	r, err := entClient.Resource.Query().
		Where(entresource.UIDEQ(uid), entresource.KindEQ("Named")).
		WithLabels().
		WithAnnotations().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load Named %s: %w", uid, err)
	}
	v, err := FromEntResource(ctx, r)
	if err != nil {
		return nil, err
	}
	return v.(*v1.Named), nil
}

// ListnamedsByLabels returns nameds matching all provided labels
func ListnamedsByLabels(ctx context.Context, labels map[string]string) ([]*v1.Named, error) {
	q, err := QueryResourcesByLabels(ctx, "Named", labels)
	if err != nil {
		return nil, err
	}
	rs, err := q.WithLabels().WithAnnotations().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*v1.Named, 0, len(rs))
	for _, r := range rs {
		v, err := FromEntResource(ctx, r)
		if err != nil {
			continue
		}
		out = append(out, v.(*v1.Named))
	}
	return out, nil
}
