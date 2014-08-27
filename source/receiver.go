package main

import(
  "os"
  "log"
  "fmt"
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
  log.Printf("New release directory: %s", release.NewReleaseDirectory())

  release.EnsureDirectoryStructureExists()
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

func (r Release) NewReleaseDirectory() string {
  return fmt.Sprintf("%s/%s", r.ReleaseDirectory(), r.NewRevision)
}

func (r Release) EnsureDirectoryStructureExists() {
  log.Println("Ensuring directory structure exists...")
  log.Println(os.MkdirAll(r.NewReleaseDirectory(), 0700))
  log.Println(os.MkdirAll(r.LogDirectory(), 0700))
}
