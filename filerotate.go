package filerotate

import (
	"os"
	"sync"
)

type Rotate struct {
	FromFile   string
	ToDir      string
	Permission os.FileMode
	ToFile     func() string

	lastToFile string
	fp         *os.File
	mutex      sync.Mutex
}

func (r *Rotate) Write(content []byte) (bytes int, err error) {
	if r.ToFile == nil {
		panic("Rotate.ToFile is nil")
	}

	//1.当Write()第一次被调用时
	//   1.1 创建 fp *os.File
	//   1.2 lastToFile 赋值
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.fp == nil {
		//新建文件
		r.fp, err = os.OpenFile(r.FromFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, r.Permission)
		if err != nil {
			return 0, err
		}
		r.lastToFile = r.ToFile()

		return r.fp.Write(content)
	}

	toFile := r.ToFile()
	if toFile == r.lastToFile {
		return r.fp.Write(content)
	}
	//2. 当Write()遇到需要更换文件时
	//关闭原文件
	if err = r.fp.Close(); err != nil {
		return 0, err
	}
	//重命名文件
	if err = os.Rename(r.FromFile, r.ToDir+"/"+r.lastToFile); err != nil {
		return 0, err
	}
	r.lastToFile = toFile
	//新建文件
	r.fp, err = os.OpenFile(r.FromFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, r.Permission)
	if err != nil {
		return 0, err
	}
	return r.fp.Write(content)

}
