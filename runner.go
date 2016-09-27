package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Job in execution states
const (
	JobWorking   = "working"
	JobSubmitted = "submitted"
	JobNotFound  = "not found"
	JobFailed    = "failed"
	JobTimeout   = "timeout"
	JobCompleted = "completed"
)

// Job structure to track ran script
type Job struct {
	Name        string
	Params      []string
	UUID        string
	Created     time.Time
	Status      string
	ErrorString string
	Pid         int
	SystemTime  time.Duration
	UserTime    time.Duration
	Timeout     int
	outChan     chan string
	stateChan   chan string
	ip          string
}

//JobInterface ...
type JobInterface interface {
	Run(name, UUID string, args []string) int
	State() *chan string
	Log() *chan string
}

var logContextRunner LoggerContext
var localScriptPath string

//hasScript Check if a script exists
func hasScript(script string) (string, bool) {
	files, err := ioutil.ReadDir(localScriptPath)
	var exist = false
	var name = ""

	if err != nil {
		LogErr(logContextRunner, "No scripts directory found")
	} else {
		for _, file := range files {
			if strings.Contains(file.Name(), script) {
				name = file.Name()
				exist = true
			}
		}
	}
	return name, exist
}

func runLoop(job *Job, scriptPath string) {
	cmd := exec.Command(scriptPath, job.Params...)
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		LogErr(logContextRunner, "Error occurred while reading stdout")
		job.stateChan <- JobFailed
	}
	stderr, e := cmd.StderrPipe()
	if e != nil {
		LogErr(logContextRunner, "Error occurred while reading stderr")
		job.stateChan <- JobFailed
	}

	multi := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(multi)

	err := cmd.Start()
	if err != nil {
		LogErr(logContextRunner, "Error occurred during execution [%v]", err)
		job.stateChan <- JobFailed
	} else {
		job.stateChan <- JobWorking
	}

	for scanner.Scan() {
		out := scanner.Text()
		job.outChan <- out
	}
	err = cmd.Wait()
	if err != nil {
		LogErr(logContextRunner, "Error occurred during execution [%v]", err)
		job.stateChan <- JobFailed
	}

	if cmd.ProcessState != nil {
		LogInf(logContextRunner, "Script PID[%v]", cmd.ProcessState.Pid())
		if cmd.ProcessState.Exited() {
			if cmd.ProcessState.Success() {
				job.stateChan <- JobCompleted
			} else {
				job.stateChan <- JobFailed
			}
		}
	}
	LogErr(logContextRunner, "niente..")
	job.stateChan <- JobCompleted
}

//Run put in working the script
func (job *Job) Run(script, UUID, ip string, args []string) int {
	job.Name = script
	job.UUID = UUID
	job.Params = args
	job.Status = JobSubmitted
	job.outChan = make(chan string, 1)
	job.stateChan = make(chan string, 1)
	job.ip = ip

	LogInf(logContextRunner, "Run [%v], State[%v]", script, job.Status)

	name, exist := hasScript(job.Name)
	if !exist {
		LogErr(logContextRunner, "Script [%v] does not exist", script)
		return -1
	}

	scriptPath := filepath.Join(localScriptPath, name)
	go runLoop(job, scriptPath)
	return 0
}

//Log Return the current stdout and stderr
func (job *Job) Log() chan string {
	return job.outChan
}

//State Handle the status of script
func (job *Job) State() chan string {
	return job.stateChan
}

// List Get script list that we could run.§
func List() []string {
	files, err := ioutil.ReadDir(localScriptPath)
	var scripts []string

	if err != nil {
		LogErr(logContextRunner, "No scripts directory found")
	} else {
		for _, file := range files {
			n := strings.LastIndexByte(file.Name(), '.')
			if n > 0 {
				scripts = append(scripts, file.Name()[:n])
			} else {
				scripts = append(scripts, file.Name())
			}
		}
	}
	return scripts
}

// GetSet ...
func GetSet(set string) []string {
	file, err := os.Open(filepath.Join("sets", set))
	var list []string

	if err != nil {
		LogErr(logContextRunner, "Set file not found")
		return list
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	return list
}

// SetsList ...
func SetsList() []string {
	sets, err := ioutil.ReadDir("sets")
	var list []string

	if err != nil {
		LogErr(logContextRunner, "No sets dir found")
		return list
	}

	for _, set := range sets {
		list = append(list, set.Name())
	}

	return list
}

// RunnerInit ...
func RunnerInit(cm *ConfigModule) {
	localScriptPath = cm.Get("runner", "script_path", "scripts")
	logContextRunner = LoggerContext{
		name:  "RUNNER",
		level: cm.GetLogLevel("runner", 3)}

	LogInf(logContextRunner, "Start")
	LogInf(logContextRunner, "Script path[%v]", localScriptPath)
}
