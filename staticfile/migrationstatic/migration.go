// Code generated for package migrationstatic by go-bindata DO NOT EDIT. (@generated)
// sources:
// migrations/1.0.0/20210907143608_create_v2ray_endpoint_table.down.sql
// migrations/1.0.0/20210907143608_create_v2ray_endpoint_table.up.sql
package migrationstatic

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _migrations10020210907143608_create_v2ray_endpoint_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4a\x29\xca\x2f\x50\x28\x49\x4c\xca\x49\x55\xc8\x4c\x53\x48\xad\xc8\x2c\x2e\x29\x56\x28\x33\x2a\x4a\xac\x8c\x4f\xcd\x4b\x29\xc8\xcf\xcc\x2b\xb1\x06\x04\x00\x00\xff\xff\x30\xdd\x7c\xb5\x24\x00\x00\x00")

func migrations10020210907143608_create_v2ray_endpoint_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations10020210907143608_create_v2ray_endpoint_tableDownSql,
		"migrations/1.0.0/20210907143608_create_v2ray_endpoint_table.down.sql",
	)
}

func migrations10020210907143608_create_v2ray_endpoint_tableDownSql() (*asset, error) {
	bytes, err := migrations10020210907143608_create_v2ray_endpoint_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/1.0.0/20210907143608_create_v2ray_endpoint_table.down.sql", size: 36, mode: os.FileMode(420), modTime: time.Unix(1631498292, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations10020210907143608_create_v2ray_endpoint_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x93\xc1\x6e\xb3\x30\x10\x84\xef\x3c\xc5\xde\x92\x48\x39\xfc\xc9\x1f\xf5\xd2\x87\xb1\x1c\x7b\x93\x58\x98\xb5\xb5\x5e\x5a\x78\xfb\x0a\x52\x12\xea\x9a\xa8\xf8\x84\xc4\xcc\x37\x78\xd0\x18\x46\x2d\x08\xa2\xcf\x1e\xc1\x5d\x80\x82\x00\x76\x2e\x49\x82\x8f\x23\xeb\x5e\x21\xd9\x18\x1c\x49\xb5\xad\x00\x00\x9c\x85\xf9\x71\x24\x78\x45\x1e\x6d\xd4\x7a\x5f\x4d\x2f\x4c\xa0\x24\xac\x1d\x49\xc6\x51\xb1\xae\xe6\x84\xc8\xae\xd1\xdc\x43\x8d\x3d\xe8\x56\x82\x23\xc3\xd8\x20\xc9\x7e\x94\x19\x1f\x5a\xbb\x9c\x77\x17\x4d\xec\x97\x22\x1e\xee\xf9\x3c\x82\xdd\xa4\xb7\x78\xd1\xad\x17\xd8\x6c\x72\x0b\x36\x9a\xeb\x55\x96\x5b\x48\x52\x48\xf9\x29\x8a\x81\xa5\x50\xe2\x8c\x7b\x3a\xfd\xcf\x3c\x6d\x42\x56\xcf\xf6\xcb\x60\xed\x65\xae\x2a\x80\xdf\x4e\xbf\xb9\x4a\x7c\x7a\xf1\x2d\x87\xcc\x91\xc8\xc1\xca\x22\x85\x35\xa5\xe1\xd6\x4a\xfa\x88\xa5\x90\x7f\xb9\xc3\xc4\xb5\x21\x9f\x78\x56\x29\x98\x1a\xe5\xcf\x96\x7a\x7d\xca\x4d\x24\x1e\xd7\x59\xae\x1c\xcd\xca\x14\x8b\x1e\x05\x6d\xf6\x57\xb6\x87\xdd\x62\x61\xf7\x21\x2b\x71\x0d\x42\x69\x98\xbb\xf7\xaa\xfa\x1e\xbb\x23\x8b\x1d\x38\xdb\xa9\x71\x5d\xcf\x89\x0f\xc6\x40\xd9\x60\x61\x3b\xaa\xf6\x8f\x95\x15\x49\x8f\x79\x8f\xe2\x25\xd2\xf4\xb4\xbf\xef\x7a\x20\x7d\x05\x00\x00\xff\xff\xaa\xe3\x24\xc4\x81\x04\x00\x00")

func migrations10020210907143608_create_v2ray_endpoint_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations10020210907143608_create_v2ray_endpoint_tableUpSql,
		"migrations/1.0.0/20210907143608_create_v2ray_endpoint_table.up.sql",
	)
}

func migrations10020210907143608_create_v2ray_endpoint_tableUpSql() (*asset, error) {
	bytes, err := migrations10020210907143608_create_v2ray_endpoint_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/1.0.0/20210907143608_create_v2ray_endpoint_table.up.sql", size: 1153, mode: os.FileMode(420), modTime: time.Unix(1631498292, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"migrations/1.0.0/20210907143608_create_v2ray_endpoint_table.down.sql": migrations10020210907143608_create_v2ray_endpoint_tableDownSql,
	"migrations/1.0.0/20210907143608_create_v2ray_endpoint_table.up.sql":   migrations10020210907143608_create_v2ray_endpoint_tableUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"migrations": &bintree{nil, map[string]*bintree{
		"1.0.0": &bintree{nil, map[string]*bintree{
			"20210907143608_create_v2ray_endpoint_table.down.sql": &bintree{migrations10020210907143608_create_v2ray_endpoint_tableDownSql, map[string]*bintree{}},
			"20210907143608_create_v2ray_endpoint_table.up.sql":   &bintree{migrations10020210907143608_create_v2ray_endpoint_tableUpSql, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
