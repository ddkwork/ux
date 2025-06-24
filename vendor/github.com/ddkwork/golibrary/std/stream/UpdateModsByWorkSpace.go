package stream

import (
	"archive/zip"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/safemap"
	"os/user"
	"sync"

	"golang.org/x/mod/modfile"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func IsRunningOnGitHubActions() bool {
	// 1. 首先检查标准环境变量（非sudo场景）
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return true
	}

	// 2. 检查当前用户和工作目录特征（处理sudo场景）
	currentUser := mylog.Check2(user.Current())
	// GitHub Actions的运行器用户
	if strings.HasPrefix(currentUser.Username, "runner") {
		return true
	}

	// GitHub Actions的标准工作目录
	workdir := mylog.Check2(os.Getwd())
	if strings.Contains(workdir, "/home/runner/work") ||
		strings.Contains(workdir, "/__w/") {
		return true
	}

	// 3. 检查运行器工具缓存目录（备用检测）
	if IsDir("/opt/hostedtoolcache") {
		return true
	}
	if IsDir("/Users/runner/hostedtoolcache") {
		return true
	}

	return false
}

func UpdateAllLocalRep() {
	if IsRunningOnGitHubActions() {
		return
	}
	reps := []string{
		"D:\\workspace\\workspace\\golibrary",
		"D:\\workspace\\workspace\\ux",
	}
	w := sync.WaitGroup{}
	for _, rep := range reps {
		w.Go(func() {
			RunCommandSafe("go get -x github.com/ddkwork/" + filepath.Base(rep) + "@" + GetLastCommitHashLocal(rep))
		})
	}
	w.Wait()
}

func GetLastCommitHashLocal(repositoryDir string) string {
	return RunCommandWithDir(repositoryDir, "git rev-parse HEAD").String()
}

func UpdateWorkSpace(isUpdateAll bool) {
	mylog.Call(func() { updateWorkSpace(isUpdateAll) })
}
func updateWorkSpace(isUpdateAll bool) {
	if !FileExists("go.work") {
		updateMod(mylog.Check2(os.Getwd()))
		return
	}
	RunCommand("go work use -r .")
	mods := make([]string, 0)
	for line := range ReadFileToLines("go.work") {
		for _, skip := range skips {
			if line == skip {
				mylog.Warning("skip", skip)
				continue
			}
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ".") {
			abs := mylog.Check2(filepath.Abs(line))
			if filepath.Base(abs) == "golibrary" {
				continue
			}
			mods = append(mods, abs)
		}
	}

	modChan := make(chan string, len(mods))

	g := sync.WaitGroup{}
	for _, modPath := range mods {
		g.Go(func() { // 每个模块单独跑,这里不能加锁，否则很慢，谨慎使用读写锁
			updateMod(modPath) // 锁应该在这里面
			if isUpdateAll {
				RunCommandSafe("go get -u -x all") // need lock,但是不使用这个，太慢了
			}
			modChan <- modPath
		})
	}
	g.Go(func() {
		for mod := range modChan {
			mylog.Success("updated mod", strconv.Quote(mod))
		}
		close(modChan)
	})
	g.Wait()
	mylog.Success("all work finished")
}
func updateMod(dir string) { // 实现替换，不要网络访问了，太慢了
	if filepath.Base(dir) == "golibrary" {
		return
	}
	r := mylog.Check2(zip.OpenReader("D:\\workspace\\workspace\\mod.zip"))
	newBody := mylog.Check2(io.ReadAll(mylog.Check2(r.Open("go.mod"))))
	originMod := filepath.Join(dir, "go.mod")
	body := mylog.Check2(os.ReadFile(originMod))
	f := mylog.Check2(modfile.Parse(originMod, body, nil))
	newModMap := ParseGoMod("new.mod", newBody)
	for oldName, oldVersion := range ParseGoMod(originMod, body).Range() {
		newVersion, exist := newModMap.Get(oldName)
		if exist {
			if oldVersion != newVersion {
				for i, require := range f.Require {
					if require.Mod.Path == oldName {
						require.Mod.Version = newVersion
						// f.Require[i] = require
						setVersion(f.Require[i], newVersion)
					}
				}
			}
		}
	}
	f.Cleanup()
	f.SortBlocks()
	updateModFile := mylog.Check2(f.Format())
	// println(string(updateModFile))
	WriteTruncate(originMod, updateModFile)
	g := sync.WaitGroup{}
	g.Go(func() {
		RunCommandWithDir(dir, "go mod tidy")
		v := newModMap.GetMust("github.com/ddkwork/golibrary")
		b := NewBuffer(originMod)
		if !b.Contains("golibrary") {
			line := "require github.com/ddkwork/golibrary  " + v
			mylog.Info("add golibrary", line)
			b.WriteStringLn(line).ReWriteSelf()
		}
		// https://github.com/ddkwork/tools/blob/master/gopls/doc/analyzers.md
		RunCommandWithDir(dir, "go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -diff ./...")
		RunCommandWithDir(dir, "go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix ./...")
	})
	g.Wait()
}

func ParseGoMod(file string, data []byte) *safemap.M[string, string] {
	f := mylog.Check2(modfile.Parse(file, data, nil))
	return safemap.NewOrdered[string, string](func(yield func(string, string) bool) {
		for _, req := range f.Require {
			yield(req.Mod.Path, req.Mod.Version)
		}
	})
}

func setVersion(r *modfile.Require, v string) {
	r.Mod.Version = v
	if line := r.Syntax; len(line.Token) > 0 {
		if line.InBlock {
			// If the line is preceded by an empty line, remove it; see
			// https://golang.org/issue/33779.
			if len(line.Comments.Before) == 1 && len(line.Comments.Before[0].Token) == 0 {
				line.Comments.Before = line.Comments.Before[:0]
			}
			if len(line.Token) >= 2 { // example.com v1.2.3
				line.Token[1] = v
			}
		} else {
			if len(line.Token) >= 3 { // require example.com v1.2.3
				line.Token[2] = v
			}
		}
	}
}

var skips = []string{
	"module gioui.org",
	"module gioui.org/cmd",
	"module gioui.org/example",
	"module gioui.org/x",
}

func GetDesktopDir() string {
	// 获取用户主目录
	homeDir := mylog.Check2(os.UserHomeDir())
	// 根据操作系统处理路径
	switch runtime.GOOS {
	case "windows", "darwin":
		// Windows和macOS直接拼接Desktop
		return filepath.Join(homeDir, "Desktop")
	case "linux":
		// Linux优先检查XDG环境变量
		if xdgDir := os.Getenv("XDG_DESKTOP_DIR"); xdgDir != "" {
			return xdgDir
		}
		// 默认使用主目录下的Desktop
		return filepath.Join(homeDir, "Desktop")
	default:
		panic("unsupported platform")
	}
}
