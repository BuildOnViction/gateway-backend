package project

import (
	"context"
	"fmt"
	"testing"

	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	store "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	"github.com/globalsign/mgo/bson"
)

func Test_service_Create(t *testing.T) {
	mongoConnection, _ := store.NewMongo("mongodb://tomobridgejob:anhlavip@localhost:27017/tomobridgejob", "tomobridgejob")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "User", "0xd106159eC58BD2EAf5B62eF4e9cDb286170B0Bb9")

	type fields struct {
		db *store.Mongo
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Project
		wantErr bool
	}{
		{
			name: "simple_create",
			fields: fields{
				db: mongoConnection,
			},
			args: args{
				ctx:  ctx,
				name: "simple proect name",
			},
			want: entity.Project{
				Name: "simple proect name",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				db: tt.fields.db,
			}
			got, err := s.Create(tt.args.ctx, tt.args.name)
			fmt.Printf("%+v %+v", got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

// func Test_service_List(t *testing.T) {
// 	type fields struct {
// 		db *Mongo
// 	}
// 	type args struct {
// 		ctx context.Context
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    []entity.Project
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service{
// 				db: tt.fields.db,
// 			}
// 			got, err := s.List(tt.args.ctx)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("service.List() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("service.List() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_service_Delete(t *testing.T) {
// 	type fields struct {
// 		db *Mongo
// 	}
// 	type args struct {
// 		ctx context.Context
// 		id  bson.ObjectId
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := service{
// 				db: tt.fields.db,
// 			}
// 			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
// 				t.Errorf("service.Delete() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func Test_service_Update(t *testing.T) {
	mongoConnection, _ := store.NewMongo("mongodb://tomobridgejob:anhlavip@localhost:27017/tomobridgejob", "tomobridgejob")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "User", "0xd106159eC58BD2EAf5B62eF4e9cDb286170B0Bb9")
	type fields struct {
		db *Mongo
	}
	type args struct {
		ctx     context.Context
		project entity.Project
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Project
		wantErr bool
	}{
		{
			name: "only_name",
			fields: fields{
				db: mongoConnection,
			},
			args: args{
				ctx: ctx,
				project: entity.Project{
					ID:   bson.ObjectIdHex("5efc418eef88f00fc1a36c26"),
					Name: "sdfsdfsdfsdf",
				},
			},
			want: entity.Project{
				ID:   bson.ObjectIdHex("5efc418eef88f00fc1a36c26"),
				Name: "sdfsdfsdfsdf",
			},
			wantErr: false,
		},
		{
			name: "name_address",
			fields: fields{
				db: mongoConnection,
			},
			args: args{
				ctx: ctx,
				project: entity.Project{
					ID:   bson.ObjectIdHex("5efc418eef88f00fc1a36c26"),
					Name: "name_address",
					Addresses: entity.ProjectAddresses{
						MintingAddress:      "0xd106159eC58BD2EAf5B62eF4e9cDb286170B0Bb9",
						WatchSmartContracts: []string{"0xd106159eC58BD2EAf5B62eF4e9cDb286170B0Bb9"},
					},
				},
			},
			want: entity.Project{
				ID:   bson.ObjectIdHex("5efc418eef88f00fc1a36c26"),
				Name: "name_address",
				Addresses: entity.ProjectAddresses{
					MintingAddress:      "0xd106159eC58BD2EAf5B62eF4e9cDb286170B0Bb9",
					WatchSmartContracts: []string{"0xd106159eC58BD2EAf5B62eF4e9cDb286170B0Bb9"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				db: tt.fields.db,
			}

			if err := s.Update(tt.args.ctx, tt.args.project); (err != nil) != tt.wantErr {
				t.Errorf("service.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
