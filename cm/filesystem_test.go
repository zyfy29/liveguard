package cm

import (
	"os"
	"testing"
)

func Test_copyFile(t *testing.T) {
	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func() error
		cleanup func() error
	}{
		{
			name: "successful copy",
			args: args{
				src:  "test_src.txt",
				dest: "test_dest/test_dest.txt",
			},
			wantErr: false,
			setup: func() error {
				return os.WriteFile("test_src.txt", []byte("test content"), 0644)
			},
			cleanup: func() error {
				os.Remove("test_src.txt")
				return os.RemoveAll("test_dest")
			},
		},
		{
			name: "source file does not exist",
			args: args{
				src:  "non_existent_src.txt",
				dest: "test_dest/non_existent_dest.txt",
			},
			wantErr: true,
			setup:   func() error { return nil },
			cleanup: func() error { return nil },
		},
		{
			name: "destination directory does not exist",
			args: args{
				src:  "test_src2.txt",
				dest: "non_existent_dir/test_dest.txt",
			},
			wantErr: false,
			setup: func() error {
				return os.WriteFile("test_src2.txt", []byte("another test content"), 0644)
			},
			cleanup: func() error {
				os.Remove("test_src2.txt")
				return os.RemoveAll("non_existent_dir")
			},
		},
		{
			name: "overwrite existing file",
			args: args{
				src:  "test_src3.txt",
				dest: "test_dest/test_dest_overwrite.txt",
			},
			wantErr: false,
			setup: func() error {
				os.MkdirAll("test_dest", os.ModePerm)
				os.WriteFile("test_src3.txt", []byte("overwrite content"), 0644)
				return os.WriteFile("test_dest/test_dest_overwrite.txt", []byte("old content"), 0644)
			},
			cleanup: func() error {
				os.Remove("test_src3.txt")
				return os.RemoveAll("test_dest")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("setup() error: %v", err)
				}
			}

			if err := copyFile(tt.args.src, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.cleanup != nil {
				if err := tt.cleanup(); err != nil {
					t.Fatalf("cleanup() error: %v", err)
				}
			}
		})
	}
}

func TestGetRandomDataFilePathWithNameAndExt(t *testing.T) {
	type args struct {
		name string
		ext  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				name: "test",
				ext:  "txt",
			},
		},
		{
			name: "empty name",
			args: args{
				name: "",
				ext:  "txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := GetRandomDataFilePathWithNameAndExt(tt.args.name, "txt")
			t.Logf("Random data file path: %s", res)
		})
	}
}

func TestGetDataFilePath(t *testing.T) {
	type args struct {
		parts []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				parts: []string{"test", "test2", "test3.txt"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Data file path: %s", GetDataFilePath(tt.args.parts...))
		})
	}
}

func TestProjectRoot(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test GetProjectRoot",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := GetProjectRoot()
			t.Logf("Project root: %s", res)
		})
	}
}

func TestGetAbsDataDir(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test getAbsDataDir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := getAbsDataDir()
			t.Logf("Data directory: %s", res)
		})
	}
}
