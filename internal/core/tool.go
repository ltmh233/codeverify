package core

import (
	"codeverify/internal/analysis"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ToolName string

const (
	ToolReadFile      ToolName = "read_file"
	ToolWriteFile     ToolName = "write_file"
	ToolRunTest       ToolName = "run_test"
	ToolReadDir       ToolName = "read_dir"
	ToolRecordRisk    ToolName = "record_risk"
	ToolSubmitVerdict ToolName = "submit_verdict"
	ToolGetCall       ToolName = "get_call"
)

type ToolCall struct {
	Name    ToolName
	RawArgs string
	CallID  string
}

type ToolWriteFileArguments struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type ToolReadDirArguments struct {
	Path string `json:"path"`
}

type ToolReadFileArguments struct {
	Path string `json:"path"`
}

type RunContext struct {
	RepoPath string
	RepoAbs  string
	State    *State
	Analysis *analysis.Analysis
}

// 判断给定路径是否超出目标仓库路径
func ResolveRepoPath(repoAbs string, target string) (string, error) {
	targetAbs, err := filepath.Abs(
		filepath.Join(
			repoAbs,
			filepath.Clean(target),
		),
	)
	if err != nil {
		return "", err
	}
	//filepath.Rel返回targPath相对于basePath的位置
	rel, err := filepath.Rel(repoAbs, targetAbs)
	if err != nil {
		return "", err
	}
	//如果rel是父目录或者是父目录的子目录或自文件, 则报错
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes repository")
	}

	return targetAbs, nil
}

func SanitizePath(targetAbs string, repoAbs string) string {

	return strings.ReplaceAll(targetAbs, repoAbs, ".")
}
