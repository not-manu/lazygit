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
		{Name: "failed", CreatedAt: fresh},
		{Name: "passed", CreatedAt: old},
		{Name: "fresh-without-checks", CreatedAt: fresh},
		{Name: "old-without-checks", CreatedAt: old},
		{Name: "cached", CreatedAt: old},
	}
	cached := map[string]string{"cached": "SUCCESS"}
	fetched := map[string]string{
		"pending": "PENDING",
		"failed":  "FAILURE",
		"passed":  "SUCCESS",
		"cached":  "FAILURE",
	}

	states, toCache := settleTagChecks(recentTags, cached, fetched, now)

	assert.Equal(t, map[string]string{
		"pending": "PENDING",
		"failed":  "FAILURE",
		"passed":  "SUCCESS",
		"cached":  "SUCCESS",
	}, states)
	assert.Equal(t, map[string]string{
		"failed":             "FAILURE",
		"passed":             "SUCCESS",
		"old-without-checks": noTagChecks,
		"cached":             "SUCCESS",
	}, toCache)
	assert.Equal(t, map[string]string{"cached": "SUCCESS"}, cached)
}
