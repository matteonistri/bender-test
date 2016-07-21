package main

import ("time"
        "os/exec"
        "strings"
        "io"
        "path/filepath"
        "bufio")

type JobStatus string

const (
	JOB_QUEUED     = "queued"
	JOB_NOT_FOUND  = "not found"
	JOB_QUEUE_FULL = "queue full"
	JOB_WORKING    = "working"
	JOB_FAILED     = "failed"
	JOB_COMPLETED  = "completed"
)

type Job struct {
	Name    string
	Params  string
	Uuid    string
	Created time.Time
	Status  JobStatus
	Timeout int
}

var scriptsDir string
var run bool

func GetScriptsDir() string {
	return scriptsDir
}

func SetScriptsDir(dir string) {
	scriptsDir = dir
}

func init(){
    SetScriptsDir("scripts")
}

func FakeRun(job *Job, script, uuid, args string) int{
    job.Name = script
    job.Uuid = uuid
    job.Params = args
    job.Status = JOB_WORKING

    var exit int

    if FakeHasScript(job.Name){
        run = true
        go func(){
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

//Check if a script exists
func FakeHasScript(script string) bool{
    return true
}

//Return the current stdout
func FakeLog(job *Job) string{
    buf := make([]byte, 100)
    //reading from stdout pipe
    return string(buf)
}

//Handle the status of script
func FakeState(job *Job){
    if run {
        job.Status = JOB_WORKING
    } else {
        job.Status = JOB_COMPLETED
    }

}

var cmd = exec.Command("")
var outChan = make(chan string, 1)

func Run(job *Job, script, uuid, args string) int{
    job.Name = script
    job.Uuid = uuid
    job.Params = args
    job.Status = JOB_WORKING

    var exit int

    if FakeHasScript(job.Name){
        params := strings.Split(job.Params, " ")
        script_path := filepath.Join(GetScriptsDir(), job.Name)
        cmd = exec.Command(script_path, params...)
        run = true
        go func(){
            cmd.Start()
            cmd.Wait()
            run = false
        }()
        exit = 0
    } else {
        exit = -1
    }

    return exit
}


//Return the current stdout
func Log() *chan string{
    go func(){
        stdout, _ := cmd.StdoutPipe()
        stderr, _ := cmd.StderrPipe()
        multi := io.MultiReader(stdout, stderr)
        scanner := bufio.NewScanner(multi)

        for scanner.Scan() {
            outChan <- scanner.Text()
        }
    }()

    return &outChan
}