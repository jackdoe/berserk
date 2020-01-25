package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var BASIC_CHROOT = []string{
	"/lib64/ld-linux-x86-64.so.2",
	"/lib/terminfo/v/vt102",
	"/lib/terminfo/v/vt52",
	"/lib/terminfo/v/vt220",
	"/lib/terminfo/v/vt100",
	"/lib/terminfo/x/xterm-mono",
	"/lib/terminfo/x/xterm-r5",
	"/lib/terminfo/x/xterm-256color",
	"/lib/terminfo/x/xterm-r6",
	"/lib/terminfo/x/xterm-xfree86",
	"/lib/terminfo/x/xterm",
	"/lib/terminfo/x/xterm-color",
	"/lib/terminfo/x/xterm-vt220",
	"/lib/terminfo/d/dumb",
	"/lib/terminfo/c/cons25-debian",
	"/lib/terminfo/c/cons25",
	"/lib/terminfo/c/cygwin",
	"/lib/terminfo/E/Eterm",
	"/lib/terminfo/a/ansi",
	"/lib/terminfo/h/hurd",
	"/lib/terminfo/r/rxvt",
	"/lib/terminfo/r/rxvt-unicode-256color",
	"/lib/terminfo/r/rxvt-basic",
	"/lib/terminfo/r/rxvt-unicode",
	"/lib/terminfo/m/mach-color",
	"/lib/terminfo/m/mach-bold",
	"/lib/terminfo/m/mach-gnu",
	"/lib/terminfo/m/mach",
	"/lib/terminfo/m/mach-gnu-color",
	"/lib/terminfo/l/linux",
	"/lib/terminfo/s/screen-w",
	"/lib/terminfo/s/screen",
	"/lib/terminfo/s/sun",
	"/lib/terminfo/s/screen.xterm-256color",
	"/lib/terminfo/s/screen-256color",
	"/lib/terminfo/s/screen-256color-bce",
	"/lib/terminfo/s/screen-s",
	"/lib/terminfo/s/screen-bce",
	"/lib/terminfo/p/pcansi",
	"/lib/terminfo/w/wsvt25",
	"/lib/terminfo/w/wsvt25m",
	"/lib/x86_64-linux-gnu/libnss_nis.so.2",
	"/lib/x86_64-linux-gnu/libnss_hesiod.so.2",
	"/lib/x86_64-linux-gnu/libnss_files.so.2",
	"/lib/x86_64-linux-gnu/libnss_nisplus.so.2",
	"/lib/x86_64-linux-gnu/libcom_err.so.2",
	"/lib/x86_64-linux-gnu/libnss_compat.so.2",
	"/lib/x86_64-linux-gnu/libgpg-error.so.0",
	"/lib/x86_64-linux-gnu/libutil.so.1",
	"/lib/x86_64-linux-gnu/libattr.so.1",
	"/lib/x86_64-linux-gnu/libncurses.so.5",
	"/lib/x86_64-linux-gnu/libtinfo.so.5",
	"/lib/x86_64-linux-gnu/libpthread.so.0",
	"/lib/x86_64-linux-gnu/libnss_files-2.27.so",
	"/lib/x86_64-linux-gnu/libnss_nis-2.27.so",
	"/lib/x86_64-linux-gnu/libselinux.so.1",
	"/lib/x86_64-linux-gnu/libm.so.6",
	"/lib/x86_64-linux-gnu/libnss_compat-2.27.so",
	"/lib/x86_64-linux-gnu/libncursesw.so.5",
	"/lib/x86_64-linux-gnu/libdl.so.2",
	"/lib/x86_64-linux-gnu/libnss_dns-2.27.so",
	"/lib/x86_64-linux-gnu/libnss_nisplus-2.27.so",
	"/lib/x86_64-linux-gnu/libacl.so.1",
	"/lib/x86_64-linux-gnu/libnss_dns.so.2",
	"/lib/x86_64-linux-gnu/libexpat.so.1",
	"/lib/x86_64-linux-gnu/libc.so.6",
	"/lib/x86_64-linux-gnu/libpcre.so.3",
	"/lib/x86_64-linux-gnu/libreadline.so.7",
	"/lib/x86_64-linux-gnu/libkeyutils.so.1",
	"/lib/x86_64-linux-gnu/liblzma.so.5",
	"/lib/x86_64-linux-gnu/libbz2.so.1.0",
	"/lib/x86_64-linux-gnu/libnss_hesiod-2.27.so",
	"/lib/x86_64-linux-gnu/libz.so.1",
	"/lib/x86_64-linux-gnu/libnss_systemd.so.2",
	"/lib/x86_64-linux-gnu/libidn.so.11",
	"/lib/x86_64-linux-gnu/libresolv.so.2",
	"/bin/rm",
	"/bin/hostname",
	"/bin/tar",
	"/bin/ls",
	"/bin/bash",
	"/bin/sh",
	"/bin/cat",
	"/bin/gzip",
	"/bin/grep",
	"/bin/mkdir",
	"/bin/mv",
	"/bin/nano",
	"/bin/echo",
	"/bin/more",
	"/bin/uname",
	"/usr/lib/x86_64-linux-gnu/libtasn1.so.6",
	"/usr/lib/x86_64-linux-gnu/libunistring.so.2",
	"/usr/lib/x86_64-linux-gnu/libidn2.so.0",
	"/usr/lib/x86_64-linux-gnu/libffi.so.6",
	"/usr/lib/x86_64-linux-gnu/libpython3.6m.so.1.0",
	"/usr/lib/x86_64-linux-gnu/libsasl2.so.2",
	"/usr/lib/x86_64-linux-gnu/libhogweed.so.4",
	"/usr/lib/x86_64-linux-gnu/libgmp.so.10",
	"/usr/lib/x86_64-linux-gnu/libgpgme.so.11",
	"/usr/lib/x86_64-linux-gnu/libkrb5support.so.0",
	"/usr/lib/x86_64-linux-gnu/libp11-kit.so.0",
	"/usr/lib/x86_64-linux-gnu/libnettle.so.6",
	"/usr/lib/x86_64-linux-gnu/libgnutls.so.30",
	"/usr/lib/x86_64-linux-gnu/libgssapi_krb5.so.2",
	"/usr/lib/x86_64-linux-gnu/libk5crypto.so.3",
	"/usr/lib/x86_64-linux-gnu/libunwind.so.8",
	"/usr/lib/x86_64-linux-gnu/libgpm.so.2",
	"/usr/lib/x86_64-linux-gnu/libassuan.so.0",
	"/usr/lib/x86_64-linux-gnu/libkrb5.so.3",
	"/usr/lib/x86_64-linux-gnu/libunwind-ptrace.so.0",
	"/usr/lib/x86_64-linux-gnu/libtokyocabinet.so.9",
	"/usr/lib/x86_64-linux-gnu/libunwind-x86_64.so.8",
	"/usr/lib/locale/C.UTF-8/LC_MEASUREMENT",
	"/usr/lib/locale/C.UTF-8/LC_NUMERIC",
	"/usr/lib/locale/C.UTF-8/LC_COLLATE",
	"/usr/lib/locale/C.UTF-8/LC_ADDRESS",
	"/usr/lib/locale/C.UTF-8/LC_MESSAGES/SYS_LC_MESSAGES",
	"/usr/lib/locale/C.UTF-8/LC_PAPER",
	"/usr/lib/locale/C.UTF-8/LC_TIME",
	"/usr/lib/locale/C.UTF-8/LC_CTYPE",
	"/usr/lib/locale/C.UTF-8/LC_TELEPHONE",
	"/usr/lib/locale/C.UTF-8/LC_MONETARY",
	"/usr/lib/locale/C.UTF-8/LC_IDENTIFICATION",
	"/usr/lib/locale/C.UTF-8/LC_NAME",
	"/usr/lib/locale/locale-archive",
	"/usr/bin/env",
	"/usr/bin/id",
	"/usr/bin/less",
	"/usr/bin/nnn",
	"/usr/bin/mutt",
	"/usr/bin/touch",
	"/usr/bin/strace",
	"/usr/bin/vim",
	"/usr/bin/tree",
	"/usr/bin/talk",
	"/usr/bin/locale",
	"/usr/bin/clear",
	"/usr/games/nudoku",
	"/etc/resolv.conf",
	"/etc/services",
	"/etc/nsswitch.conf",
	"/etc/hosts",
}

func chroot(u *User) error {
	dev := path.Join(u.Home, "dev")
	err := os.MkdirAll(dev, 0775)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(u.Home, "log"), 0700)
	if err != nil {
		return err
	}

	etc := path.Join(u.Home, "etc")
	err = os.MkdirAll(etc, 0775)
	if err != nil {
		return err
	}

	err = mknod(path.Join(dev, "null"), 1, 3)
	if err != nil {
		return err
	}

	err = mknod(path.Join(dev, "tty"), 5, 0)
	if err != nil {
		return err
	}

	err = mknod(path.Join(dev, "zero"), 1, 5)
	if err != nil {
		return err
	}

	err = mknod(path.Join(dev, "random"), 1, 8)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(etc, "passwd"), []byte(fmt.Sprintf("%s:x:%d:%d:GECOS,,,:/:/bin/bash\n", u.Name, u.Uid, u.Gid)), 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(etc, "profile"), []byte(`
case "$TERM" in
        "dumb")
                export PS1="> "
                ;;
        *)
                export PS1="\033[1;31m\]\\h:\\w\\$\033[00m\] "
                ;;
esac
export HOME=/
export MAIL=/Maildir/

`), 0755)
	if err != nil {
		return err
	}

	err = fcopy(u.Home, BASIC_CHROOT...)
	if err != nil {
		return err
	}

	return nil
}
