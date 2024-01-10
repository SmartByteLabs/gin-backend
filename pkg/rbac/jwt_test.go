package rbac

import (
	"fmt"
	"testing"
	"time"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
)

func TestJWTAuthValidate(t *testing.T) {
	type user struct {
		database.TableID[int64]
		Name string `json:"name"`
	}

	type args struct {
		secret string
	}
	tests := []struct {
		name string
		args args
		want *user
		err  error
	}{
		{
			name: "simple test",
			args: args{
				secret: "pradip",
			},
			want: &user{Name: "John Doe"},
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.want, 10*time.Minute, tt.args.secret)
			if err != nil {
				t.Errorf("GenerateJWT() error = %v", err)
			}

			got, err := JWTAuthValidate[user](fmt.Sprintf("Bearer %s", token), tt.args.secret)
			if err != nil {
				t.Errorf("JWTAuthValidate() error = %v", err)
				return
			}

			if got.Name != tt.want.Name {
				t.Errorf("JWTAuthValidate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
