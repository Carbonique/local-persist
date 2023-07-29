package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
	log "github.com/sirupsen/logrus"
)

const STATEFILE = "local-persist.json"

type localPersistDriver struct {
	Name          string
	volumes       map[string]string
	mutex         *sync.Mutex
	stateFilePath string
	dataPath      string
}

type saveData struct {
	State map[string]string `json:"state"`
}

func NewLocalPersistDriver(statePath string, dataPath string) (*localPersistDriver, error) {
	log.Info("Starting")
	debug := os.Getenv("DEBUG")
	if ok, _ := strconv.ParseBool(debug); ok {
		log.SetLevel(log.DebugLevel)
	}

	driver := localPersistDriver{
		Name:          "local-persist",
		volumes:       map[string]string{},
		mutex:         &sync.Mutex{},
		stateFilePath: path.Join(statePath, STATEFILE),
		dataPath:      dataPath,
	}

	var err error

	err = ensureDir(statePath, 0700)
	if err != nil {
		return nil, err
	}

	err = ensureDir(dataPath, 0755)
	if err != nil {
		return nil, err
	}

	driver.volumes, _ = driver.findExistingVolumesFromStatefile()
	log.Infof("Found %d volumes on startup", len(driver.volumes))
	return &driver, nil
}

func (driver *localPersistDriver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
	log.Debug("Get called")

	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	if !driver.exists(req.Name) {
		log.Errorf("Could not find %s", req.Name)

		return &volume.GetResponse{}, fmt.Errorf("no volume found with the name %s", req.Name)
	}

	log.Debugf("Found %s", req.Name)

	return &volume.GetResponse{Volume: driver.volume(req.Name)}, nil
}

func (driver *localPersistDriver) List() (*volume.ListResponse, error) {
	log.Debug("List called")

	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	var volumes []*volume.Volume
	for name := range driver.volumes {
		volumes = append(volumes, driver.volume(name))
	}

	log.Debugf("Found %d volumes", len(driver.volumes))

	return &volume.ListResponse{Volumes: volumes}, nil
}

func (driver *localPersistDriver) Create(req *volume.CreateRequest) error {
	log.Debug("Create called")

	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	if driver.exists(req.Name) {
		return fmt.Errorf("the volume %s already exists", req.Name)
	}

	mountpoint := req.Options["mountpoint"]

	switch {
	case mountpoint == "":
		mountpoint = path.Join(driver.dataPath, req.Name)
		log.Debugf("No %s option provided. Setting mountpoint to %s", "mountpoint", mountpoint)

	case mountpoint != "":
		mountpoint = path.Join(driver.dataPath, mountpoint)
		log.Debugf("Mountpoint is %s", mountpoint)

	}

	isSubDir, err := isSubDir(driver.dataPath, mountpoint)
	if err != nil {
		return err
	}
	if !isSubDir {
		return err
	}

	err = ensureDir(mountpoint, 0755)
	if err != nil {
		return err
	}

	log.Debugf("Ensuring directory %s exists", mountpoint)

	if err != nil {
		return fmt.Errorf("%17s could not create directory %s", " ", mountpoint)
	}

	driver.volumes[req.Name] = mountpoint
	e := driver.saveState(driver.volumes)
	if e != nil {
		return fmt.Errorf("error %s", e)
	}

	log.Infof("Created volume %s with mountpoint %s", req.Name, mountpoint)

	return nil
}

func (driver *localPersistDriver) Remove(req *volume.RemoveRequest) error {
	log.Debug("Remove called")

	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	_, ok := driver.volumes[req.Name]
	// If the key exists
	if !ok {
		return fmt.Errorf("error deleting volume %s failed as it does not exist", req.Name)
	}
	delete(driver.volumes, req.Name)

	err := driver.saveState(driver.volumes)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}

	log.Infof("Removed volume %s", req.Name)

	return nil
}

func (driver *localPersistDriver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
	log.Debug("Mount called")

	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	p, ok := driver.volumes[req.Name]

	if !ok {
		return &volume.MountResponse{}, fmt.Errorf("volume %s not found", req.Name)
	} // Now check if the path still exists
	f, err := os.Stat(p)

	// If the path does not exist
	if errors.Is(err, fs.ErrNotExist) {
		return &volume.MountResponse{}, fmt.Errorf("Path %s for volume %s not found", p, req.Name)
	}

	// If the path is a file
	if f != nil && !f.IsDir() {
		return &volume.MountResponse{}, fmt.Errorf("Path %s for volume %s is a file, not a directory", p, req.Name)
	}

	log.Debugf("Mounted %s", req.Name)

	return &volume.MountResponse{Mountpoint: p}, nil
}

func (driver *localPersistDriver) Path(req *volume.PathRequest) (*volume.PathResponse, error) {
	log.Debug("Mount called")

	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	v, ok := driver.volumes[req.Name]
	if !ok {
		return &volume.PathResponse{}, fmt.Errorf("volume %s not found", req.Name)
	}
	log.Debugf("Returned path %s", v)

	return &volume.PathResponse{Mountpoint: driver.volumes[req.Name]}, nil
}

func (driver *localPersistDriver) Unmount(req *volume.UnmountRequest) error {
	log.Debug("Unmount called")

	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	_, ok := driver.volumes[req.Name]
	if !ok {
		return fmt.Errorf("volume %s not found", req.Name)
	}

	log.Infof("Unmounted %s", req.Name)

	return nil
}

func (driver *localPersistDriver) Capabilities() *volume.CapabilitiesResponse {
	log.Debug("Capabilities called")

	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "local"}}
}

func (driver *localPersistDriver) exists(name string) bool {
	return driver.volumes[name] != ""
}

func (driver *localPersistDriver) volume(name string) *volume.Volume {
	return &volume.Volume{
		Name:       name,
		Mountpoint: driver.volumes[name],
	}
}

func (driver *localPersistDriver) findExistingVolumesFromStatefile() (map[string]string, error) {
	path := driver.stateFilePath
	fileData, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{}, err
	}

	var data saveData
	e := json.Unmarshal(fileData, &data)
	if e != nil {
		return map[string]string{}, e
	}

	return data.State, nil
}

func (driver *localPersistDriver) saveState(volumes map[string]string) error {
	data := saveData{
		State: volumes,
	}

	fileData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(driver.stateFilePath, fileData, 0600)
}

func ensureDir(path string, perm os.FileMode) error {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Debugf("Trying to create path: %s with permissions: %o", path, perm)
		err := os.MkdirAll(path, perm)
		if err != nil {
			return err
		}
		return err
	}

	return nil
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isSubDir(basepath string, targetpath string) (bool, error) {
	// TODO: is there a differnt way to this instead of comparing strings
	absBasepath, err := filepath.Abs(basepath)
	if err != nil {
		return false, err
	}
	log.Debugf("absolute base path is %s", absBasepath)

	absTargetpath, err := filepath.Abs(targetpath)
	if err != nil {
		return false, err
	}
	log.Debugf("absolute target path is %s", absTargetpath)
	var isSubdir bool
	switch {

	//If they are the same path, absTargetpath is by definition not a subdirectory
	case absBasepath == absTargetpath:
		isSubdir = false

	// Now test whether the targetpath is prefixed by the basepath
	case strings.HasPrefix(absTargetpath, absBasepath):
		isSubdir = true
	}

	log.Debugf("%s is subpath of %s: %v", absTargetpath, absBasepath, isSubdir)

	if !isSubdir {
		return isSubdir, fmt.Errorf("targetpath %s is not relative to basepath %s", absTargetpath, absBasepath)
	}
	return isSubdir, nil
}
