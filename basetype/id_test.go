package basetype

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func TestId_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		id      Id
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			id:      1,
			want:    []byte(`"1"`),
			wantErr: false,
		},
		{
			name:    "2",
			id:      10000000000000000000,
			want:    []byte(`"10000000000000000000"`),
			wantErr: false,
		},
		{
			name:    "max",
			id:      math.MaxUint64,
			want:    []byte(`"18446744073709551615"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.id.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Id.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Id.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestId_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	var id Id
	tests := []struct {
		name    string
		id      *Id
		args    args
		want    Id
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			id:   &id,
			args: args{
				data: []byte(`"1"`),
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "max",
			id:   &id,
			args: args{
				data: []byte(`"18446744073709551615"`),
			},
			want:    18446744073709551615,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.id.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Id.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(id, tt.want) {
				t.Errorf("Id.MarshalJSON() = %d, want %d", id, tt.want)
			}
		})
	}
}

func TestId_MarshalAndUnmarshalJSON(t *testing.T) {
	id := Id(math.MaxUint64)
	data, err := json.Marshal(id)
	if err != nil {
		t.Fatal(err)
	}

	var nid Id
	if err := json.Unmarshal(data, &nid); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(id, nid) {
		t.Errorf("Id.MarshalAndUnmarshalJSON() = %d, want %d", id, nid)
	}
}
