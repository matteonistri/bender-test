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
	JobNotFound  = "not found"
	JobWorking   = "working"
	JobFailed    = "failed"
	JobCompleted = "completed"
)

// Job structure to track ran script
type Job struct {
	Name    string
	Params  []string
	UUID    string
	Created time.Time
	Status  string
	Timeout int
}

//JobInterface ..
type JobInterface interface {
	Run(name, UUID string, args []string) int
	UpdateState()
}

var scriptsDir string
var run bool

//GetScriptsDir ...
func GetScriptsDir() string {
	return scriptsDir
}

//SetScriptsDir  ...
func SetScriptsDir(dir string) {
	scriptsDir = dir
}

//FakeRun ..
func FakeRun(job *Job, script, uuid string, args []string) int {
	job.Name = script
	job.UUID = uuid
	job.Params = args
	job.Status = JobWorking

	var exit int

	if FakeHasScript(job.Name) {
		run = true
		go func() {
			time.Sleep(3 * time.Second)
			//execution...
			run = false
		}()
		exit = 0
	} else {
		exit = -1
	}

	return exit
}

//FakeHasScript Check if a script exists
func FakeHasScript(script string) bool {
	return true
}

//FakeLog Return the current stdout and stderr
func FakeLog(job *Job) string {
	buf := make([]byte, 100)
	//reading from stdout pipe
	return string(buf)
}

//FakeState Handle the status of script
func FakeState(job *Job) {
	if run {
		job.Status = JobWorking
	} else {
		job.Status = JobCompleted
	}
}

var cmd = exec.Command("")
var outChan = make(chan string, 1)
var syncChan = make(chan bool)
var endReadStart = make(chan bool)
var logContextRunner LoggerContext

//Run put in working the script
func (job *Job) Run(script, UUID string, args []string) int {
	job.Name = script
	job.UUID = UUID
	job.Params = args
	job.Status = JobWorking

	var exit int

	if name, exist := HasScript(job.Name); exist {
		scriptPath := filepath.Join(GetScriptsDir(), name)

		cmd = exec.Command(scriptPath, job.Params...)
		go Start()
		exit = 0
	} else {
		LogErr(logContextRunner, "Script does not exist")
		exit = -1
	}

	return exit
}

//Start go rutine to exe the command
func Start() {
	<-syncChan
	time.Sleep(100 * time.Millisecond)
	cmd.Start()
	LogInf(logContextRunner, "Execution started...")
	<-endReadStart
	err := cmd.Wait()
	LogInf(logContextRunner, "Execution finished")

	if err != nil {
		LogErr(logContextRunner, "Error occurred during execution")
	}
}

//HasScript Check if a script exists
func HasScript(script string) (string, bool) {
	files, err := ioutil.ReadDir(GetScriptsDir())
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

//Log Return the current stdout and stderr
func Log() *chan string {
	go func() {
		syncChan <- true
		stdout, err := cmd.StdoutPipe()
		stderr, err := cmd.StderrPipe()
		multi := io.MultiReader(stdout, stderr)
		scanner := bufio.NewScanner(multi)

		if err != nil {
			LogErr(logContextRunner, "Error occurred while reading stdout/stderr")
		}

		for scanner.Scan() {
			out := scanner.Text()
			outChan <- out
		}

		endReadStart <- true
		LogDeb(logContextRunner, "finished reading, sent sync to chan")
		endReadLog <- true
		LogDeb(logContextRunner, "finished reading, sent sync to chan tmp")
	}()

	return &outChan
}

//State Handle the status of script
func (job *Job) State() {
	if cmd.ProcessState == nil {
		job.Status = JobWorking
	} else if cmd.ProcessState.Success() {
		job.Status = JobCompleted
	} else {
		job.Status = JobFailed
	}
}

// List Get script list that we could run.§
func List() []string {
	files, err := ioutil.ReadDir(GetScriptsDir())
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

func init() {
	SetScriptsDir("scripts")
	logContextRunner = LoggerContext{
		name:  "RUNNER",
		level: 3}
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
