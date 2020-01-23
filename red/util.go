package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

func run(c string, args ...string) error {
	cmd := exec.Command(c, args...)
	err := cmd.Run()
	log.Printf("executing '%v %v, err: %v'", c, args, err)
	return err
}

func isAZ(s string) bool {
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

func mkdev(major int64, minor int64) uint32 {
	return uint32(unix.Mkdev(uint32(major), uint32(minor)))
}

func mknod(p string, maj int64, min int64) error {
	return syscall.Mknod(p, syscall.S_IFCHR|uint32(os.FileMode(0666)), int(mkdev(maj, min)))
}

func fcopy(dstRoot string, many ...string) error {
	for _, src := range many {
		dst := path.Join(dstRoot, src)
		err := os.MkdirAll(filepath.Dir(dst), 0775)
		if err != nil {
			return err
		}

		in, err := os.Open(src)
		if err != nil {
			return err
		}

		out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
		if err != nil {
			in.Close()
			return err
		}

		_, err = io.Copy(out, in)
		if err != nil {
			in.Close()
			out.Close()
			return err
		}
		err = out.Close()
		if err != nil {
			in.Close()
			return err
		}
		in.Close()
	}
	return nil
}

func uidgid(u string) (int, int, error) {
	x, err := user.Lookup(u)
	if err != nil {
		return 0, 0, err
	}
	uid, err := strconv.ParseInt(x.Uid, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	gid, err := strconv.ParseInt(x.Gid, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(uid), int(gid), nil
}

func chown(u string, dirs ...string) error {
	uid, gid, err := uidgid(u)
	if err != nil {
		return err
	}

	for _, p := range dirs {
		err = os.Chown(p, int(uid), int(gid))
		if err != nil {
			return err
		}
	}

	return nil
}
