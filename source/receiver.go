package main

import(
  "os"
  "os/exec"
  "log"
  "fmt"
  "io/ioutil"
  "strconv"
  "text/template"
  "strings"
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
  log.Printf("New release name: %s", release.NewReleaseName())
  log.Printf("Previous release name: %s", release.PreviousReleaseName())

  release.EnsureDirectoryStructureExists()
  release.WriteNewReleaseFiles()
  release.Build()
  // TODO: move that later so we can have zero down time.
  //       What is missing for that is multi-port, as the
  //       same port can't be used more than once, obviously.
  release.RemoveOldContainer()
  release.Run()
  release.SwitchRoute() // Switch nginx to route to this new one
  release.BumpReleaseId()
}

type Release struct {
  Repository, NewRevision, Owner, Fingerprint string;
}

func RootDirectory() string {
  return "/paas"
}

func (r Release) ProjectDirectory() string {
  return fmt.Sprintf("%s/apps/%s", RootDirectory(), r.Repository)
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

func (r Release) ReleaseVersionFilePath() string {
  return fmt.Sprintf("%s/%s", r.ProjectDirectory(), "current_release_id")
}

func (r Release) NginxSitesEnabledPath() string {
  return fmt.Sprintf("%s/config/nginx/sites-enabled", RootDirectory())
}

func (r Release) NginxLogsPath() string {
  return fmt.Sprintf("%s/logs/nginx/", RootDirectory())
}

func (r Release) NewReleaseName() string {
  releaseID := strconv.FormatInt(r.NewReleaseId(), 10);
  return fmt.Sprintf("%s-%s", r.Repository, releaseID)
}

func (r Release) PreviousReleaseName() string {
  releaseID := strconv.FormatInt(r.PreviousReleaseId(), 10);
  return fmt.Sprintf("%s-%s", r.Repository, releaseID)
}

func (r Release) NewReleaseId() int64 {
  return r.PreviousReleaseId() + 1
}

func (r Release) PreviousReleaseId() int64 {
  if _, err := os.Stat(r.ReleaseVersionFilePath()); err == nil {
    content,_ := ioutil.ReadFile(r.ReleaseVersionFilePath())
    releaseId, err := strconv.ParseInt(string(content), 10, 64)
    if err != nil { panic(err) }
    return releaseId;
  } else {
    return 0;
  }
}

func (r Release) BumpReleaseId() {
  newReleaseId := strconv.FormatInt(r.NewReleaseId() + 1, 10)
  log.Printf("Writing new release version: %s", newReleaseId)
  releaseId := strconv.FormatInt(r.PreviousReleaseId() + 1, 10)
  err := ioutil.WriteFile(r.ReleaseVersionFilePath(), []byte(releaseId), 0640)
  if err != nil { panic(err) }
}

func (r Release) EnsureDirectoryStructureExists() {
  log.Println("Ensuring directory structure exists...")
  // TODO: remove the nil output and only output when errors occur
  // TODO: improve security on this...I would prefer 0750 at most.
  log.Println(os.MkdirAll(r.NewReleaseDirectory(), 0755))
  log.Println(os.MkdirAll(r.CurrentReleaseDirectory(), 0755))
  log.Println(os.MkdirAll(r.LogDirectory(), 0755))
  log.Println(os.MkdirAll(r.NginxSitesEnabledPath(), 0755))
  log.Println(os.MkdirAll(r.NginxLogsPath(), 0755))
}

func (r Release) WriteTarFile() string {
  bytes, err := ioutil.ReadAll(os.Stdin)
  tarPath := fmt.Sprintf("%s/app.tar", r.NewReleaseDirectory())
  err = ioutil.WriteFile(tarPath, bytes, 0640)
  if err != nil { panic(err) }
  return tarPath
}

// This will write the files for your new release in the "current" directory
func (r Release) WriteNewReleaseFiles() {
  tarPath := r.WriteTarFile()
  // TODO: clear the current directory before untarring
  // exec.Command("rm", "-rf", r.CurrentReleaseDirectory())
  _, err := exec.Command("/bin/tar", "-xf", tarPath, "-C", r.CurrentReleaseDirectory()).Output()
  if err != nil { panic(err) }
}

// This will build a new container for your app
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

// Inspect the new release directory to see if there is a Dockerfile that can be used
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
}

// Starting the container, with the environment variables
func (r Release) Run() {
  // TODO: We will need to use $PORT to start stuff inside the container, for now, assuming 3000
  cmd := exec.Command("/usr/bin/docker", "run", "-d", "-p", "3000:3000", "--name", r.NewReleaseName(), r.Repository)
  RunCommand(cmd)
}

// Abstract to run command with output and error "management"
func RunCommand(cmd *exec.Cmd) {
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  err := cmd.Run()
  if err != nil { panic(err) }
}

func (r Release) SwitchRoute() {
  r.GenerateNginxConf()
  r.EnsureNginxConfLink()
  exec.Command("/usr/bin/docker", "start", "nginx").Output()
  exec.Command("/usr/bin/docker", "restart", "nginx").Output()
}

func (r Release) EnsureNginxConfLink() {
  cmd := exec.Command("/bin/cp", r.NginxConfPath(), r.NginxSymbolicLinkPath())
  cmd.Run()
}

func (r Release) RemoveOldContainer() {
  log.Println("Removing old container.")
  cmd := exec.Command("/usr/bin/docker", "kill", r.PreviousReleaseName())
  cmd.Run() // We don't care if it fails, at least, I am assuming the app is new or was closed
  // TODO: This should be queued, somehow, and run a little later, to let the kill do its jobs
  // cmd = exec.Command("/usr/bin/docker", "rm", r.PreviousReleaseName())
  // RunCommand(cmd)
}

func (r Release) NewReleaseIP() string {
  cmd := exec.Command("/usr/bin/docker", "inspect", "--format", "'{{.NetworkSettings.IPAddress}}'", r.NewReleaseName())
  ipBytes, err := cmd.Output()
  if err != nil { panic(err) }
  return string(ipBytes)
}

func (r Release) NginxProxyPass() string {
  cleanedIp := strings.Replace(r.NewReleaseIP(), "'", "", -1)
  cleanedIp = strings.Replace(cleanedIp, "\n", "", -1)
  return fmt.Sprintf("http://%s%s", cleanedIp, ":3000")
}

func (r Release) NginxServerName() string {
  return fmt.Sprintf("server_name localhost %s.jipiboily.net %s;", r.Repository, "http://app.jipiboily.com")
}

func (r Release) NginxConfPath() string {
  return fmt.Sprintf("%s/nginx.conf", r.ProjectDirectory())
}

func (r Release) NginxTemplatePath() string {
  return "/vagrant/templates/nginx-site.conf"
}

func (r Release) NginxSymbolicLinkPath() string {
  return fmt.Sprintf("/paas/config/nginx/sites-enabled/%s.conf", r.Repository)
}

func (r Release) GenerateNginxConf() {
  log.Println("Generating routing conf.")
  tmpl, errTmpl := template.ParseFiles(r.NginxTemplatePath())
  if errTmpl != nil { panic(errTmpl) }

  tmplFile, errTmplFile := os.OpenFile(r.NginxConfPath(), os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0644)
  if errTmplFile != nil { panic(errTmplFile) }

  errTmpl = tmpl.Execute(tmplFile, r)
  if errTmpl != nil { panic(errTmpl) }
}
