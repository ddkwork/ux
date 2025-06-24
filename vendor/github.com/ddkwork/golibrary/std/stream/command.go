package stream

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ddkwork/golibrary/std/mylog"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var mu sync.Mutex

func RunCommandSafe(arg ...string) (stdOut *GeneratedFile) {
	mu.Lock()
	defer mu.Unlock()
	return RunCommand(arg...)
}
func RunCommand(arg ...string) (stdOut *GeneratedFile) {
	return runCommand("", arg...)
}

func RunCommandWithDir(dir string, arg ...string) (stdOut *GeneratedFile) {
	if !IsDir(dir) {
		panic(dir + " is not exist,please check your arg ")
	}
	return runCommand(dir, arg...)
}

func runCommand(dir string, arg ...string) (stdOut *GeneratedFile) {
	if strings.Contains(arg[0], " ") && len(arg) == 1 { //命令已被合并，需要分割，取出第一个命令作为执行的path
		arg = strings.Split(arg[0], " ")
	}
	type setup struct {
		init       func()
		fastModel  func()
		slowModel  func()
		handleWait func() //merge exit code.and error
	}
	cmdKey := "command" //fast the key is clang target file path
	fast := true
	stdOut = NewGeneratedFile()
	stderr := NewGeneratedFile()

	var (
		cmd         *exec.Cmd
		stdoutPipe  io.ReadCloser
		stderrPipe  io.ReadCloser
		output      = make(chan string)
		errorOutput = make(chan string)
		done        = make(chan struct{})
	)

	const (
		execTimeout = 1 * time.Minute // 命令最大执行时间
		waitDelay   = 5 * time.Second // 结束后的最大等待时间
	)

	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	s := setup{
		init: func() {
			binaryPath := arg[0]
			switch binaryPath { //todo add more need fast model
			case "clang", "clang-format":
				cmdKey = filepath.Base(arg[len(arg)-1])
			case "ping", "go":
				fast = false
			}

			cmd = exec.CommandContext(ctx, binaryPath, arg[1:]...)
			cmd.Dir = dir //需要切换到对应目录，避免使用os.chdir,应用场景：批量更新工作区下的mod
			cmd.WaitDelay = waitDelay

			mylog.Success(cmdKey, cmd.String())
			WriteAppend("cmd.cmd", binaryPath+" "+strings.Join(arg[1:], " ")+"\n")

			stdoutPipe = mylog.Check2(cmd.StdoutPipe())
			stderrPipe = mylog.Check2(cmd.StderrPipe())
		},
		fastModel: func() {
			mylog.Check2(stdOut.ReadFrom(stdoutPipe))
			mylog.Check2(stderr.ReadFrom(stderrPipe))
		},
		slowModel: func() {
			g := sync.WaitGroup{} //这里只是处理单个文件的命令输出，不要加锁，否则会无限等待。应该在外部的文件列表遍历的地方加锁来保证同一个mod的并发更新读写模块安全

			// 启动 goroutine 读取 stdout
			g.Go(func() {
				mylog.Call(func() {
					scanner := bufio.NewScanner(stdoutPipe)
					for scanner.Scan() {
						output <- ConvertUtf82Gbk(scanner.Text())
					}
				})
			})

			// 启动 goroutine 读取 stderr
			g.Go(func() {
				mylog.Call(func() {
					scanner := bufio.NewScanner(stderrPipe)
					for scanner.Scan() {
						errorOutput <- ConvertUtf82Gbk(scanner.Text())
					}
				})
			})

			// 启动 goroutine 统一处理输出
			g.Go(func() {
				mylog.Call(func() {
					g.Wait()
					close(output)
				})
			})

			g.Go(func() {
				for line := range output {
					println(line)
					stdOut.P(line)
				}
				done <- struct{}{}
			})

			g.Go(func() {
				for line := range errorOutput {
					stderr.P(line)
				}
				done <- struct{}{}
			})

		},
		handleWait: func() {
			select { // 等待 goroutine 完成,不写在这里会导致全部携程死锁，原因不明
			case <-ctx.Done():
			case <-done:
			default:
			}
			e := cmd.Wait()
			if e != nil {
				bug := stderr.String() + "\n" + e.Error()
				mylog.Check(bug)
			}
		},
	}
	s.init()
	mylog.Check(cmd.Start())
	defer func() {
		s.handleWait()
		stdOut.TrimSuffix("\n")
	}()
	if fast {
		s.fastModel()
		return
	}
	s.slowModel()
	return
}

func ConvertUtf82Gbk(src string) string {
	if IsWindows() {
		c := mylog.Check2(simplifiedchinese.GB18030.NewDecoder().String(src)) // todo test rune
		return strings.TrimSuffix(c, "\r\n")
	}
	return src
}
