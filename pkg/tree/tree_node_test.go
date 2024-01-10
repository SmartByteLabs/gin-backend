package tree

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestCreateFieldTree(t *testing.T) {
	type args struct {
		ctx      context.Context
		mockFunc func(sqlmock.Sqlmock)
		configID int
	}
	tests := []struct {
		name    string
		args    args
		want    *Fields
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx:      context.Background(),
				configID: 1,
				mockFunc: func(mock sqlmock.Sqlmock) {
					// this should set id", "parent_id", "name", "type", "required fields and value should be 1, null, "rootField", "object", false
					mock.ExpectQuery("SELECT ").
						WithArgs(1).
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "parent_id", "name", "type", "required"}).
								AddRow(1, nil, "rootField", Object, false),
						)
				},
			},
			want:    &Fields{ID: 1, Name: "rootField", Type: Object, Required: false},
			wantErr: false,
		},
		{
			name: "3 level test",
			args: args{
				ctx:      context.Background(),
				configID: 1,
				mockFunc: func(mock sqlmock.Sqlmock) {
					// this should set id", "parent_id", "name", "type", "required fields and value should be 1, null, "rootField", Object, false
					mock.ExpectQuery("SELECT ").
						WithArgs(1).
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "parent_id", "name", "type", "required"}).
								AddRow(1, nil, "root", Object, false).
								AddRow(2, 1, "child1", Object, false).
								AddRow(3, 1, "child2", Object, false).
								AddRow(4, 2, "child1_1", Object, false).
								AddRow(5, 2, "child1_2", Object, false).
								AddRow(6, 3, "child2_1", Object, false).
								AddRow(7, 3, "child2_2", Object, false),
						)
				},
			},
			want: &Fields{
				ID:       1,
				Name:     "root",
				Type:     Object,
				Required: false,
				ChildField: []*Fields{
					{
						ID:       2,
						Name:     "child1",
						Type:     Object,
						Required: false,
						ChildField: []*Fields{
							{
								ID:       4,
								Name:     "child1_1",
								Type:     Object,
								Required: false,
							},
							{
								ID:       5,
								Name:     "child1_2",
								Type:     Object,
								Required: false,
							},
						},
					},
					{
						ID:       3,
						Name:     "child2",
						Type:     Object,
						Required: false,
						ChildField: []*Fields{
							{
								ID:       6,
								Name:     "child2_1",
								Type:     Object,
								Required: false,
							},
							{
								ID:       7,
								Name:     "child2_2",
								Type:     Object,
								Required: false,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock database
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.args.mockFunc(mock)

			got, err := CreateFieldTree(tt.args.ctx, db, tt.args.configID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFieldTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Assert that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			// compare got with want with deep equal
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFieldTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateNodeTree(t *testing.T) {
	type args struct {
		ctx       context.Context
		mockFunc  func(sqlmock.Sqlmock)
		configID  int64
		templeID  int64
		tableName string
	}
	tests := []struct {
		name    string
		args    args
		want    *Node
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx:       context.Background(),
				configID:  1,
				templeID:  1,
				tableName: "tree_values",
				mockFunc: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT ").
						WithArgs(1, 1).
						WillReturnRows(
							sqlmock.NewRows([]string{"field_id", "field_parent_id", "field_name", "field_type", "field_required", "value_parent_id", "value_id", "value"}).
								AddRow(1, nil, "rootField", Object, false, nil, 1, ""),
						)
				},
			},
			want: &Node{
				Fields: Fields{
					ID:       1,
					Name:     "rootField",
					Type:     Object,
					Required: false,
				},
				Value: "",
			},
			wantErr: false,
		},
		{
			name: "3 level test",
			args: args{
				ctx:       context.Background(),
				configID:  1,
				templeID:  1,
				tableName: "tree_values",
				mockFunc: func(mock sqlmock.Sqlmock) {
					mock.ExpectQuery("SELECT ").
						WithArgs(1, 1).
						WillReturnRows(
							sqlmock.NewRows([]string{"field_id", "field_parent_id", "field_name", "field_type", "field_required", "value_parent_id", "value_id", "value"}).
								AddRow(1, nil, "root", Object, false, nil, 1, "").
								AddRow(2, 1, "child1", Object, false, 1, 2, "").
								AddRow(3, 1, "child2", Object, false, 1, 3, "").
								AddRow(4, 2, "child1_1", Object, false, 2, 4, "").
								AddRow(5, 2, "child1_2", Object, false, 2, 5, "").
								AddRow(6, 3, "child2_1", Object, false, 3, 6, "").
								AddRow(7, 3, "child2_2", Object, false, 3, 7, ""),
						)
				},
			},
			want: &Node{
				Fields: Fields{
					ID:       1,
					Name:     "root",
					Type:     Object,
					Required: false,
				},
				Value: "",
				Children: []*Node{
					{
						Fields: Fields{
							ID:       2,
							Name:     "child1",
							Type:     Object,
							Required: false,
						},
						Value: "",
						Children: []*Node{
							{
								Fields: Fields{
									ID:       4,
									Name:     "child1_1",
									Type:     Object,
									Required: false,
								},
								Value: "",
							},
							{
								Fields: Fields{
									ID:       5,
									Name:     "child1_2",
									Type:     Object,
									Required: false,
								},
								Value: "",
							},
						},
					},
					{
						Fields: Fields{
							ID:       3,
							Name:     "child2",
							Type:     Object,
							Required: false,
						},
						Value: "",
						Children: []*Node{
							{
								Fields: Fields{
									ID:       6,
									Name:     "child2_1",
									Type:     Object,
									Required: false,
								},
								Value: "",
							},
							{
								Fields: Fields{
									ID:       7,
									Name:     "child2_2",
									Type:     Object,
									Required: false,
								},
								Value: "",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.args.mockFunc(mock)

			got, err := CreateNodeTree(tt.args.ctx, db, tt.args.configID, tt.args.templeID, tt.args.tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNodeTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Assert that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				g, _ := json.Marshal(got)
				w, _ := json.Marshal(tt.want)
				t.Errorf("CreateNodeTree() = %v, want %v", string(g), string(w))
			}
		})
	}
}
