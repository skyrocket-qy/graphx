package redis

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/skyrocketOoO/zanazibar-dag/domain"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) (*RedisRepository, error) {
	return &RedisRepository{
		client: client,
	}, nil
}

func (r *RedisRepository) Ping(c context.Context) error {
	return r.client.Ping(c).Err()
}

func (r *RedisRepository) Get(c context.Context, edge domain.Edge,
	queryMode bool) ([]domain.Edge, error) {
	from, to := edgeToKeyValue(edge)
	if queryMode {
		if edge == (domain.Edge{}) {
			result := r.client.Keys(c, "!reverse%*")
			if err := result.Err(); err != nil {
				return nil, err
			}
			keys := []string{}
			if err := result.ScanSlice(&keys); err != nil {
				return nil, err
			}
			edges := []domain.Edge{}
			for _, key := range keys {
				res := r.client.SMembers(c, key)
				if err := res.Err(); err != nil {
					return nil, err
				}
				var tos []string
				if err := res.ScanSlice(&tos); err != nil {
					return nil, err
				}
				fromSplit := strings.Split(key, "%")
				for _, to := range tos {
					toSplit := strings.Split(to, "%")
					edges = append(edges, domain.Edge{
						ObjNs:   toSplit[0],
						ObjName: toSplit[1],
						ObjRel:  toSplit[2],
						SbjNs:   fromSplit[0],
						SbjName: fromSplit[1],
						SbjRel:  fromSplit[2],
					})
				}
			}
			return edges, nil
		} else {
			if to != "%%" {
				to = "reverse%" + to
				memberStrings := []string{}
				stringSliceCmd := r.client.SMembers(c, to)
				if stringSliceCmd.Err() != nil {
					return nil, stringSliceCmd.Err()
				}
				if err := stringSliceCmd.ScanSlice(&memberStrings); err != nil {
					return nil, err
				}
				edges := []domain.Edge{}
				for _, member := range memberStrings {
					strSplit := strings.Split(member, "%")
					edges = append(edges, domain.Edge{
						SbjNs:   strSplit[0],
						SbjName: strSplit[1],
						SbjRel:  strSplit[2],
						ObjNs:   edge.SbjNs,
						ObjName: edge.SbjName,
						ObjRel:  edge.SbjRel,
					})
				}
				return edges, nil
			} else {
				memberStrings := []string{}
				stringSliceCmd := r.client.SMembers(c, from)
				if stringSliceCmd.Err() != nil {
					return nil, stringSliceCmd.Err()
				}
				if err := stringSliceCmd.ScanSlice(&memberStrings); err != nil {
					return nil, err
				}
				edges := []domain.Edge{}
				for _, member := range memberStrings {
					strSplit := strings.Split(member, "%")
					edges = append(edges, domain.Edge{
						ObjNs:   strSplit[0],
						ObjName: strSplit[1],
						ObjRel:  strSplit[2],
						SbjNs:   edge.SbjNs,
						SbjName: edge.SbjName,
						SbjRel:  edge.SbjRel,
					})
				}
				return edges, nil
			}
		}
	} else {
		rdsBoolCmd := r.client.SIsMember(c, from, to)
		if rdsBoolCmd.Err() != nil {
			return nil, rdsBoolCmd.Err()
		}
		if rdsBoolCmd.Val() {
			return []domain.Edge{edge}, nil
		} else {
			return nil, domain.ErrRecordNotFound{}
		}
	}
}

func (r *RedisRepository) Create(c context.Context, edge domain.Edge) error {
	from, to := edgeToKeyValue(edge)
	if err := r.client.SAdd(c, from, to).Err(); err != nil {
		return err
	}
	return r.client.SAdd(c, "reverse%"+to, from).Err()
}

func (r *RedisRepository) Delete(c context.Context, edge domain.Edge,
	queryMode bool) error {
	from, to := edgeToKeyValue(edge)
	if queryMode {
		result := r.client.Keys(c, from)
		if err := result.Err(); err != nil {
			return nil, err
		}
		keys := []string{}
		if err := result.ScanSlice(&keys); err != nil {
			return nil, err
		}
	} else {
		return r.client.SRem(c, from, to).Err()
	}
}

func (r *RedisRepository) ClearAll(c context.Context) error {
	return r.client.FlushDB(c).Err()
}

func vertexToString(v domain.Vertex) string {
	return v.Ns + "%" + v.Name + "%" + v.Rel
}

func edgeToKeyValue(edge domain.Edge) (from string, to string) {
	from = vertexToString(domain.Vertex{
		Ns:   edge.SbjNs,
		Name: edge.SbjName,
		Rel:  edge.SbjRel,
	})
	to = vertexToString(domain.Vertex{
		Ns:   edge.ObjNs,
		Name: edge.ObjName,
		Rel:  edge.ObjRel,
	})
	return
}

func vertexToPattern(vertex domain.Vertex) string {
	if vertex == (domain.Vertex{}) {
		return "*"
	}
	res := ""
	if vertex.Ns == "" {
		res += "*"
	} else {
		res += vertex.Ns
	}
	res += "%"
	if vertex.Name == "" {
		res += "*"
	} else {
		res += vertex.Name
	}
	res += "%"

}
