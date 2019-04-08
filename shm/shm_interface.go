package shm

import (
	"os"
	"syscall"
	"unsafe"
)

const shm_Len = 128

type ShmChannel struct {
	pFile    *os.File
	shm_Name string
	shm_WBuf []byte
}

func New() ShmChannel {
	return ShmChannel{}
}

func (this *ShmChannel) Create(name string) error {
	// 创建
	var err error
	Unlink(name)
	this.pFile, err = Open(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	// 容量
	if err := syscall.Ftruncate(int(this.pFile.Fd()), int64(1024)); err != nil {
		return err
	}
	// 获取([]byte)
	fd := int(uintptr(this.pFile.Fd()))
	this.shm_WBuf, err = syscall.Mmap(fd, 0, 1024, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return err
	}

	return nil
}

func (this *ShmChannel) Set(byt []byte) bool {
	if this.shm_WBuf != nil && len(byt) < shm_Len {
		buf := (*[shm_Len]byte)(unsafe.Pointer(&this.shm_WBuf[0]))[:]
		for i := 0; i < shm_Len; i++ {
			buf[i] = 0
		}
		for i := 0; i < len(byt); i++ {
			buf[i] = byt[i]
		}
		return true
	}
	return false
}

func (this *ShmChannel) Get(name string) []byte {
	var byt []byte
	file, err := Open(name, os.O_RDONLY, 0600)
	if err != nil {
		return nil
	}
	defer file.Close()
	// 获取
	fd := int(uintptr(this.pFile.Fd()))
	buf, err := syscall.Mmap(fd, 0, 1024, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil
	}
	if buf != nil {
		org := (*[shm_Len]byte)(unsafe.Pointer(&buf[0]))[:]
		copy(byt, org)
		return byt
	}
	return nil
}

func (this *ShmChannel) Close() {
	if this.pFile != nil {
		this.pFile.Close()
		Unlink(this.pFile.Name())
	}
}
