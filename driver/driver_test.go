package driver

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
)

const (
	BASEDIR                      = "./test"
	STATEPATH                    = "./test/state/test-local-persist.json"
	DATAPATH                     = "./test/data"
	VOLUME1_NAME                 = "test-volume-1"
	VOLUME1_REQUESTED_MOUNTPOINT = "volume-1"
	VOLUME1_ACTUAL_MOUNTPOINT    = "./test/data/volume-1"
	VOLUME2_NAME                 = "test-volume-2"
	VOLUME2_REQUESTED_MOUNTPOINT = ""
	VOLUME2_ACTUAL_MOUNTPOINT    = "./test/data/test-volume-2"
)

type fields struct {
	Name      string
	volumes   map[string]string
	mutex     *sync.Mutex
	debug     bool
	statePath string
	dataPath  string
}

func returnFieldsEmptyVolume() fields {
	vol := make(map[string]string)

	f := fields{
		Name:      "local-persist-test",
		volumes:   vol,
		mutex:     &sync.Mutex{},
		statePath: STATEPATH,
		dataPath:  DATAPATH,
	}
	return f
}

func returnFieldsOneVolume() fields {
	vol := make(map[string]string)

	vol[VOLUME1_NAME] = VOLUME1_ACTUAL_MOUNTPOINT

	f := fields{
		Name:      "local-persist-test",
		volumes:   vol,
		mutex:     &sync.Mutex{},
		statePath: STATEPATH,
		dataPath:  DATAPATH,
	}
	return f
}

func returnFieldsTwoVolumes() fields {
	vol := make(map[string]string)

	vol[VOLUME1_NAME] = VOLUME1_ACTUAL_MOUNTPOINT
	vol[VOLUME2_NAME] = VOLUME2_ACTUAL_MOUNTPOINT

	f := fields{
		Name:      "local-persist-test",
		volumes:   vol,
		mutex:     &sync.Mutex{},
		statePath: STATEPATH,
		dataPath:  DATAPATH,
	}
	return f
}

var volume1 = volume.Volume{
	Name:       VOLUME1_NAME,
	Mountpoint: VOLUME1_ACTUAL_MOUNTPOINT,
}

var volume2 = volume.Volume{
	Name:       VOLUME2_NAME,
	Mountpoint: VOLUME2_ACTUAL_MOUNTPOINT,
}

func Test_localPersistDriver_Create(t *testing.T) {

	volume1_option := make(map[string]string)
	volume1_option["mountpoint"] = VOLUME1_REQUESTED_MOUNTPOINT

	volume3_option := make(map[string]string)
	volume3_option["mountpoint"] = "../../joe"

	type args struct {
		req *volume.CreateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Create volume with mountpoint option",
			fields: returnFieldsEmptyVolume(),
			args: args{
				req: &volume.CreateRequest{
					Name:    VOLUME1_NAME,
					Options: volume1_option,
				},
			},
			wantErr: false,
		},
		{
			name:   "Create volume without mountpoint option",
			fields: returnFieldsEmptyVolume(),
			args: args{
				&volume.CreateRequest{
					Name: VOLUME2_NAME,
				},
			},
			wantErr: false,
		},
		{
			name:   "Create volume outside mountpoint option",
			fields: returnFieldsEmptyVolume(),
			args: args{
				req: &volume.CreateRequest{
					Name:    "i-try-to-escape",
					Options: volume3_option,
				},
			},
			wantErr: false,
		},
		{
			name:   "Create volume outside basedir mountpoint option",
			fields: returnFieldsEmptyVolume(),
			args: args{
				&volume.CreateRequest{
					Name: VOLUME2_NAME,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			if err := driver.Create(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("localPersistDriver.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_localPersistDriver_Remove(t *testing.T) {
	type fields struct {
		Name      string
		volumes   map[string]string
		mutex     *sync.Mutex
		debug     bool
		statePath string
		dataPath  string
	}
	type args struct {
		req *volume.RemoveRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			if err := driver.Remove(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("localPersistDriver.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_localPersistDriver_Mount(t *testing.T) {
	type fields struct {
		Name      string
		volumes   map[string]string
		mutex     *sync.Mutex
		debug     bool
		statePath string
		dataPath  string
	}
	type args struct {
		req *volume.MountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *volume.MountResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			got, err := driver.Mount(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("localPersistDriver.Mount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("localPersistDriver.Mount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_localPersistDriver_Path(t *testing.T) {
	type fields struct {
		Name      string
		volumes   map[string]string
		mutex     *sync.Mutex
		debug     bool
		statePath string
		dataPath  string
	}
	type args struct {
		req *volume.PathRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *volume.PathResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			got, err := driver.Path(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("localPersistDriver.Path() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("localPersistDriver.Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_localPersistDriver_Unmount(t *testing.T) {
	type fields struct {
		Name      string
		volumes   map[string]string
		mutex     *sync.Mutex
		debug     bool
		statePath string
		dataPath  string
	}
	type args struct {
		req *volume.UnmountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			if err := driver.Unmount(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("localPersistDriver.Unmount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_localPersistDriver_Capabilities(t *testing.T) {
	type fields struct {
		Name      string
		volumes   map[string]string
		mutex     *sync.Mutex
		debug     bool
		statePath string
		dataPath  string
	}
	tests := []struct {
		name   string
		fields fields
		want   *volume.CapabilitiesResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			if got := driver.Capabilities(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("localPersistDriver.Capabilities() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_localPersistDriver_exists(t *testing.T) {
	type fields struct {
		Name      string
		volumes   map[string]string
		mutex     *sync.Mutex
		debug     bool
		statePath string
		dataPath  string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			if got := driver.exists(tt.args.name); got != tt.want {
				t.Errorf("localPersistDriver.exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_localPersistDriver_volume(t *testing.T) {
	type fields struct {
		Name      string
		volumes   map[string]string
		mutex     *sync.Mutex
		debug     bool
		statePath string
		dataPath  string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *volume.Volume
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			if got := driver.volume(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("localPersistDriver.volume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_localPersistDriver_findExistingVolumesFromStateFile(t *testing.T) {
	type fields struct {
		Name      string
		volumes   map[string]string
		mutex     *sync.Mutex
		debug     bool
		statePath string
		dataPath  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			got, err := driver.findExistingVolumesFromStateFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("localPersistDriver.findExistingVolumesFromStateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("localPersistDriver.findExistingVolumesFromStateFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_localPersistDriver_saveState(t *testing.T) {
	type fields struct {
		Name      string
		volumes   map[string]string
		mutex     *sync.Mutex
		debug     bool
		statePath string
		dataPath  string
	}
	type args struct {
		volumes map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			if err := driver.saveState(tt.args.volumes); (err != nil) != tt.wantErr {
				t.Errorf("localPersistDriver.saveState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ensureDir(t *testing.T) {
	type args struct {
		path string
		perm os.FileMode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ensureDir(tt.args.path, tt.args.perm); (err != nil) != tt.wantErr {
				t.Errorf("ensureDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_localPersistDriver_List(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    *volume.ListResponse
		wantErr bool
	}{
		{
			name:    "List empty volumes",
			fields:  returnFieldsEmptyVolume(),
			want:    &volume.ListResponse{},
			wantErr: false,
		},
		{
			name:   "List one volume",
			fields: returnFieldsOneVolume(),
			want: &volume.ListResponse{
				Volumes: []*volume.Volume{&volume1},
			},
			wantErr: false,
		},
		{
			name:   "List two volumes",
			fields: returnFieldsTwoVolumes(),
			want: &volume.ListResponse{
				Volumes: []*volume.Volume{&volume1, &volume2},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := &localPersistDriver{
				Name:      tt.fields.Name,
				volumes:   tt.fields.volumes,
				mutex:     tt.fields.mutex,
				debug:     tt.fields.debug,
				statePath: tt.fields.statePath,
				dataPath:  tt.fields.dataPath,
			}
			got, err := driver.List()

			if (err != nil) != tt.wantErr {
				t.Errorf("localPersistDriver.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("localPersistDriver.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
