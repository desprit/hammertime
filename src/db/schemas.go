// Code generated by go-bindata. DO NOT EDIT.
// sources:
// src/db/schedule/schema.sql (281B)
// src/db/subscription/schema.sql (283B)

package db

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
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

var _srcDbScheduleSchemaSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x8e\x41\x6b\x84\x30\x10\x46\xef\xfe\x8a\xef\xa8\xe0\xa1\x3d\xf7\x14\x65\xda\x06\xd3\xd8\xc6\x08\x7a\x12\xab\x43\x1b\x68\xa5\xa4\xd9\x05\xff\xfd\xe2\xca\xae\xbb\xe0\xf9\xbd\x79\xf3\xe5\x86\x84\x25\x58\x91\x29\x82\x7c\x86\x2e\x2d\xa8\x91\x95\xad\xf0\x3f\x7c\xf3\x78\xf8\x61\xc4\x11\x00\xb8\x11\x52\x5b\x7a\x21\x83\x77\x23\xdf\x84\x69\x51\x50\x9b\x9e\x59\x3f\x04\x77\x74\x61\xee\x6e\xa4\xa5\xa4\x6b\xa5\x50\x6b\xf9\x51\xd3\x2a\x8e\x7d\xe0\xe0\x7e\x19\x96\x1a\x7b\x55\x56\x16\x7c\xef\x26\xf6\x7b\xe8\xd2\xdf\x63\x7f\x9e\x3b\x9e\x82\x9f\x91\x95\xa5\x22\xa1\xb7\xcf\xf9\x2b\xe5\x05\xe2\xcd\x90\x1a\xf1\x43\x8a\xc7\x24\x59\x6f\x3f\xf9\xcb\x4d\xdd\x32\xea\xbe\x1c\x25\x4f\xa7\x00\x00\x00\xff\xff\x54\xf2\xa8\xac\x19\x01\x00\x00")

func srcDbScheduleSchemaSqlBytes() ([]byte, error) {
	return bindataRead(
		_srcDbScheduleSchemaSql,
		"src/db/schedule/schema.sql",
	)
}

func srcDbScheduleSchemaSql() (*asset, error) {
	bytes, err := srcDbScheduleSchemaSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "src/db/schedule/schema.sql", size: 281, mode: os.FileMode(0644), modTime: time.Unix(1699266632, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x50, 0xf0, 0xf6, 0x20, 0x67, 0x50, 0x28, 0x25, 0xc1, 0x7, 0x82, 0x9a, 0x8e, 0x93, 0xba, 0x4a, 0xda, 0x30, 0x8f, 0x2a, 0xbc, 0x54, 0xab, 0xa1, 0x89, 0x97, 0x83, 0x2f, 0xed, 0x3, 0x20, 0x8d}}
	return a, nil
}

var _srcDbSubscriptionSchemaSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8e\xc1\x4a\xc3\x40\x14\x45\xf7\xfd\x8a\xbb\x4c\xa0\x7f\xe0\x6a\x9c\xde\x94\xc1\xf8\xa2\x6f\x5e\xc0\xae\x0a\xa6\x01\x07\xa4\x6a\x26\xf3\xff\x92\x56\x42\x36\xae\xef\xb9\x87\xe3\x95\xce\x08\x73\x8f\x2d\x11\x1a\x48\x67\xe0\x5b\x88\x16\x91\xcb\x7b\x1e\xa6\xf4\x3d\xa7\xaf\x2b\xaa\x1d\x00\xa4\x0b\x82\x18\x8f\x54\xbc\x68\x78\x76\x7a\xc2\x13\x4f\xfb\xdb\x56\xf2\x38\x9d\x37\xc0\x62\x92\xbe\x6d\xef\x6b\x1e\x3e\xc6\x4b\xf9\x1c\xff\x27\x9a\x4e\x19\x8e\xb2\x08\x51\x6d\xf0\x1a\xca\x86\x4a\xf1\x8c\xab\x06\xd5\x32\x74\x82\x03\x5b\x1a\xe1\x5d\xf4\xee\xc0\xbb\xc9\x77\x12\x4d\x5d\x10\x43\xb9\xa6\x9f\x32\x9e\x6f\x6d\xeb\xb7\x97\xf0\xda\x13\xd5\x5f\xf1\x7e\x1b\x57\xef\x6a\xe4\x79\x4a\xc3\xfc\xf0\x1b\x00\x00\xff\xff\x3a\xb9\x42\x58\x1b\x01\x00\x00")

func srcDbSubscriptionSchemaSqlBytes() ([]byte, error) {
	return bindataRead(
		_srcDbSubscriptionSchemaSql,
		"src/db/subscription/schema.sql",
	)
}

func srcDbSubscriptionSchemaSql() (*asset, error) {
	bytes, err := srcDbSubscriptionSchemaSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "src/db/subscription/schema.sql", size: 283, mode: os.FileMode(0644), modTime: time.Unix(1699358564, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xa8, 0xda, 0xee, 0x60, 0x31, 0xfd, 0x44, 0xdd, 0x69, 0x4, 0x67, 0x16, 0x44, 0x3d, 0x63, 0x73, 0x1, 0x51, 0x0, 0x81, 0x6b, 0x28, 0x7e, 0x95, 0xa3, 0xeb, 0xb4, 0xc9, 0x61, 0x9c, 0x7f, 0xe3}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"src/db/schedule/schema.sql":     srcDbScheduleSchemaSql,
	"src/db/subscription/schema.sql": srcDbSubscriptionSchemaSql,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"src": {nil, map[string]*bintree{
		"db": {nil, map[string]*bintree{
			"schedule": {nil, map[string]*bintree{
				"schema.sql": {srcDbScheduleSchemaSql, map[string]*bintree{}},
			}},
			"subscription": {nil, map[string]*bintree{
				"schema.sql": {srcDbSubscriptionSchemaSql, map[string]*bintree{}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
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
	err = os.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
