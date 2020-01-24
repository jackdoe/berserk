package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"

	"golang.org/x/crypto/ssh"
)

var BLACKLIST = map[string]bool{
	"berserk":             true,
	"admin":               true,
	"register":            true,
	"toor":                true,
	"root":                true,
	"daemon":              true,
	"bin":                 true,
	"sys":                 true,
	"sync":                true,
	"games":               true,
	"man":                 true,
	"lp":                  true,
	"mail":                true,
	"news":                true,
	"uucp":                true,
	"proxy":               true,
	"www-data":            true,
	"backup":              true,
	"list":                true,
	"irc":                 true,
	"gnats":               true,
	"nobody":              true,
	"systemd-timesync":    true,
	"systemd-network":     true,
	"systemd-resolve":     true,
	"syslog":              true,
	"_apt":                true,
	"messagebus":          true,
	"uuidd":               true,
	"avahi-autoipd":       true,
	"usbmux":              true,
	"dnsmasq":             true,
	"rtkit":               true,
	"cups-pk-helper":      true,
	"speech-dispatcher":   true,
	"whoopsie":            true,
	"geoclue":             true,
	"kernoops":            true,
	"saned":               true,
	"pulse":               true,
	"nm-openvpn":          true,
	"avahi":               true,
	"colord":              true,
	"hplip":               true,
	"gnome-initial-setup": true,
	"gdm":              true,
	"systemd-coredump": true,
	"sshd":             true,
	"postgres":         true,
}

func chroot(u string, p string) error {
	dev := path.Join(p, "dev")
	err := os.MkdirAll(dev, 0775)
	if err != nil {
		return err
	}

	etc := path.Join(p, "etc")
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

	uid, gid, err := uidgid(u)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(etc, "passwd"), []byte(fmt.Sprintf("%s:x:%d:%d:GECOS,,,:/:/bin/bash\n", u, uid, gid)), 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(etc, "profile"), []byte(`
PS1="\033[1;31m\]\\h:\\w\\$\033[00m\] "
HOME=/
MAIL=/Maildir/

`), 0755)
	if err != nil {
		return err
	}

	err = fcopy(p, BASIC_CHROOT...)
	if err != nil {
		return err
	}

	return nil
}
func keyIsValid(key []byte) error {
	_, _, _, _, err := ssh.ParseAuthorizedKey(key)
	return err
}

func appendAuthorizedKey(p string, key []byte) error {
	err := keyIsValid(key)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(key); err != nil {
		return err
	}

	if _, err = f.WriteString("\n"); err != nil {
		return err
	}
	return nil
}

func userIsValid(u string) error {
	if len(u) < 1 || len(u) > 8 || !isAZ(u) || BLACKLIST[u] {
		return fmt.Errorf("user is blacklisted, only <=8 a-z usernames are allowed")
	}
	return nil
}

func addUser(u string, key []byte) error {
	err := keyIsValid(key)
	if err != nil {
		return err
	}

	err = userIsValid(u)
	if err != nil {
		return err
	}

	_, err = user.Lookup(u)
	if err == nil {
		return fmt.Errorf("already exists")
	}

	if _, ok := err.(user.UnknownUserError); !ok && err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(ROOT), 0775)
	if err != nil {
		return err
	}

	home := path.Join(ROOT, u)

	err = os.MkdirAll(path.Join(home, "public_html"), 0755)
	if err != nil {
		return err
	}

	for _, dir := range []string{"tmp", "Mail", "private", "Maildir", path.Join("Maildir", "cur"), path.Join("Maildir", "new"), path.Join("Maildir", "tmp"), "log"} {
		err = os.MkdirAll(path.Join(home, dir), 0700)
		if err != nil {
			return err
		}
	}

	err = run("/usr/sbin/adduser", "--firstuid", "1000", "--gecos", "GECOS", "--home", home, "--no-create-home", "--disabled-password", "--add_extra_groups", u)
	if err != nil {
		return err
	}

	err = run("/usr/sbin/setquota", u, "1G", "1G", "10000", "10000", ROOT)
	if err != nil {
		return err
	}
	authorizedKeyFile := path.Join("etc", "ssh", "authorized_keys", u)

	err = appendAuthorizedKey(authorizedKeyFile, key)
	if err != nil {
		return err
	}

	err = chown(u, authorizedKeyFile)
	if err != nil {
		return err
	}

	return chroot(u, home)
}

func appendUserLog(u string, logname string, data []byte) error {
	filename := path.Join(ROOT, u, "log", logname)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(fmt.Sprintf("------------------ %v\n", time.Now())))
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte("\n------------------\n"))
	if err != nil {
		return err
	}

	return f.Close()
}
