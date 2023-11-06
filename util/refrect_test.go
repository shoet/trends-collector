package util

import (
	"testing"

	"github.com/shoet/trends-collector/testutil"
)

func Test_MergeStruct(t *testing.T) {
	src := &struct {
		FieldA string
		FieldB int
		FieldC string
	}{
		FieldA: "a",
		FieldB: 1,
		FieldC: "c",
	}

	opt := &struct {
		FieldA string
		FieldB int
	}{
		FieldA: "A",
		FieldB: 2,
	}

	MergeStruct(src, opt)

	want := &struct {
		FieldA string
		FieldB int
		FieldC string
	}{
		FieldA: "A",
		FieldB: 2,
		FieldC: "c",
	}

	testutil.AssertObject(t, src, want)
}
