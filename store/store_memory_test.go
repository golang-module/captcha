package store

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	s := NewStoreMemory(GCLimitNumber, Expiration)
	id := "captcha id"
	d := "random-string"
	_ = s.Set(id, d)
	d2 := s.Get(id, false)
	if d2 != d {
		t.Errorf("saved %v, getDigits returned got %v", d, d2)
	}
}

func TestGetClear(t *testing.T) {
	s := NewStoreMemory(GCLimitNumber, Expiration)
	id := "captcha id"
	d := "932839jfffjkdss"
	_ = s.Set(id, d)
	d2 := s.Get(id, true)
	if d != d2 {
		t.Errorf("saved %v, getDigitsClear returned got %v", d, d2)
	}
	d2 = s.Get(id, false)
	if d2 != "" {
		t.Errorf("getDigitClear didn't clear (%q=%v)", id, d2)
	}
}

func BenchmarkSetCollect(b *testing.B) {
	b.StopTimer()
	d := "fdskfew9832232r"
	s := NewStoreMemory(9999, -1)
	ids := make([]string, 1000)
	for i := range ids {
		ids[i] = fmt.Sprintf("%d", rand.Int63())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_ = s.Set(ids[j], d)
		}
	}
}

func TestMemoryStore_SetGoCollect(t *testing.T) {
	s := NewStoreMemory(10, -1)
	for i := 0; i <= 100; i++ {
		_ = s.Set(fmt.Sprint(i), fmt.Sprint(i))
	}
}

func TestMemoryStore_CollectNotExpire(t *testing.T) {
	s := NewStoreMemory(10, time.Hour)
	for i := 0; i < 50; i++ {
		_ = s.Set(fmt.Sprint(i), fmt.Sprint(i))
	}

	// let background goroutine to go
	time.Sleep(time.Second)

	if v := s.Get("0", false); v != "0" {
		t.Error("mem store get failed")
	}
}

func TestNewMemoryStore(t *testing.T) {
	type args struct {
		collectNum int
		expiration time.Duration
	}
	tests := []struct {
		name string
		args args
		want Store
	}{
		{"", args{20, time.Hour}, nil},
		{"", args{20, time.Hour * 5}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStoreMemory(tt.args.collectNum, tt.args.expiration); got == nil {
				t.Errorf("NewMemoryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryStore_Verify(t *testing.T) {
	thisStore := NewStoreMemory(10, time.Hour)
	_ = thisStore.Set("xx", "xx")
	got := thisStore.Verify("xx", "xx", false)
	if !got {
		t.Error("failed1")
	}
	got = thisStore.Verify("xx", "xx", true)

	if !got {
		t.Error("failed2")
	}
	got = thisStore.Verify("xx", "xx", true)

	if got {
		t.Error("failed3")
	}
	got = DefaultStoreMemory.Verify("saaf", "", true)
	if got {
		t.Error("CVE-2023-45292 GO-2023-2386")
	}
}

func Test_memoryStore_Get(t *testing.T) {
	thisStore := NewStoreMemory(10, time.Hour)
	_ = thisStore.Set("xx", "xx")
	got := thisStore.Get("xx", false)
	if got != "xx" {
		t.Error("failed1")
	}
	got = thisStore.Get("xx", true)
	if got != "xx" {
		t.Error("failed2")
	}
	got = thisStore.Get("xx", false)
	if got == "xx" {
		t.Error("failed3")
	}
}
