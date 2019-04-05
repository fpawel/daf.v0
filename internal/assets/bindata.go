// Code generated by go-bindata.
// sources:
// assets/exe.manifest
// assets/png16/checkmark.png
// assets/png16/error.png
// assets/png16/forward.png
// assets/png16/pin_off.png
// assets/png16/pin_on.png
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

	info := bindataFileInfo{name: "assets/exe.manifest", size: 738, mode: os.FileMode(438), modTime: time.Unix(1554222934, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsPng16CheckmarkPng = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xea\x0c\xf0\x73\xe7\xe5\x92\xe2\x62\x60\x60\xe0\xf5\xf4\x70\x09\x62\x60\x60\x10\x00\x61\x0e\x66\x06\x06\x06\x0d\x5d\xfe\x60\x06\x06\x06\xe6\x62\x27\xcf\x10\x0e\x0e\x8e\xdb\x0f\xfd\x1f\x30\x30\x30\x70\x16\x78\x44\x16\x33\x30\x30\xe4\x83\x30\xe3\xc7\x45\x77\x9c\x19\x18\x18\x24\x4b\x5c\x23\x4a\x82\xf3\xd3\x4a\xca\x13\x8b\x52\x19\xca\xcb\xcb\xf5\x32\xf3\xb2\x8b\x93\x13\x0b\x52\xf5\xf2\x8b\xd2\x67\xbf\xb3\x91\x62\x60\x60\x98\x19\xe0\x13\xe2\xfa\xff\xff\x7f\x86\xff\xff\x19\x66\x9e\x61\x58\x75\x55\x74\xd5\x2b\x91\xb9\x37\x78\x97\x3d\xe3\x59\x74\x87\x63\xe9\x43\xde\x65\xf7\x79\x96\x3d\xe4\x5e\x7a\x9f\x6b\xf9\x03\xae\x25\x8f\x38\x17\xdd\xe7\x59\xf4\x80\x7b\xe9\x43\xee\xc5\x0f\xb9\x96\x3e\xe0\x59\xf2\x90\x07\xc4\x78\xc8\xb3\xf4\x3e\xcf\xe2\x07\xdc\x8b\x1f\x70\x2f\xb9\xcf\x03\x12\x7f\xc0\xbd\xe4\x01\xcf\x62\x10\x9b\x7b\xe9\x03\xae\x25\x20\x95\xdc\x8b\xef\x73\x2f\xb9\xcf\x05\x52\x03\x41\x0f\x61\x0c\x10\xe2\x42\x62\xc3\x51\xae\xfa\xae\x52\x06\x06\x06\xa3\x92\x20\xbf\x60\x06\x46\x56\x36\x1e\x5e\x11\x31\x25\x0d\x6d\x03\x43\x53\x73\x27\x17\x57\xaf\x90\xb0\xb4\xcc\x9c\xa2\xca\xa6\x96\x9e\x99\xb3\xe7\xad\xdd\xb0\xe9\xd0\xe1\xcb\x0f\x1f\x3d\x79\xf3\xf6\xc3\xe7\x2f\x5f\x7f\xfe\xfe\x7b\x39\xe2\x74\x28\x03\x03\x43\x91\xa7\x8b\x63\x88\x84\x64\xed\xc1\x9b\x62\x8d\x01\x0c\x6d\x17\xed\x83\x8f\xb5\xf4\x04\x9e\xb0\x09\xd9\xf7\xbd\xde\x42\xfd\xf4\xa9\xc2\x83\xf7\x7e\xad\xde\x97\x79\xf2\x59\xce\x4e\xbe\x54\xce\xc6\xad\x76\x4c\x09\xf1\xc1\x2b\x27\x1f\xda\xa9\xb9\x78\x61\x53\x75\xc3\x1e\xfd\xe4\xc2\x2c\xef\x55\x2a\x07\x42\x8e\xf1\xe8\xcd\xdb\xc2\xa4\x24\x3a\xed\x40\xad\x01\x63\x59\xe8\xde\x33\x8d\x3a\xed\x8b\x57\x6f\xff\xc8\xac\xfc\x76\x95\xc9\xa9\x06\x07\xb1\xff\x6e\x2b\xf4\xf7\xde\xb8\xa8\xef\xdf\xb6\x8a\xe7\x97\x11\x67\x6b\x04\x03\x03\x03\x83\xa7\xab\x9f\xcb\x3a\xa7\x84\x26\x40\x00\x00\x00\xff\xff\xaf\x82\x80\xc9\xd7\x01\x00\x00")

func assetsPng16CheckmarkPngBytes() ([]byte, error) {
	return bindataRead(
		_assetsPng16CheckmarkPng,
		"assets/png16/checkmark.png",
	)
}

func assetsPng16CheckmarkPng() (*asset, error) {
	bytes, err := assetsPng16CheckmarkPngBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/png16/checkmark.png", size: 471, mode: os.FileMode(438), modTime: time.Unix(1554395435, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsPng16ErrorPng = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x00\x9b\x01\x64\xfe\x89\x50\x4e\x47\x0d\x0a\x1a\x0a\x00\x00\x00\x0d\x49\x48\x44\x52\x00\x00\x00\x10\x00\x00\x00\x10\x08\x06\x00\x00\x00\x1f\xf3\xff\x61\x00\x00\x00\x04\x73\x42\x49\x54\x08\x08\x08\x08\x7c\x08\x64\x88\x00\x00\x00\x09\x70\x48\x59\x73\x00\x00\x0b\x13\x00\x00\x0b\x13\x01\x00\x9a\x9c\x18\x00\x00\x01\x3d\x49\x44\x41\x54\x38\x8d\x9d\x93\xbf\x6a\xc2\x50\x14\xc6\xbf\x73\x70\xa9\xb1\x6f\x20\x16\x87\x42\x08\x0e\xa2\x6b\x71\xc8\x64\x5f\xa3\x83\xae\x7d\x93\x0c\x05\x91\x3c\x49\x32\x28\xa5\x73\x07\x11\xc9\x90\xa1\x68\xc9\x23\xf4\x9a\x52\xe1\x9c\x0e\x4d\x6c\x88\x26\xa5\xfe\xa6\xcb\xbd\xe7\xfb\xce\x77\xff\x11\x4a\x44\xb6\x3d\x50\x91\x29\x03\x2e\x80\x9b\x6c\x7a\xab\xcc\x4b\x61\x9e\xf7\xa2\x68\x55\xac\xa7\x7c\xf0\x3a\x1c\x36\x9b\xc6\x3c\x11\xf0\x50\x36\x2d\x22\x22\xfe\xf5\x7e\xff\xd8\x49\x92\xf4\x68\x90\x89\x43\x02\xee\xea\xc4\x47\x54\x5f\x2c\x63\xc6\x9d\x24\x49\x19\x00\xb2\xce\x27\x62\xcb\xf3\x60\x79\xde\xa9\x01\xd1\xe8\xa3\xd5\xf2\x00\x80\x23\xdb\x1e\xfc\x15\xfb\x1c\x04\x4c\x37\x8e\xd3\x67\x15\x99\xfe\x57\x9c\xc3\x22\x13\xce\x4e\xfb\x22\x14\x70\x19\xbf\x57\x75\x49\x82\x2e\x5f\x2a\xce\x69\x00\xd8\x02\xb8\x3d\xb7\x78\x08\xc3\x5a\xb1\x30\xbf\x35\x94\x79\x49\x15\x06\x5f\x41\x50\x6b\xc0\xaa\x0b\x16\xe6\x79\x55\x41\xe5\x3b\xc8\x20\x55\x9f\x7b\x51\xb4\x12\x11\xff\x5c\xc1\x21\x0c\xab\xb7\xa1\x3a\xb3\xe3\x78\x4d\x00\xf0\xde\x6e\x5f\x19\xcb\x0a\x40\x34\xaa\xcd\x9c\x21\xc0\x73\x33\x4d\xef\xbb\xbb\xdd\x27\x03\x40\x27\x49\x52\xcb\x98\x71\x55\x92\x72\xe7\x5c\x0c\x14\x7e\x63\xce\xc6\x71\xfa\x2c\x32\x51\xc0\x65\x91\x2e\xf0\x73\xda\xac\xba\x20\x55\xdf\x8e\xe3\x75\xb1\xfe\x1b\xf7\xed\x87\xb0\x33\x1d\xd8\x89\x00\x00\x00\x00\x49\x45\x4e\x44\xae\x42\x60\x82\x01\x00\x00\xff\xff\x40\xe1\xb8\x0d\x9b\x01\x00\x00")

func assetsPng16ErrorPngBytes() ([]byte, error) {
	return bindataRead(
		_assetsPng16ErrorPng,
		"assets/png16/error.png",
	)
}

func assetsPng16ErrorPng() (*asset, error) {
	bytes, err := assetsPng16ErrorPngBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/png16/error.png", size: 411, mode: os.FileMode(438), modTime: time.Unix(1554395435, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsPng16ForwardPng = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xea\x0c\xf0\x73\xe7\xe5\x92\xe2\x62\x60\x60\xe0\xf5\xf4\x70\x09\x62\x60\x60\x10\x00\x61\x0e\x66\x06\x06\x06\x0d\x5d\xfe\x60\x06\x06\x06\xe6\x62\x27\xcf\x10\x0e\x0e\x8e\xdb\x0f\xfd\x1f\x30\x30\x30\x70\x16\x78\x44\x16\x33\x30\xb0\x58\x80\x30\xe3\xb3\xa9\x8b\x03\x18\x18\x18\x24\x4b\x5c\x23\x4a\x82\xf3\xd3\x4a\xca\x13\x8b\x52\x19\xca\xcb\xcb\xf5\x32\xf3\xb2\x8b\x93\x13\x0b\x52\xf5\xf2\x8b\xd2\x67\xbf\xb3\x91\x62\x60\x60\xd8\x18\xe0\x13\xe2\xfa\xff\xff\x7f\x86\xff\xff\x1d\xf6\x37\x18\x9f\x99\x29\xb5\x79\xa6\xd2\xee\x99\x0a\xfb\xe7\xeb\x6e\x99\xa6\xbd\x63\x8e\xc6\x9e\x29\x2a\xdb\x66\xab\x6c\x9b\xa9\xb2\x63\xa6\xda\x8e\xd9\xaa\x5b\xe7\x68\x6c\x9f\xad\xbe\x63\xa6\xfa\xce\xd9\x6a\xbb\x66\xa8\xee\x9c\xad\xb6\x63\x96\xda\xae\x59\xaa\x3b\x67\xa9\x81\xd9\x20\xc6\x2e\x98\xc8\x8e\x59\x6a\x50\xf1\x99\x6a\x10\xf1\x9d\x60\xc6\x4e\x30\x63\x07\x58\x0a\x84\x40\x8a\xd5\x77\xc1\xc4\xa1\x0a\x66\xa9\xef\x9c\xa9\x06\x15\x84\x1b\x05\x45\x62\x91\x4e\x2b\x19\x18\x18\xac\x4a\x82\xfc\x82\x19\x18\x59\x58\xb9\xf8\x05\x04\x85\x84\x65\x94\x8d\x6c\xed\x1c\x1c\x9d\x9c\x5d\x53\x8b\x4a\xea\xe7\x2f\x59\xb6\x7c\xe7\xa1\x93\xa7\xcf\x9d\x3f\x7f\xe1\xc2\xc5\x4b\x57\x6f\xdc\xbe\xf3\xf0\xe1\xe3\x27\x4f\x9e\x3e\x7b\xf7\xfe\xe3\xb7\x9f\xbf\xff\xfe\x4b\xcb\x73\x9f\xc3\xc0\xc0\xd0\xe4\xe9\xe2\x18\x22\x21\x19\x7b\xb0\x53\xac\xd1\x81\x81\xf5\x82\xb8\xd6\xa6\xf7\xbb\x7d\x52\xe4\x53\x3c\xef\xfd\xaf\x9f\x91\x69\x55\x74\x6c\xef\x89\xd0\x8f\xe9\xb7\x5b\x35\xf7\x94\x24\x08\xd8\x6c\x9f\xf0\xf9\xd6\x5e\x66\xcd\xd9\x17\xd6\xd4\x49\x1c\xf5\x39\xc9\x36\xe5\xd1\x9f\x93\xfa\x91\xec\x06\x95\x86\xdf\xbf\x35\x3e\x9c\x6b\x7f\xc8\xed\xa9\xd7\x53\x41\xd9\x7c\x15\x66\x69\x09\x95\x25\x32\xd7\xbd\x38\x4e\x4c\x65\x4e\x38\x5c\x11\xd6\xa6\xb4\xa3\x6c\xc1\x9b\xc9\x82\x5b\x62\x5f\x28\x95\x46\x45\x5c\x58\xa6\xb7\x65\x69\xb3\x0f\xa3\x24\xff\x32\xfe\xce\x68\xe1\xdf\x16\x4a\x3f\xe3\x18\x18\x18\x18\x3c\x5d\xfd\x5c\xd6\x39\x25\x34\x01\x02\x00\x00\xff\xff\x11\x80\x79\x9e\x07\x02\x00\x00")

func assetsPng16ForwardPngBytes() ([]byte, error) {
	return bindataRead(
		_assetsPng16ForwardPng,
		"assets/png16/forward.png",
	)
}

func assetsPng16ForwardPng() (*asset, error) {
	bytes, err := assetsPng16ForwardPngBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/png16/forward.png", size: 519, mode: os.FileMode(438), modTime: time.Unix(1534863454, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsPng16Pin_offPng = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xea\x0c\xf0\x73\xe7\xe5\x92\xe2\x62\x60\x60\xe0\xf5\xf4\x70\x09\x62\x60\x60\x10\x00\x61\x0e\x66\x06\x06\x06\x0d\x5d\xfe\x60\x06\x06\x06\x96\x74\x47\x5f\x47\x06\x86\x8d\xfd\xdc\x7f\x12\x59\x19\x18\x18\x14\x92\x3d\x82\x7c\x19\x18\xaa\xd4\x18\x18\x1a\x5a\x18\x18\x7e\x31\x30\x30\x34\xbc\x60\x60\x28\x35\x60\x60\x78\x95\xc0\xc0\x60\x35\x83\x81\x41\xbc\x60\xce\xae\x40\x1b\x06\x06\x86\xc9\x01\x3e\x21\xae\x0c\x0c\x0c\x4e\x07\x7e\x0f\x1e\xf4\xff\xff\xff\x78\x3e\xd6\x0d\x0c\x0c\x0c\xfa\x25\x41\x7e\xc1\x0c\x0c\x6a\x5d\xe7\x3e\x5c\x9a\xa0\xc5\xd2\xf1\xeb\xcf\x14\x8e\x75\xdf\x0c\xd8\x34\x9a\x3e\xed\x6a\x7b\xa4\xa3\x72\x6b\x86\x9e\x52\x4f\xdd\xb4\x63\x26\x37\x5e\xf0\x7c\xe1\x3a\x61\xd4\x57\xd3\xb0\xa7\xf0\x67\xd6\x32\x06\x06\x06\xc6\x24\x6f\x77\x17\x83\x75\x77\x74\x9f\x30\x30\x30\xb0\x97\x78\xfa\xba\xb2\x3f\xe4\x94\xe2\x56\xe7\xf8\x6d\x9c\xb4\x9b\x81\x81\x61\x8a\xa7\x8b\x63\x88\xc4\xe5\xd0\xfe\x6e\xbe\x26\x03\x1e\x57\xed\xe4\x07\x72\x4a\x1d\xc2\x81\x92\x0f\x7e\xdf\xfd\xff\xdf\xfa\xb5\x5b\xc4\x0e\xcf\xd9\x6f\x96\x6c\xdf\xcd\xaa\xe0\xaa\xf3\xcc\xc6\xf3\xf9\x03\xd5\x25\x4e\xea\x7b\x4f\xdc\xa9\xee\xb4\xd4\x3e\xb2\x2a\x46\xfa\xce\x1f\xcf\xd6\xdb\xb7\xfa\xdf\x7d\xab\x58\xb0\xee\x53\xda\x87\xa6\xa7\xbb\x16\x36\x8a\xb9\x95\xde\x38\x2c\xe8\xcf\x18\x61\x73\x98\x75\x3d\xe3\x99\x2b\xbd\x07\x6a\xb9\xe4\x7e\xcc\x89\xd2\xaa\x61\xfd\xb0\x6c\x69\xd9\xaf\x5c\x46\x85\xe3\x09\xd7\x58\x8b\x2f\xa6\xea\xba\xfc\x34\xfc\xf1\xf5\xb8\x81\x95\x5f\xf6\x2a\x3b\xe9\x3d\x81\xd2\xc5\x93\x97\x9c\xfb\xc9\x2c\x22\xc0\xdb\x2c\x23\xc4\xb6\x98\x81\x81\x41\xb5\xc4\x35\xa2\x24\x25\xb1\x24\xd5\x2a\xb9\x28\x35\xb1\x24\x95\xc1\xc8\xc0\xd0\x5c\xd7\xc0\x52\xd7\xc8\x2c\xc4\xd0\xd0\xca\xd8\xd2\xca\xc0\x42\xd7\xc0\xc4\xca\xc0\x80\xf3\x8b\x57\x0a\x8a\x86\xdc\xfc\x94\xcc\xb4\x4a\xdc\x1a\x2a\x56\x7e\xba\xc1\xc0\xc0\xc0\xe0\xe9\xea\xe7\xb2\xce\x29\xa1\x09\x10\x00\x00\xff\xff\xf4\x2e\xa6\x6b\x65\x02\x00\x00")

func assetsPng16Pin_offPngBytes() ([]byte, error) {
	return bindataRead(
		_assetsPng16Pin_offPng,
		"assets/png16/pin_off.png",
	)
}

func assetsPng16Pin_offPng() (*asset, error) {
	bytes, err := assetsPng16Pin_offPngBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/png16/pin_off.png", size: 613, mode: os.FileMode(438), modTime: time.Unix(1534863454, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsPng16Pin_onPng = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xea\x0c\xf0\x73\xe7\xe5\x92\xe2\x62\x60\x60\xe0\xf5\xf4\x70\x09\x62\x60\x60\x10\x00\x61\x0e\x66\x06\x06\x06\x0d\x5d\xfe\x60\x06\x06\x06\x96\x74\x47\x5f\x47\x06\x86\x8d\xfd\xdc\x7f\x12\x59\x19\x18\x18\x14\x92\x3d\x82\x7c\x19\x18\xaa\xd4\x18\x18\x1a\x5a\x18\x18\x7e\x31\x30\x30\x34\xbc\x60\x60\x28\x35\x60\x60\x78\x95\xc0\xc0\x60\x35\x83\x81\x41\xbc\x60\xce\xae\x40\x1b\x06\x06\x86\xb4\x00\x9f\x10\x57\x06\x06\x06\xa7\x03\xbf\x69\x8a\xfe\xff\xff\xaf\xd7\x55\x7d\x01\xe4\xba\x92\x20\xbf\x60\x06\xb5\xae\x73\x1f\x2e\x4d\xd0\x62\xe9\xf8\xf5\x67\x0a\xc7\xba\x5d\x6d\x33\x54\xf4\xa6\x1d\xbb\xf1\xe2\xcb\x09\x25\x83\x3d\x3d\x9f\x74\xce\x4e\x99\xa3\xc5\xc0\xc0\xc0\x98\xe4\xed\xee\xa2\x78\x24\x87\x57\x8c\x81\x81\x81\xbd\xc4\xd3\xd7\x95\xfd\x21\xa7\x14\xb7\x0a\x9f\x65\xed\x14\x5f\x06\x06\x86\x54\x4f\x17\xc7\x10\x89\xcb\xa9\xfd\xd7\x85\x1a\x14\x78\x58\xfc\x8f\xde\x13\xad\x3e\xfd\xbf\x7e\x29\xd3\x2b\x89\x9b\xea\x9b\x98\x1f\xf2\xed\x28\x88\xf1\xe7\xb0\xe7\x48\xec\x78\xd4\x62\x92\x39\xed\xf5\xce\x05\xd3\x83\x66\x76\xd9\x5d\x71\x98\x5c\xe3\x7a\xeb\x13\xb7\x8b\x90\xac\xb8\x2d\x9b\x9e\x5b\xb4\xcf\x7f\xf5\x6b\x97\x66\xd6\xef\x66\xf8\xa8\x2f\xd1\xbd\x6f\x96\xd1\xbb\x60\xbf\x67\xe5\xf1\xe9\x9f\x8f\xbc\x63\x6f\xb7\x9a\x9f\xfb\x49\x7c\xef\x5e\x9b\x0f\x99\x0c\x0c\x0c\xaa\x25\xae\x11\x25\x29\x89\x25\xa9\x56\xc9\x45\xa9\x89\x25\xa9\x0c\x46\x06\x86\xe6\xba\x06\x96\xba\x46\x66\x21\x86\x86\x56\xc6\x66\x56\x86\x26\xba\x06\x26\x56\x06\x06\x9f\xbf\x56\x33\xa3\x68\xc8\xcd\x4f\xc9\x4c\xab\xc4\xad\xa1\x69\xc5\xe1\xfd\x0c\x0c\x0c\x0c\x9e\xae\x7e\x2e\xeb\x9c\x12\x9a\x00\x01\x00\x00\xff\xff\xf8\x8c\xaa\x35\xfa\x01\x00\x00")

func assetsPng16Pin_onPngBytes() ([]byte, error) {
	return bindataRead(
		_assetsPng16Pin_onPng,
		"assets/png16/pin_on.png",
	)
}

func assetsPng16Pin_onPng() (*asset, error) {
	bytes, err := assetsPng16Pin_onPngBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/png16/pin_on.png", size: 506, mode: os.FileMode(438), modTime: time.Unix(1534863454, 0)}
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
	"assets/exe.manifest":        assetsExeManifest,
	"assets/png16/checkmark.png": assetsPng16CheckmarkPng,
	"assets/png16/error.png":     assetsPng16ErrorPng,
	"assets/png16/forward.png":   assetsPng16ForwardPng,
	"assets/png16/pin_off.png":   assetsPng16Pin_offPng,
	"assets/png16/pin_on.png":    assetsPng16Pin_onPng,
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
	"assets": {nil, map[string]*bintree{
		"exe.manifest": {assetsExeManifest, map[string]*bintree{}},
		"png16": {nil, map[string]*bintree{
			"checkmark.png": {assetsPng16CheckmarkPng, map[string]*bintree{}},
			"error.png":     {assetsPng16ErrorPng, map[string]*bintree{}},
			"forward.png":   {assetsPng16ForwardPng, map[string]*bintree{}},
			"pin_off.png":   {assetsPng16Pin_offPng, map[string]*bintree{}},
			"pin_on.png":    {assetsPng16Pin_onPng, map[string]*bintree{}},
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
