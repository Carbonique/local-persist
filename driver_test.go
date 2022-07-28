package main

import (
	"os"
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
)

var (
	defaultTestName       = "test-volume"
	defaultTestMountpoint = "/tmp/data/local-persist-test"
)

func TestCreate(t *testing.T) {
	driver := newLocalPersistDriver()

	defaultCreateHelper(driver, t)

	// test that a directory is created
	_, err := os.Stat(defaultTestMountpoint)
	if os.IsNotExist(err) {
		t.Error("Mountpoint directory was not created:", err.Error())
	}

	// test that volumes has one
	if len(driver.volumes) != 1 {
		t.Error("Driver should have exactly 1 volume")
	}

	defaultCleanupHelper(driver, t)

	req := &volume.CreateRequest{Name: "defaultTestName"}
	// test that options are required
	err = driver.Create(req)

	if err.Error() != "the `mountpoint` option is required" {
		t.Error(err)
	}
	defaultCleanupHelper(driver, t)
}

func TestGet(t *testing.T) {
	driver := newLocalPersistDriver()

	defaultCreateHelper(driver, t)

	req := &volume.GetRequest{Name: defaultTestName}

	res, err := driver.Get(req)
	if err != nil {
		t.Error("Should have found a volume!")
	}

	if res.Volume.Name != defaultTestName {
		t.Error("Incorrect volume name was returned")
	}

	defaultCleanupHelper(driver, t)
}

func TestList(t *testing.T) {
	driver := newLocalPersistDriver()

	name := defaultTestName + "2"
	mountpoint := defaultTestMountpoint + "2"

	defaultCreateHelper(driver, t)

	res, err := driver.List()

	if err != nil {
		t.Error("List returned error")
	}

	if len(res.Volumes) != 1 {
		t.Error("Should have found 1 volume!")
	}

	createHelper(driver, t, name, mountpoint)

	res, err = driver.List()

	if err != nil {
		t.Error("List returned error2")
	}

	if len(res.Volumes) != 2 {
		t.Error("Should have found 1 volume!")
	}

	defaultCleanupHelper(driver, t)
	cleanupHelper(driver, t, name, mountpoint)
}

func TestMount(t *testing.T) {
	driver := newLocalPersistDriver()

	defaultCreateHelper(driver, t)
	req := &volume.MountRequest{Name: defaultTestName}
	_, err := driver.Mount(req)

	if err != nil {
		t.Error("Error on mount")
	}

	err = os.Remove(defaultTestMountpoint)

	if err != nil {
		t.Error("Could not remove mountpoint")
	}

	_, err = driver.Mount(req)

	if err == nil {
		t.Error("Mountpoint was deleted but test did not error")
	}

	_, err = os.Create(defaultTestMountpoint)
	if err != nil {
		t.Error("Could create mountpoint as file")
	}

	_, err = driver.Mount(req)
	if err == nil {
		t.Error("Mountpoint is a file but test did not error")
	}
	defaultCleanupHelper(driver, t)
}

func TestUnmount(t *testing.T) {
	driver := newLocalPersistDriver()

	defaultCreateHelper(driver, t)

	reqFail := &volume.UnmountRequest{Name: defaultTestName + "does_not_exist"}
	err := driver.Unmount(reqFail)

	if err == nil {
		t.Error("test should fail as volume does not exist")
	}

	req := &volume.UnmountRequest{Name: defaultTestName}

	err = driver.Unmount(req)

	if err != nil {
		t.Error("Error on unmount")
	}

	defaultCleanupHelper(driver, t)

}
func TestPath(t *testing.T) {
	driver := newLocalPersistDriver()

	defaultCreateHelper(driver, t)

	reqFail := &volume.PathRequest{Name: defaultTestName + "does_not_exist"}
	_, err := driver.Path(reqFail)

	if err == nil {
		t.Error("test should fail as volume does not exist")
	}

	req := &volume.PathRequest{Name: defaultTestName}

	v, err := driver.Path(req)

	if err != nil {
		t.Error("error on path")
	}

	if v.Mountpoint != defaultTestMountpoint {
		t.Error("mountpoint should be equal to defaultTestMountpoint")
	}

	defaultCleanupHelper(driver, t)

}

func createHelper(driver *localPersistDriver, t *testing.T, name string, mountpoint string) {

	req := &volume.CreateRequest{
		Name: name,
		Options: map[string]string{
			"mountpoint": mountpoint,
		},
	}

	err := driver.Create(req)

	if err != nil {
		t.Error(err)
	}
}

func defaultCreateHelper(driver *localPersistDriver, t *testing.T) {
	createHelper(driver, t, defaultTestName, defaultTestMountpoint)
}

func cleanupHelper(driver *localPersistDriver, t *testing.T, name string, mountpoint string) {
	os.RemoveAll(mountpoint)

	_, err := os.Stat(mountpoint)
	if !os.IsNotExist(err) {
		t.Error("[Cleanup] Mountpoint still exists:", err.Error())
	}

	removeReq := &volume.RemoveRequest{Name: name}

	err = driver.Remove(removeReq)
	if err != nil {
		t.Error("Remove failed", err)
	}

	getReq := &volume.GetRequest{Name: name}

	//Volume should be nil, as it is deleted
	v, err := driver.Get(getReq)

	if v.Volume != nil {
		t.Error(err)
	}
}

func defaultCleanupHelper(driver *localPersistDriver, t *testing.T) {
	cleanupHelper(driver, t, defaultTestName, defaultTestMountpoint)
}
