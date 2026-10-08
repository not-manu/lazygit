package helpers

import (
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/stretchr/testify/assert"
)

func TestSettleTagChecks(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	old := now.Add(-time.Hour)

	recentTags := []git_commands.RecentTag{
		{Name: "pending", CreatedAt: fresh},
		{Name: "was-pending", CreatedAt: fresh},
		{Name: "failed", CreatedAt: fresh},
		{Name: "passed", CreatedAt: old},
		{Name: "fresh-without-checks", CreatedAt: fresh},
		{Name: "old-without-checks", CreatedAt: old},
		{Name: "cached", CreatedAt: old},
	}
	cached := map[string]string{"cached": "SUCCESS", "was-pending": "PENDING"}
	fetched := map[string]string{
		"pending":     "PENDING",
		"was-pending": "SUCCESS",
		"failed":      "FAILURE",
		"passed":      "SUCCESS",
		"cached":      "FAILURE",
	}

	states, toCache := settleTagChecks(recentTags, cached, fetched, now)

	assert.Equal(t, map[string]string{
		"pending":              "PENDING",
		"was-pending":          "SUCCESS",
		"failed":               "FAILURE",
		"passed":               "SUCCESS",
		"fresh-without-checks": "PENDING",
		"old-without-checks":   noTagChecks,
		"cached":               "SUCCESS",
	}, states)
	assert.Equal(t, map[string]string{
		"pending":            "PENDING",
		"was-pending":        "SUCCESS",
		"failed":             "FAILURE",
		"passed":             "SUCCESS",
		"old-without-checks": noTagChecks,
		"cached":             "SUCCESS",
	}, toCache)
	assert.Equal(t, map[string]string{"cached": "SUCCESS", "was-pending": "PENDING"}, cached)
}

func TestSettleTagChecksIsNotOptimisticWithoutTagChecks(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	recentTags := []git_commands.RecentTag{
		{Name: "fresh", CreatedAt: now.Add(-time.Minute)},
		{Name: "old", CreatedAt: now.Add(-time.Hour)},
	}

	states, _ := settleTagChecks(recentTags, map[string]string{}, map[string]string{}, now)

	assert.Equal(t, map[string]string{"old": noTagChecks}, states)
}
