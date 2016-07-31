package main

import (
	"bufio"
	"fmt"
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
	Name    string
	Params  []string
	UUID    string
	Created time.Time
	Status  string
	pid     int
	Timeout int
}

//JobInterface ...
type JobInterface interface {
	Run(name, UUID string, args []string) int
	State() *chan string
	Log() *chan string
}

var outChan = make(chan string, 1)
var stateChan = make(chan string, 1)
var cmdStartChannel = make(chan *exec.Cmd)
var cmdStateChannel = make(chan *exec.Cmd)
var cmdSyncChannel = make(chan *exec.Cmd)
var cmdStateErrorChannel = make(chan string, 1)
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

//Start go rutine to exe the command
func runScritpLoop() {
	c := <-cmdSyncChannel
	c.Start()
	LogInf(logContextRunner, "Execution started...")

	// Start to track state of script
	cmdStateChannel <- c

	err := c.Wait()
	if err != nil {
		LogErr(logContextRunner, "Error occurred during execution [%v]", err)
		cmdStateErrorChannel <- JobFailed
	}
}

func logRunLoop() {
	//Wait start sync from run func
	c := <-cmdStartChannel
	// Start script exec
	cmdSyncChannel <- c

	stdout, err := c.StdoutPipe()
	if err != nil {
		LogErr(logContextRunner, "Error occurred while reading stdout")
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		LogErr(logContextRunner, "Error occurred while reading stderr")
	}

	multi := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(multi)
	for scanner.Scan() {
		out := scanner.Text()
		outChan <- out
	}
}

func stateRunLoop() {
	//Wait start sync from run func
	c := <-cmdStateChannel

	for {
		select {
		case s := <-cmdStateErrorChannel:
			LogErr(logContextRunner, "Script Error occurred")
			stateChan <- s
			return
		default:
			if c.ProcessState == nil {
				stateChan <- JobWorking
			} else {
				fmt.Println("Process", c.ProcessState.Pid())
				fmt.Println("Process", c.ProcessState.String())
				fmt.Println("Process", c.ProcessState.Success())
				fmt.Println("Process", c.ProcessState.Exited())
				fmt.Println("Process", c.ProcessState.SystemTime())
				fmt.Println("Process", c.ProcessState.UserTime())
				stateChan <- JobCompleted
				return
			}
		}
	}
}

//Run put in working the script
func (job *Job) Run(script, UUID string, args []string) int {
	job.Name = script
	job.UUID = UUID
	job.Params = args
	job.Status = JobSubmitted

	LogInf(logContextRunner, "Run [%v], State[%v]", script, job.Status)

	name, exist := hasScript(job.Name)
	if !exist {
		LogErr(logContextRunner, "Script [%v] does not exist", script)
		return -1
	}

	LogInf(logContextRunner, "Prepare exec of [%v]", script)
	scriptPath := filepath.Join(localScriptPath, name)
	cmd := exec.Command(scriptPath, job.Params...)

	go runScritpLoop()
	go logRunLoop()
	go stateRunLoop()

	LogInf(logContextRunner, "Start exec of [%v]", script)
	cmdStartChannel <- cmd
	return 0
}

//Log Return the current stdout and stderr
func (job *Job) Log() *chan string {
	return &outChan
}

//State Handle the status of script
func (job *Job) State() *chan string {
	return &stateChan
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
