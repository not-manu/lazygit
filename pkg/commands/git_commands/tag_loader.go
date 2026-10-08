package git_commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type TagLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewTagLoader(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *TagLoader {
	return &TagLoader{
		Common: common,
		cmd:    cmd,
	}
}

func (self *TagLoader) GetTags() ([]*models.Tag, error) {
	// get remote branches, sorted  by creation date (descending)
	// see: https://git-scm.com/docs/git-tag#Documentation/git-tag.txt---sortltkeygt
	cmdArgs := NewGitCmd("tag").Arg("--list", "-n", "--sort=-creatordate").ToArgv()
	tagsOutput, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	split := utils.SplitLines(tagsOutput)

	lineRegex := regexp.MustCompile(`^([^\s]+)(\s+)?(.*)$`)

	tags := lo.Map(split, func(line string, _ int) *models.Tag {
		matches := lineRegex.FindStringSubmatch(line)
		tagName := matches[1]
		message := ""
		if len(matches) > 3 {
			message = matches[3]
		}

		return &models.Tag{
			Name:    tagName,
			Message: message,
		}
	})

	return tags, nil
}

type RecentTag struct {
	Name      string
	CreatedAt time.Time
}

func (self *TagLoader) GetRecentTags(count int) ([]RecentTag, error) {
	cmdArgs := NewGitCmd("for-each-ref").
		Arg("--sort=-creatordate", fmt.Sprintf("--count=%d", count), "--format=%(creatordate:unix) %(refname:strip=2)", "refs/tags").
		ToArgv()
	output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	return lo.FilterMap(utils.SplitLines(output), func(line string, _ int) (RecentTag, bool) {
		timestamp, name, found := strings.Cut(line, " ")
		unix, err := strconv.ParseInt(timestamp, 10, 64)
		if !found || err != nil {
			return RecentTag{}, false
		}
		return RecentTag{Name: name, CreatedAt: time.Unix(unix, 0)}, true
	}), nil
}
