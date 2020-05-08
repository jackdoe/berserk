
experimental, use it at your own risk

also people have local shell on the machine, despite automatic
security updates, it is safe to assume the machine hacked.


------------------

Hi,

https://berserk.red is a shell + web hosting service

The price is 1€ per month for 1GB of space.
(trial 0.1€ for the first month)

you can register by sending your pub key to:

    cat ~/.ssh/id_rsa.pub | curl -d@- https://berserk.red/register/:username
    # usernames are lowercase a-z up to 8 characters

    # By registering *YOU AGREE* with the Terms Of Service, visible at:
    curl https://berserk.red/tos

then follow the instructions.

After you register and pay you can access it by:

    ssh username@berserk.red "echo hi > public_html/index.html"

    sftp username@berserk.red

    sshfs .. etc

All the logs about your payment/registration are stored in log/ and
are not accessible by anyone but you and me.

Please send feedback to:
    jack@baxx.dev or https://github.com/jackdoe/berserk

Usage:
    public_html/
       everything under this directory can be accessed via web on
       https://berserk.red/~username/

    private/
       your private files

    available commands: (at the moment)

    ├── bin
    │   ├── bash
    │   ├── cat
    │   ├── echo
    │   ├── grep
    │   ├── gzip
    │   ├── ls
    │   ├── mkdir
    │   ├── more
    │   ├── mv
    │   ├── nano
    │   └── tar
    ├── dev
    │   ├── null
    │   ├── random
    │   ├── tty
    │   └── zero
    └── usr
        ├── bin
        │   ├── clear
        │   ├── id
        │   ├── less
        │   ├── mutt
        │   ├── locale
        │   ├── nnn
        │   ├── talk
        │   ├── touch
        │   ├── tree
        │   ├── vim

    * there is no outgoing internet from the machine
    * users are chrooted to their homedir

-b


------------------

there is no docker, no cloud, no replication, no nothing

using 1 machine with attached volume on digital ocean (maybe will move
it to hetzner)


/etc/security/limits.conf:
        memory, nprocs, cpu, etc..

quota:
        usrquota,grpquota

chroot:
        Match Group berserk
                AuthorizedKeysFile /etc/ssh/authorized_keys/%u
                ChrootDirectory %h

        this is questionable, I feel a bit safer, but it is more
        annoying to let people interact, esp not having /proc mounted

security:
        up-to-date ubuntu, with daily security updates etc.. but that
        only gets us so far, assume the machine hacked.

        /home/user is owner by root (and user is chrooted into it)
        /home/user/private is owned by user:user and mode is 0700
        /home/user/public_html is owned by user:user
        /home/user/log owned by root mode 700 and files are 0600

        log contains logs of the http request registering the user
        and payment subscription events from paypal.

        secuirty updates are automatically installed and if they
        require restart the machine is automatially restarted
        (immidiately after the update)

internet:
        there is only 443,80 and dne outgoing and 443,80 and ssh incoming

talk:
        you can talk with someone by typing
        $ talk jack@127.0.0.1
        it is pretty cool

backups:
        there are no backups


games:

   nudoku - sudoku

   +---+---+---+---+---+---+---+---+---+     nudoku 0.2.5
   | 7 | 3 |   | 2 | 4 |   |   | 8 | 9 |     level: easy
   +---+---+---+---+---+---+---+---+---+
   | 8 |   | 2 |   | 5 | 7 | 3 | 6 | 4 |     Commands
   +---+---+---+---+---+---+---+---+---+      Q - Quit
   | 4 |   | 6 | 3 |   |   |   |   | 7 |      r - Redraw
   +---+---+---+---+---+---+---+---+---+      h - Move left
   | 1 |   | 3 | 4 |   |   | 9 |   | 8 |      l - Move right
   +---+---+---+---+---+---+---+---+---+      j - Move down
   |   | 4 | 8 | 6 |   | 9 | 2 | 1 |   |      k - Move up
   +---+---+---+---+---+---+---+---+---+      x - Delete number
   | 2 |   |   | 8 | 1 |   |   | 4 | 5 |      c - Check solution
   +---+---+---+---+---+---+---+---+---+      N - New puzzle
   | 9 | 8 |   | 5 |   |   |   | 3 | 1 |      S - Solve puzzle
   +---+---+---+---+---+---+---+---+---+      H - Give a hint
   | 6 | 1 | 5 | 7 |   | 4 |   |   | 2 |
   +---+---+---+---+---+---+---+---+---+
   |   | 2 | 7 | 9 | 8 | 1 |   | 5 | 6 |
   +---+---+---+---+---+---+---+---+---+


NB: use it at your own risk


TODO:

* mud
* ircd
* bbs

* dont chroot?

------------------

Why?

> why not?


I feel I have lost my edge, havent been reading the exploit lists, and
just chilling as if everything is taken care of.  Having a box where
people can get a shell will surely keep me on my toes.

Shell + public_html was my first experience on the web, in ~1999, I
had an ugly website where I used to write poems for my girlfriend (now
my wife). Sadly when they shut it down I did not download it and now
they are lost forever (I dont remember what I wrote, but I doubt it is
a big loss). Anyway, it was a lot of fun trying to quit vi for the
first time haha.
