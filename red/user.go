package main

import (
	"bytes"
	"fmt"
	"net/mail"
	"os"
	"os/user"
	"path"
	"sync/atomic"
	"time"
)

type User struct {
	Name string
	Uid  int
	Gid  int
	Home string
}

func NewUser(u string) (*User, error) {
	err := userIsValid(u)
	if err != nil {
		return nil, err
	}

	uid, gid, err := uidgid(u)
	if err != nil {
		return nil, err
	}

	if uid < 1000 {
		return nil, fmt.Errorf("invalid user")
	}

	return &User{Name: u, Home: path.Join(ROOT, u), Uid: uid, Gid: gid}, nil
}

func (u *User) LogP(logname string, data []byte) {
	err := u.Log(logname, data)
	if err != nil {
		panic(err)
	}

}
func (u *User) Log(logname string, data []byte) error {
	filename := path.Join(ROOT, u.Name, "log", logname)
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

func (u *User) SetAuthorizedKeyFile(key []byte) error {
	authorizedKeyFile := path.Join("etc", "ssh", "authorized_keys", u.Name)

	err := appendAuthorizedKey(authorizedKeyFile, key)
	if err != nil {
		return err
	}

	err = chown(u.Uid, u.Gid, authorizedKeyFile)
	if err != nil {
		return err
	}

	return nil
}

var counter = uint64(0)

func (u *User) Mail(m []byte) error {
	_, err := mail.ReadMessage(bytes.NewBuffer(m))
	if err != nil {
		return err
	}

	basename := fmt.Sprintf("%v.M%vP%v_%v.%v", time.Now().Unix(), time.Now().Nanosecond()/1000, os.Getpid(), atomic.AddUint64(&counter, 1), "berserk")

	fn := path.Join(u.Home, "Maildir", "tmp", basename)
	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	size, err := f.Write(m)
	if err != nil {
		f.Close()
		_ = os.Remove(fn)

		return err
	}
	f.Close()

	err = chown(u.Uid, u.Gid, fn)
	if err != nil {
		_ = os.Remove(fn)

		return err
	}

	newname := path.Join(u.Home, "Maildir", "new", fmt.Sprintf("%v,S=%v", basename, size))
	err = os.Rename(fn, newname)
	if err != nil {
		_ = os.Remove(fn)
		return err
	}
	return nil
}

var ALLOWED = map[string]os.FileMode{
	"tmp":                       0700,
	"Mail":                      0700,
	"private":                   0700,
	"Maildir":                   0700,
	path.Join("Maildir", "cur"): 0700,
	path.Join("Maildir", "new"): 0700,
	path.Join("Maildir", "tmp"): 0700,
	"public_html":               0700,
}

func (u *User) Enable() error {
	for dir, perm := range ALLOWED {
		p := path.Join(u.Home, dir)
		err := os.MkdirAll(p, perm)
		if err != nil {
			return err
		}

		// FIXME: recursive
		err = chown(u.Uid, u.Gid, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) Disable() error {
	for dir := range ALLOWED {
		p := path.Join(u.Home, dir)

		// FIXME: recursive
		err := chown(0, 0, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func userIsValid(u string) error {
	if len(u) < 1 || len(u) > 8 || !isAZ(u) || BLACKLIST[u] {
		return fmt.Errorf("user is blacklisted, only lte 8 characters and a-z usernames are allowed")
	}
	return nil
}

func CreateSystemUser(username string, key []byte) error {
	err := keyIsValid(key)
	if err != nil {
		return err
	}

	err = userIsValid(username)
	if err != nil {
		return err
	}

	_, err = user.Lookup(username)
	if err == nil {
		return fmt.Errorf("already exists")
	}

	if _, ok := err.(user.UnknownUserError); !ok && err != nil {
		return err
	}

	err = run("/usr/sbin/adduser", "--firstuid", "1000", "--gecos", "GECOS", "--home", path.Join(ROOT, username), "--no-create-home", "--disabled-password", "--add_extra_groups", username)
	if err != nil {
		return err
	}

	err = run("/usr/sbin/setquota", username, "1G", "1G", "10000", "10000", ROOT)
	if err != nil {
		return err
	}

	u, err := NewUser(username)
	if err != nil {
		return err
	}

	err = u.SetAuthorizedKeyFile(key)
	if err != nil {
		return err
	}

	return chroot(u)
}
