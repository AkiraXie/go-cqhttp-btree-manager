package base

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/AkiraXie/go-cqhttp-btree-manager/cache"

	"github.com/pkg/errors"
)

type Image struct {
	Md5  string
	Size uint32
	Id   string
	Url  string
}

func (i *Image) String() string {
	return fmt.Sprintf("image md5:%s\nsize:%d\nid:%s\nurl:%s\n", i.Md5, i.Size, i.Id, i.Url)
}

func parseData(data []byte) *Image {
	Md5 := fmt.Sprintf("%x", data[:16])
	Size := binary.BigEndian.Uint32(data[16:20])
	idSize := binary.BigEndian.Uint32(data[20:24])
	Id := string(data[24 : idSize+23])
	Url := string(data[idSize+24:])
	return &Image{Md5, Size, Id, Url}
}

func saveData(data []byte, dst string) (string, error) {
	var dstStr string
	md5 := fmt.Sprintf("%x", data[:16])
	if dst == "" {
		dstStr = md5 + ".image"
	} else {
		dstStr = dst
	}
	return dstStr, errors.Wrap(ioutil.WriteFile(dstStr, data, 0664), "save image failed")
}

func InsertCacheToDb(
	src, db string) error {
	if !filepath.IsAbs(db) {
		db, _ = filepath.Abs(db)
	}
	c := cache.Init(db)
	defer c.Close()
	if !filepath.IsAbs(src) {
		src, _ = filepath.Abs(src)
	}
	fd, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "open cache fail")
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return errors.Wrap(err, "read cache fail")
	}
	c.Insert(data[:16], data)
	return nil
}

func ShowImg(src string) (string, error) {
	if !filepath.IsAbs(src) {
		src, _ = filepath.Abs(src)
	}
	fd, err := os.Open(src)
	if err != nil {
		return "", errors.Wrap(err, "open cache fail")
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return "", errors.Wrap(err, "read cache fail")
	}
	defer func() {
		a := recover()
		err = errors.Errorf("parse cache data error:%v", a)
	}()
	str := parseData(data).String()
	return str, err
}

func ShowImageFromDb(src, db string) (string, error) {
	if !filepath.IsAbs(db) {
		db, _ = filepath.Abs(db)
	}
	c := cache.Init(db)
	defer c.Close()
	md5, err := hex.DecodeString(src)
	if err != nil {
		return "", errors.Wrap(err, "convert string to md5 fail")
	}
	data := c.Get(md5)
	if data == nil {
		return "", errors.Errorf("No data selected from db %s can match md5 %s", db, src)
	}
	defer func() {
		a := recover()
		err = errors.Errorf("parse cache data error:%v", a)
	}()
	str := parseData(data).String()
	return str, err
}

func SaveImageFromDb(src, db, dst string) (string, error) {
	if !filepath.IsAbs(db) {
		db, _ = filepath.Abs(db)
	}
	c := cache.Init(db)
	defer c.Close()
	md5, err := hex.DecodeString(src)
	if err != nil {
		err = errors.Wrap(err, "convert string to md5 fail")
		return "", err
	}
	data := c.Get(md5)
	if data == nil {
		err = errors.Errorf("No data selected from db %s can match md5 %s", db, src)
		return "", err
	}
	defer func() {
		a := recover()
		err = errors.Errorf("parse cache data error:%v", a)
	}()
	str, err := saveData(data, dst)
	return str, err
}

func ShowAllFromDb(db string) (int, string) {
	if !filepath.IsAbs(db) {
		db, _ = filepath.Abs(db)
	}
	c := cache.Init(db)
	defer c.Close()
	var got []string
	var i int
	c.Foreach(func(key [16]byte, value []byte) {
		i++
		s := fmt.Sprintf("Image %d:\n%s--------", i, parseData(value).String())
		got = append(got, s)
	})
	return i, fmt.Sprintf("total Image Number:%v,Detail:\n%s", i, strings.Join(got, "\n"))
}

func DumpAllToDb(db, dst string) int {
	if !filepath.IsAbs(db) {
		db, _ = filepath.Abs(db)
	}
	if !filepath.IsAbs(dst) {
		dst, _ = filepath.Abs(dst)
	}
	c := cache.Init(db)
	defer c.Close()
	d := cache.Init(dst)
	defer d.Close()
	var i int
	c.Foreach(func(key [16]byte, value []byte) {
		i++
		d.Insert(key[:], value)
	})
	return i
}
