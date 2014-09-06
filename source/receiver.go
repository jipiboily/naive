package main

import(
  "os"
  "os/exec"
  "log"
  "fmt"
  "io/ioutil"
)

func main() {
  release := Release { Repository: os.Args[1],
                      NewRevision: os.Args[2],
                      Owner: os.Args[3],
                      Fingerprint: os.Args[4] }

  log.Printf("Project directory: %s", release.ProjectDirectory())
  log.Printf("Log directory: %s", release.LogDirectory())
  log.Printf("History log path: %s", release.HistoryLogPath())
  log.Printf("Release directory: %s", release.ReleaseDirectory())
  log.Printf("Current release directory: %s", release.CurrentReleaseDirectory())
  log.Printf("New release directory: %s", release.NewReleaseDirectory())

  release.EnsureDirectoryStructureExists()
  release.WriteNewReleaseFiles()
  release.Build()
}

type Release struct {
  Repository, NewRevision, Owner, Fingerprint string;
}

func RootDirectory() string {
  return "/paas"
}

func (r Release) ProjectDirectory() string {
  return fmt.Sprintf("%s/%s", RootDirectory(), r.Repository)
}

func (r Release) LogDirectory() string {
  return fmt.Sprintf("%s/logs", r.ProjectDirectory())
}

func (r Release) HistoryLogPath() string {
  return fmt.Sprintf("%s/history.log", r.LogDirectory())
}

func (r Release) ReleaseDirectory() string {
  return fmt.Sprintf("%s/releases", r.ProjectDirectory())
}

func (r Release) CurrentReleaseDirectory() string {
  return fmt.Sprintf("%s/current", r.ReleaseDirectory())
}

func (r Release) NewReleaseDirectory() string {
  return fmt.Sprintf("%s/%s", r.ReleaseDirectory(), r.NewRevision)
}

func (r Release) EnsureDirectoryStructureExists() {
  log.Println("Ensuring directory structure exists...")
  // TODO: remove the nil output and only output when errors occur
  log.Println(os.MkdirAll(r.NewReleaseDirectory(), 0750))
  log.Println(os.MkdirAll(r.CurrentReleaseDirectory(), 0750))
  log.Println(os.MkdirAll(r.LogDirectory(), 0750))
}

func (r Release) WriteTarFile() string {
  bytes, err := ioutil.ReadAll(os.Stdin)
  tarPath := fmt.Sprintf("%s/app.tar", r.NewReleaseDirectory())
  err = ioutil.WriteFile(tarPath, bytes, 0640)
  if err != nil { panic(err) }
  return tarPath
}

func (r Release) WriteNewReleaseFiles() {
  tarPath := r.WriteTarFile()
  // TODO: clear the current directory before untarring
  // exec.Command("rm", "-rf", r.CurrentReleaseDirectory())
  _, err := exec.Command("/bin/tar", "-xf", tarPath, "-C", r.CurrentReleaseDirectory()).Output()
  if err != nil { panic(err) }
}

func (r Release) Build() {
  if r.HasDockerFile(){
    log.Println("Launching the docker build...")
    // TODO: error handling
    // TODO: output with `log`?
    cmd := exec.Command("/usr/bin/docker", "build", "-t", r.Repository, r.CurrentReleaseDirectory())
    RunCommand(cmd)
  } else {
    log.Println("Buildstep stuff...")
  }
}

func (r Release) HasDockerFile() bool {
  path := fmt.Sprintf("%s/Dockerfile", r.CurrentReleaseDirectory())
  finfo, err := os.Stat(path)
  if err == nil && !finfo.IsDir() {
    log.Println("Dockerfile found...")
    return true
  } else {
    log.Println("No Dockerfile found...:(")
    log.Println("Trying to build with a build pack...")
    return false
  }

// Abstract to run command with output and error "management"
func RunCommand(cmd *exec.Cmd) {
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  err := cmd.Run()
  if err != nil { panic(err) }
}