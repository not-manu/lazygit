package helpers

import (
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
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
	pending := models.TagChecks{State: "PENDING", StartedAt: fresh}
	passed := models.TagChecks{State: "SUCCESS", StartedAt: old, CompletedAt: old.Add(time.Minute)}
	failed := models.TagChecks{State: "FAILURE", StartedAt: fresh, CompletedAt: now}
	cached := map[string]models.TagChecks{"cached": passed, "was-pending": pending}
	fetched := map[string]models.TagChecks{
		"pending":     pending,
		"was-pending": passed,
		"failed":      failed,
		"passed":      passed,
		"cached":      failed,
	}

	checks, toCache := settleTagChecks(recentTags, cached, fetched, now)

	assert.Equal(t, map[string]models.TagChecks{
		"pending":              pending,
		"was-pending":          passed,
		"failed":               failed,
		"passed":               passed,
		"fresh-without-checks": {State: "PENDING"},
		"old-without-checks":   {State: noTagChecks},
		"cached":               passed,
	}, checks)
	assert.Equal(t, map[string]models.TagChecks{
		"pending":            pending,
		"was-pending":        passed,
		"failed":             failed,
		"passed":             passed,
		"old-without-checks": {State: noTagChecks},
		"cached":             passed,
	}, toCache)
	assert.Equal(t, map[string]models.TagChecks{"cached": passed, "was-pending": pending}, cached)
}

func TestSettleTagChecksIsNotOptimisticWithoutTagChecks(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	recentTags := []git_commands.RecentTag{
		{Name: "fresh", CreatedAt: now.Add(-time.Minute)},
		{Name: "old", CreatedAt: now.Add(-time.Hour)},
	}

	checks, _ := settleTagChecks(recentTags, map[string]models.TagChecks{}, map[string]models.TagChecks{}, now)

	assert.Equal(t, map[string]models.TagChecks{"old": {State: noTagChecks}}, checks)
}

func TestHasRunningTagChecks(t *testing.T) {
	started := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	assert.True(t, hasRunningTagChecks(map[string]models.TagChecks{"a": {State: "PENDING", StartedAt: started}}))
	assert.False(t, hasRunningTagChecks(map[string]models.TagChecks{"a": {State: "PENDING"}}))
	assert.False(t, hasRunningTagChecks(map[string]models.TagChecks{"a": {State: "SUCCESS", StartedAt: started, CompletedAt: started}}))
}
