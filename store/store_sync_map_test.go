package store

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

var duration = time.Second * 2
var mapStore = NewStoreSyncMap(duration)

func TestNewStoreSyncMap(t *testing.T) {
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want *storeSyncMap
	}{
		{"new", args{duration}, mapStore},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStoreSyncMap(tt.args.duration); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreSyncMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreSyncMap_Get(t *testing.T) {
	mapStore.Set("1", "1")
	mapStore.Set("2", "2")

	type fields struct {
		d time.Duration
		m *sync.Map
	}
	type args struct {
		id    string
		clear bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"get", fields{duration, mapStore.m}, args{"1", false}, "1"},
		{"get", fields{duration, mapStore.m}, args{"2", true}, "2"},
		{"get", fields{duration, mapStore.m}, args{"2", true}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeSyncMap{
				d: tt.fields.d,
				m: tt.fields.m,
			}
			if got := s.Get(tt.args.id, tt.args.clear); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreSyncMap_Expire(t *testing.T) {
	mapStore.Set("2", "22")
	if v := mapStore.Get("2", false); v != "22" {
		t.Error("failed")
	}
	time.Sleep(time.Second * 2)
	if v := mapStore.Get("2", false); v != "" {
		t.Error("expire failed")
	}
}

func TestStoreSyncMap_Set(t *testing.T) {
	type fields struct {
		d time.Duration
		m *sync.Map
	}
	type args struct {
		id    string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"get", fields{duration, mapStore.m}, args{"1", "1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeSyncMap{
				d: tt.fields.d,
				m: tt.fields.m,
			}
			s.Set(tt.args.id, tt.args.value)
		})
	}
}

func TestStoreSyncMap_Verify(t *testing.T) {
	mapStore.Set("1", "1")
	mapStore.Set("2", "2")
	type fields struct {
		d time.Duration
		m *sync.Map
	}
	type args struct {
		id     string
		answer string
		clear  bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"get", fields{duration, mapStore.m}, args{"1", "1", true}, true},
		{"get", fields{duration, mapStore.m}, args{"1", "1", false}, false},
		{"get", fields{duration, mapStore.m}, args{"2", "2", true}, true},
		{"get", fields{duration, mapStore.m}, args{"2", "2", false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeSyncMap{
				d: tt.fields.d,
				m: tt.fields.m,
			}
			if got := s.Verify(tt.args.id, tt.args.answer, tt.args.clear); got != tt.want {
				t.Errorf("Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreSyncMap_delete(t *testing.T) {
	type fields struct {
		d time.Duration
		m *sync.Map
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"get", fields{duration, new(sync.Map)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeSyncMap{
				d: tt.fields.d,
				m: tt.fields.m,
			}
			s.delete()
		})
	}
}
