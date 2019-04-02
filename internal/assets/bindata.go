// Code generated by go-bindata.
// sources:
// assets/exe.manifest
// DO NOT EDIT!

package assets

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _assetsExeManifest = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x91\xc1\x6a\xdb\x40\x10\x86\xcf\x0d\xe4\x1d\x96\x39\x16\xac\xb5\x63\x3b\xb4\xc6\x72\x30\x81\xd0\x50\xd2\x8b\xd3\xe6\xbc\x5e\x8d\xed\x25\xda\x99\x65\x77\x64\x47\xcf\xd6\x43\x1f\xa9\xaf\x50\x24\x39\x76\x0d\x82\xa0\xd3\x0c\xfb\xcf\xff\xfd\xfa\xff\xfe\xfe\x33\xbf\x7b\xf3\xa5\xda\x63\x4c\x8e\x29\x87\x51\x36\x04\x85\x64\xb9\x70\xb4\xcd\xe1\xe7\xf3\xc3\xe0\x0b\xa8\x24\x86\x0a\x53\x32\x61\x0e\x35\x26\xb8\x5b\x5c\x5f\xcd\x4d\x4a\xe8\xd7\x65\xad\xde\x7c\x49\x29\x87\x2a\xd2\x2c\xd9\x1d\x7a\x93\x06\xde\xd9\xc8\x89\x37\x32\xb0\xec\x67\x26\xf9\x6c\x3f\x02\xe5\x0d\xb9\x0d\x26\xf9\x75\xe1\xd6\xca\x9b\x37\xfb\xf1\x87\x47\xc6\xb0\xb8\xbe\xfa\x74\xb2\x7e\x2c\x90\xc4\x49\x7d\xc1\xdf\x7c\xa0\x42\x64\x8b\x29\x71\x5c\x46\xbb\x73\x82\x56\xaa\x88\x39\x7c\x06\x45\xc6\x63\x0e\x2b\xf6\xf8\x50\xd1\x6b\xfd\xc3\x78\xfc\x86\x11\x41\x49\x1d\x30\x87\x83\xa3\xf1\x0d\xe8\xd6\xa7\xc0\x80\x54\x20\xd9\xba\x19\xcf\xb3\x2c\x8f\x00\xed\xba\x87\xe7\xff\x53\x47\xc3\xa7\xf7\x34\xd9\x8b\xa3\x82\x0f\x29\xbb\x67\xef\x99\x06\xf7\x4c\x12\xb9\x4c\x70\x4e\x71\xfb\x71\x8a\x50\xad\x4b\x67\xbf\x63\xfd\xcc\xaf\xd8\x48\xa6\x5f\xa7\xeb\xdb\xc9\x68\x32\xb1\x76\x33\x2a\x36\xa0\x4a\x43\xdb\xca\x6c\xdb\xe7\xba\xe3\xd7\xbd\x01\xce\xeb\x2e\xe7\xbc\xed\x62\x66\x42\x28\x9d\x35\xe2\x98\x3a\x75\xb7\x3e\x74\xf4\x2b\x14\x71\xb4\x4d\xef\xed\xef\x44\xc2\x4c\xeb\x63\x77\xd9\xa9\xbb\xcc\xb2\xd7\xab\xa7\x47\x7d\x33\x1c\x4e\xf5\xcb\xa5\x18\x8e\xff\xaf\x08\x6e\x79\x30\x11\x17\x12\x2b\x9c\xeb\xd3\xd8\x41\xf7\xfa\x76\xe0\x7d\xa4\xcd\xf6\x14\xee\x5f\x00\x00\x00\xff\xff\xcb\x21\x46\xc5\xe2\x02\x00\x00")

func assetsExeManifestBytes() ([]byte, error) {
	return bindataRead(
		_assetsExeManifest,
		"assets/exe.manifest",
	)
}

func assetsExeManifest() (*asset, error) {
	bytes, err := assetsExeManifestBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/exe.manifest", size: 738, mode: os.FileMode(438), modTime: time.Unix(1518413152, 0)}
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
	"assets/exe.manifest": assetsExeManifest,
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
	"assets": &bintree{nil, map[string]*bintree{
		"exe.manifest": &bintree{assetsExeManifest, map[string]*bintree{}},
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
