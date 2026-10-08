package config

import "github.com/jesseduffield/lazygit/pkg/commands/models"

const (
	githubPullRequestsCacheFileName = "github_pull_requests.json"
	githubTagChecksCacheFileName    = "github_tag_checks.json"
)

// CachedPullRequest stores the essential fields of a GitHub pull request.
type CachedPullRequest struct {
	HeadRefName         string `json:"headRefName"`
	Number              int    `json:"number"`
	Title               string `json:"title"`
	State               string `json:"state"`
	ChecksState         string `json:"checksState,omitempty"`
	Url                 string `json:"url"`
	HeadRepositoryOwner string `json:"headRepositoryOwner"`
}

type githubPullRequestCache = repoCache[[]CachedPullRequest]

func loadGithubPullRequestCache() *githubPullRequestCache {
	return loadRepoCache[[]CachedPullRequest](githubPullRequestsCacheFileName, "GitHub pull request")
}

func githubPullRequestCachePath() (string, error) {
	return repoCachePath(githubPullRequestsCacheFileName)
}

func newGithubPullRequestCache(path string) *githubPullRequestCache {
	return newRepoCache[[]CachedPullRequest](path, "GitHub pull request")
}

type githubTagChecksCache = repoCache[map[string]models.TagChecks]

func loadGithubTagChecksCache() *githubTagChecksCache {
	return loadRepoCache[map[string]models.TagChecks](githubTagChecksCacheFileName, "GitHub tag checks")
}

func newGithubTagChecksCache(path string) *githubTagChecksCache {
	return newRepoCache[map[string]models.TagChecks](path, "GitHub tag checks")
}
