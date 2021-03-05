package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/unix"
)

func BPF(cmd int, attr unsafe.Pointer, size uintptr) (uintptr, error) {
	r1, _, errNo := unix.Syscall(unix.SYS_BPF, uintptr(cmd), uintptr(attr), size)
	runtime.KeepAlive(attr)

	var err error
	if errNo != 0 {
		err = errNo
	}

	return r1, err
}

type bpfGetFDByIDAttr struct {
	id    uint32
	next  uint32
	flags uint32
}

func bpfGetFDByID(id uint32) (uint32, error) {
	attr := bpfGetFDByIDAttr{
		id: id,
	}
	fd, err := BPF(unix.BPF_PROG_GET_FD_BY_ID, unsafe.Pointer(&attr), unsafe.Sizeof(attr))
	return uint32(fd), err
}

type bpfObjName [unix.BPF_OBJ_NAME_LEN]byte

type bpfProgramInformation struct {
	prog_type                uint32
	id                       uint32
	tag                      [unix.BPF_TAG_SIZE]byte
	jited_prog_len           uint32
	xlated_prog_len          uint32
	jited_prog_insns         unsafe.Pointer
	xlated_prog_insns        unsafe.Pointer
	load_time                uint64 // since 4.15 cb4d2b3f03d8
	created_by_uid           uint32
	nr_map_ids               uint32
	map_ids                  unsafe.Pointer
	name                     bpfObjName // since 4.15 067cae47771c
	ifindex                  uint32
	gpl_compatible           uint32
	netns_dev                uint64
	netns_ino                uint64
	nr_jited_ksyms           uint32
	nr_jited_func_lens       uint32
	jited_ksyms              unsafe.Pointer
	jited_func_lens          unsafe.Pointer
	btf_id                   uint32
	func_info_rec_size       uint32
	func_info                unsafe.Pointer
	nr_func_info             uint32
	nr_line_info             uint32
	line_info                unsafe.Pointer
	jited_line_info          unsafe.Pointer
	nr_jited_line_info       uint32
	line_info_rec_size       uint32
	jited_line_info_rec_size uint32
	nr_prog_tags             uint32
	prog_tags                unsafe.Pointer
	run_time_ns              uint64
	run_cnt                  uint64
}

type bpfObjGetInfoByFDAttr struct {
	fd      uint32
	infoLen uint32
	info    unsafe.Pointer
}

func bpfGetInfoByFD(fd uint32) (*bpfProgramInformation, error) {
	var info bpfProgramInformation
	attr := bpfObjGetInfoByFDAttr{
		fd:      fd,
		infoLen: uint32(unsafe.Sizeof(info)),
		info:    unsafe.Pointer(&info),
	}
	_, err := BPF(unix.BPF_OBJ_GET_INFO_BY_FD, unsafe.Pointer(&attr), unsafe.Sizeof(attr))
	if err != nil {
		return nil, fmt.Errorf("fd %v: %w", fd, err)
	}
	return &info, nil
}
