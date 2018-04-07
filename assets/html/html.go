// Code generated by go-bindata.
// sources:
// templates/html/inc_feeds.html
// templates/html/inc_foot.html
// templates/html/inc_head.html
// templates/html/inc_items.html
// templates/html/inc_search_form.html
// templates/html/inc_search_results.html
// DO NOT EDIT!

package html

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

var _templatesHtmlInc_feedsHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x50\x50\x50\xb0\x29\xcd\xb1\xe3\x02\x31\xaa\xab\x15\x8a\x12\xf3\xd2\x53\x15\x54\x32\x75\x14\x54\xd2\x14\xac\x6c\x15\xf4\xdc\x52\x53\x53\x8a\x15\x6a\x6b\xc1\x0a\x6c\x72\x32\xed\xaa\xab\x15\x54\xd2\xf4\x42\x32\x4b\x72\x52\x15\x6a\x6b\x6d\xf4\x73\x32\xe1\x9a\x53\xf3\x52\xe0\x2a\xf5\x41\x86\x02\x02\x00\x00\xff\xff\xe3\xb9\x4d\xe7\x5e\x00\x00\x00")

func templatesHtmlInc_feedsHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesHtmlInc_feedsHtml,
		"templates/html/inc_feeds.html",
	)
}

func templatesHtmlInc_feedsHtml() (*asset, error) {
	bytes, err := templatesHtmlInc_feedsHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/html/inc_feeds.html", size: 94, mode: os.FileMode(420), modTime: time.Unix(1523136606, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesHtmlInc_footHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb2\xd1\x4f\xc9\x2c\xb3\xe3\xe2\xb2\x49\xc9\x2c\x53\xc8\x4c\xb1\x55\x4a\xcb\xcf\x2f\x49\x2d\x52\xb2\x83\x4a\xd8\xe8\x27\xe5\xa7\x54\x82\xe8\x8c\x92\xdc\x1c\x3b\x2e\x40\x00\x00\x00\xff\xff\x1a\x9b\x11\x3e\x30\x00\x00\x00")

func templatesHtmlInc_footHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesHtmlInc_footHtml,
		"templates/html/inc_foot.html",
	)
}

func templatesHtmlInc_footHtml() (*asset, error) {
	bytes, err := templatesHtmlInc_footHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/html/inc_foot.html", size: 48, mode: os.FileMode(420), modTime: time.Unix(1523136820, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesHtmlInc_headHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x8e\x31\x8b\xc3\x30\x0c\x85\xf7\xfc\x0a\x9d\xf7\x43\xeb\x0d\x8a\x96\x6b\xe7\x76\xe8\xd2\x51\x8d\x15\x2c\xb0\x1b\x70\x9c\x40\xff\x7d\x49\x9d\x14\x5a\x32\x7d\x82\xf7\x3d\xf1\xe8\xe7\x70\xfa\xbf\x5c\xcf\x47\x08\x25\x45\x6e\xa8\x02\x80\x82\x8a\x5f\x0e\x00\x4a\x5a\x04\xba\x20\x79\xd4\xd2\xba\xa9\xf4\xbf\x7f\x6e\x8d\x8a\x95\xa8\x4c\x58\xb9\xf4\x70\x2b\xd2\x6d\xf0\x0f\x6e\xaa\xe7\x6d\x06\xf3\xad\x5b\x42\xcd\x5b\x7b\x8a\xf5\x00\xa0\x68\x4c\x02\x21\x6b\xdf\x3a\x74\x1c\x86\xa4\x84\xc2\x84\xd1\xf6\x9d\x5e\xd5\x8f\x8e\x5f\x78\x9b\x55\xdc\xf5\x47\x95\xdc\x05\xc7\x95\x9f\xbf\x09\xb7\x25\x84\xde\xe6\xef\xd1\x49\xec\xbe\x4e\x6e\x9e\x01\x00\x00\xff\xff\x9e\x50\x48\x27\x31\x01\x00\x00")

func templatesHtmlInc_headHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesHtmlInc_headHtml,
		"templates/html/inc_head.html",
	)
}

func templatesHtmlInc_headHtml() (*asset, error) {
	bytes, err := templatesHtmlInc_headHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/html/inc_head.html", size: 305, mode: os.FileMode(420), modTime: time.Unix(1523136877, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesHtmlInc_itemsHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x8c\xcd\x4a\xc5\x30\x14\x84\xf7\x7d\x8a\x21\x74\x59\x1a\xe8\x52\xd2\x6c\x14\xa4\xe0\xc2\x85\x3e\x40\xb4\x47\x73\x30\x8d\xa5\x89\xdd\x84\xf3\xee\xd2\x1f\x85\x7b\xef\xf2\x9b\xf9\x66\x00\xa0\x14\x2c\x2e\x7e\x12\x6a\x6e\x50\x73\xc6\x5d\x8f\x76\xc8\x34\x25\x88\x54\x9b\x61\x46\x5e\xf1\x1e\x5c\x4a\xbd\xe2\x4c\x93\xb2\x7b\x0c\x18\xdf\x59\xe3\xe0\x17\xfa\xe8\x55\x29\xdb\xba\x7d\xe2\xf8\x05\x11\x65\x4f\x7e\xe1\x1c\x08\x22\x46\x3b\x0b\x93\x26\x17\x82\xbd\x54\x8d\x3e\x52\xa3\x7d\xf7\xff\x3c\xff\x49\xf7\xdf\x31\x53\xcc\xbb\x37\xdf\xd6\x8f\xaf\xc3\x03\x44\x1a\x9c\xfc\xfc\xf3\x16\x38\x79\x1a\xaf\x06\x7a\xe4\xf5\x80\x52\x40\x71\xab\xab\xdf\x00\x00\x00\xff\xff\x2a\xcf\x9f\x6f\xfe\x00\x00\x00")

func templatesHtmlInc_itemsHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesHtmlInc_itemsHtml,
		"templates/html/inc_items.html",
	)
}

func templatesHtmlInc_itemsHtml() (*asset, error) {
	bytes, err := templatesHtmlInc_itemsHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/html/inc_items.html", size: 254, mode: os.FileMode(420), modTime: time.Unix(1523130288, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesHtmlInc_search_formHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\x8e\x4d\x0a\x83\x30\x14\x84\xf7\x3d\xc5\xf0\x2e\xe0\x05\x12\x77\x3d\x41\x4f\x10\x75\xda\x04\xcc\x8f\xcf\x17\xa8\xb7\x2f\x68\x71\x35\x03\x33\x7c\x7c\x0f\x00\x70\xef\xaa\x79\x3c\x2b\xe0\x52\x69\xdd\x60\x47\xa3\x17\xe3\xd7\x04\x25\x64\x7a\xd9\x04\x69\x39\xa3\xad\x61\x66\xac\xeb\x42\xf5\xf2\x2c\x46\xc5\xd6\xa9\x07\x76\xd3\x54\x3e\x88\x54\x0a\x86\x9b\x38\x75\xb3\x5a\xfe\xc8\xbd\x4f\x39\x99\x8c\x2f\x06\x9d\xa3\x1b\xae\xf1\xfe\x0e\x97\xcb\x2f\x00\x00\xff\xff\x98\x09\x6c\xbe\x98\x00\x00\x00")

func templatesHtmlInc_search_formHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesHtmlInc_search_formHtml,
		"templates/html/inc_search_form.html",
	)
}

func templatesHtmlInc_search_formHtml() (*asset, error) {
	bytes, err := templatesHtmlInc_search_formHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/html/inc_search_form.html", size: 152, mode: os.FileMode(420), modTime: time.Unix(1523133789, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesHtmlInc_search_resultsHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\x50\x50\x50\xb0\xc9\x30\xb4\x0b\x2c\x4d\x2d\xaa\x54\x28\x4a\x2d\x2e\xcd\x29\x29\x56\x48\xcb\x2f\x52\xb0\x29\xb4\xab\xae\x56\xd0\x83\x48\xd4\xd6\xda\xe8\x17\xda\xd9\xe8\x67\x18\xda\x71\x01\x02\x00\x00\xff\xff\x57\xc2\x84\x6b\x33\x00\x00\x00")

func templatesHtmlInc_search_resultsHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesHtmlInc_search_resultsHtml,
		"templates/html/inc_search_results.html",
	)
}

func templatesHtmlInc_search_resultsHtml() (*asset, error) {
	bytes, err := templatesHtmlInc_search_resultsHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/html/inc_search_results.html", size: 51, mode: os.FileMode(420), modTime: time.Unix(1523136388, 0)}
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
	"templates/html/inc_feeds.html": templatesHtmlInc_feedsHtml,
	"templates/html/inc_foot.html": templatesHtmlInc_footHtml,
	"templates/html/inc_head.html": templatesHtmlInc_headHtml,
	"templates/html/inc_items.html": templatesHtmlInc_itemsHtml,
	"templates/html/inc_search_form.html": templatesHtmlInc_search_formHtml,
	"templates/html/inc_search_results.html": templatesHtmlInc_search_resultsHtml,
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
	"templates": &bintree{nil, map[string]*bintree{
		"html": &bintree{nil, map[string]*bintree{
			"inc_feeds.html": &bintree{templatesHtmlInc_feedsHtml, map[string]*bintree{}},
			"inc_foot.html": &bintree{templatesHtmlInc_footHtml, map[string]*bintree{}},
			"inc_head.html": &bintree{templatesHtmlInc_headHtml, map[string]*bintree{}},
			"inc_items.html": &bintree{templatesHtmlInc_itemsHtml, map[string]*bintree{}},
			"inc_search_form.html": &bintree{templatesHtmlInc_search_formHtml, map[string]*bintree{}},
			"inc_search_results.html": &bintree{templatesHtmlInc_search_resultsHtml, map[string]*bintree{}},
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

