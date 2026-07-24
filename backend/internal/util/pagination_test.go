package util_test

import (
	"testing"

	"github.com/akrom/finance-backend/internal/util"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name       string
		page, limit, defLim, maxLim int
		wantPage, wantLimit, wantOffset int
	}{
		{"normal", 1, 20, 10, 100, 1, 20, 0},
		{"second page", 2, 20, 10, 100, 2, 20, 20},
		{"zero page", 0, 20, 10, 100, 1, 20, 0},
		{"negative limit", 1, -5, 15, 100, 1, 15, 0},
		{"exceed max limit", 1, 500, 10, 100, 1, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, l, o := util.Paginate(tt.page, tt.limit, tt.defLim, tt.maxLim)
			if p != tt.wantPage || l != tt.wantLimit || o != tt.wantOffset {
				t.Errorf("Paginate(%d, %d, %d, %d) = (%d, %d, %d), want (%d, %d, %d)",
					tt.page, tt.limit, tt.defLim, tt.maxLim, p, l, o, tt.wantPage, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}

func TestBuildPageResult(t *testing.T) {
	res := util.BuildPageResult(1, 10, 25)
	if res.TotalPages != 3 {
		t.Errorf("TotalPages = %d, want 3", res.TotalPages)
	}

	res2 := util.BuildPageResult(1, 10, 30)
	if res2.TotalPages != 3 {
		t.Errorf("TotalPages exact = %d, want 3", res2.TotalPages)
	}
}

func TestHasNextPrevPage(t *testing.T) {
	if !util.HasNextPage(1, 10, 25) {
		t.Error("HasNextPage should be true for p1 l10 total25")
	}
	if util.HasNextPage(3, 10, 25) {
		t.Error("HasNextPage should be false for p3 l10 total25")
	}
	if util.HasPrevPage(1) {
		t.Error("HasPrevPage should be false for page 1")
	}
	if !util.HasPrevPage(2) {
		t.Error("HasPrevPage should be true for page 2")
	}
}
