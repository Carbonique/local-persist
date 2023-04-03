package driver

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
)

const (
	baseDir = "./test"
)

var (
	stateDir = path.Join(baseDir, "state")
	dataDir  = path.Join(baseDir, "data")
)

func init() {
	if _, err := os.Stat(baseDir); !os.IsNotExist(err) {
		os.RemoveAll(baseDir)
	}
}

func TestCreate(t *testing.T) {

	tests := map[string]struct {
		name             string
		mountpointOption string
		mountpointPath   string
	}{
		"create with mountpoint option":       {name: "test-volume", mountpointOption: "/my/directory", mountpointPath: path.Join(dataDir, "/my/directory")},
		"create with empty mountpoint option": {name: "test-volume", mountpointOption: "", mountpointPath: path.Join(dataDir, "test-volume")},
	}

	driver := createDriverHelper()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			createVolumeHelper(driver, t, tc.name, tc.mountpointOption)

			// check that directory is created
			_, err := os.Stat(tc.mountpointPath)
			if os.IsNotExist(err) {
				t.Error("Mountpoint directory was not created", err.Error())
			}

			// check that volumes has one
			if len(driver.volumes) != 1 {
				t.Error("Driver should have exactly 1 volume")
			}

			volumeCleanupHelper(driver, t, tc.name, tc.mountpointPath)
		})
	}

}

func TestGet(t *testing.T) {
	tests := map[string]struct {
		name             string
		mountpointOption string
		mountpointPath   string
		expectedName     string
	}{
		"create with mountpoint option": {name: "test-volume", mountpointOption: "/my/directory", mountpointPath: path.Join(dataDir, "/my/directory"), expectedName: "test-volume"}}

	driver := createDriverHelper()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			createVolumeHelper(driver, t, tc.name, tc.mountpointOption)

			req := &volume.GetRequest{Name: tc.name}

			res, err := driver.Get(req)
			if err != nil {
				t.Error("Should have found a volume!")
			}

			if res.Volume.Name != tc.expectedName {
				t.Error("Incorrect volume name was returned")
			}

			volumeCleanupHelper(driver, t, tc.name, tc.mountpointPath)
		})
	}

}

func TestList(t *testing.T) {

	type volume struct {
		name             string
		mountpointOption string
		mountpointPath   string
	}

	vol1 := volume{name: "test-volume-1", mountpointOption: "", mountpointPath: path.Join(dataDir, "test-volume-1")}
	vol2 := volume{name: "test-volume-2", mountpointOption: "", mountpointPath: path.Join(dataDir, "test-volume-2")}

	tests := map[string]struct {
		volumes         []volume
		expected_length int
	}{
		"list one volume should return one volume":   {volumes: []volume{vol1}, expected_length: 1},
		"list two volumes should return two volumes": {volumes: []volume{vol1, vol2}, expected_length: 2},
	}

	driver := createDriverHelper()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			// Create all volumes
			for _, vol := range tc.volumes {
				createVolumeHelper(driver, t, vol.name, vol.mountpointOption)
			}

			res, err := driver.List()

			if err != nil {
				t.Error("List returned error")
			}

			if len(res.Volumes) != tc.expected_length {
				t.Errorf("Should have found %d volume(s)!", tc.expected_length)
			}

			// Cleanup all volumes
			for _, vol := range tc.volumes {
				volumeCleanupHelper(driver, t, vol.name, vol.mountpointOption)
			}

		})
	}
}

func TestMount(t *testing.T) {
	// TODO: use table testing
	testVolumeName := "test-volume"
	testMountpoint := "test"

	driver := createDriverHelper()

	createVolumeHelper(driver, t, testVolumeName, testMountpoint)

	req := &volume.MountRequest{Name: testVolumeName}
	_, err := driver.Mount(req)

	if err != nil {
		t.Error("Error on mount")
	}

	// Remove a mountpoint, while volume still exists
	err = os.Remove(path.Join(dataDir, testMountpoint))

	if err != nil {
		t.Error("Could not remove mountpoint")
	}

	_, err = driver.Mount(req)
	if err == nil {
		t.Error("Mountpoint was deleted but test did not error")
	}

	// Test to mount an existing file (should not be possible)
	_, err = os.Create(path.Join(dataDir, testMountpoint))
	if err != nil {
		t.Error("Could not create mountpoint as file")
	}

	_, err = driver.Mount(req)
	if err == nil {
		t.Error("Mountpoint is a file but test did not error")
	}
	volumeCleanupHelper(driver, t, testVolumeName, testMountpoint)
}

func TestUnmount(t *testing.T) {

	tests := map[string]struct {
		name             string
		exists           bool
		mountpointOption string
		mountpointPath   string
		error_expected   bool
	}{
		"Unmount existing volume":   {name: "test-volume-existing", exists: true, mountpointOption: "/my/directory", mountpointPath: path.Join(dataDir, "/my/directory"), error_expected: false},
		"Mount non-existing volume": {name: "test-volume-non-existing", exists: false, mountpointOption: "", mountpointPath: path.Join(dataDir, "does-not-exist"), error_expected: true},
	}

	driver := createDriverHelper()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.exists {
				createVolumeHelper(driver, t, tc.name, tc.mountpointOption)
			}

			req := &volume.UnmountRequest{Name: tc.name}
			err := driver.Unmount(req)

			if err != nil {
				if !tc.error_expected {
					t.Error("Error on unmount")
				}
			}

			if tc.exists {
				volumeCleanupHelper(driver, t, tc.name, tc.mountpointPath)
			}
		})
	}

}
func TestPath(t *testing.T) {

	tests := map[string]struct {
		name             string
		exists           bool
		mountpointOption string
		mountpointPath   string
		error_expected   bool
	}{
		"existing volume":     {name: "test-volume-existing", exists: true, mountpointOption: "/my/directory", mountpointPath: path.Join(dataDir, "/my/directory"), error_expected: false},
		"non-existing volume": {name: "test-volume-non-existing", exists: false, mountpointOption: "", mountpointPath: path.Join(dataDir, "does-not-exist"), error_expected: true},
	}

	driver := createDriverHelper()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.exists {
				createVolumeHelper(driver, t, tc.name, tc.mountpointOption)

				req := &volume.PathRequest{Name: tc.name}

				v, err := driver.Path(req)

				if err != nil {
					t.Error("Error on path")
				}

				if v.Mountpoint != tc.mountpointPath {
					t.Error("Mountpoint should be equal to defaultTestMountpoint")
				}
				volumeCleanupHelper(driver, t, tc.name, tc.mountpointPath)
			}

			if !tc.exists {
				reqFail := &volume.PathRequest{Name: tc.name}
				_, err := driver.Path(reqFail)

				if err != nil {
					if !tc.error_expected {
						t.Error("Test should fail as volume does not exist")
					}
				}
			}

		})
	}

}

func createVolumeHelper(driver *localPersistDriver, t *testing.T, name string, mountpoint string) {

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

func volumeCleanupHelper(driver *localPersistDriver, t *testing.T, name string, mountpoint string) {
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

func createDriverHelper() *localPersistDriver {
	d, err := NewLocalPersistDriver(stateDir, dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	return d
}
