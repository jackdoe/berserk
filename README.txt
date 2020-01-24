
experimental, use it at your own risk


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

using 1 machine with attached volume on digital ocean:


/etc/security/limits.conf:
        memory, nprocs, cpu, etc..

quota:
        usrquota,grpquota

chroot:
        Match Group berserk
                AuthorizedKeysFile /etc/ssh/authorized_keys/%u
                ChrootDirectory %h

security:
        up-to-date ubuntu, with daily security updates etc.. but that
        only gets us so far, assume the machine hacked.

        /home/user is owner by root (and user is chrooted into it)
        /home/user/private is owned by user:user and mode is 0700
        /home/user/public_html is owned by user:user
        /home/user/log owned by root mode 700 and files are 0600

        log contains logs of the http request registering the user
        and payment subscription events from paypal.

internet:
        there is no outgoing internet, only http, https, ssh input and
        those established connections are allowed

talk:
        you can talk with someone by typing
        $ talk jack@127.0.0.1
        it is pretty cool

NB: use it at your own risk


TODO:
* mud
* ircd
* bbs

------------------

Why?

> why not?


------------------

INBOX:

I am experimenting with some way to receive messages (without mail)

example:

cat <<EOF |  curl -XPOST --data-binary @- https://berserk.red/mail/jack
From: John Doe <jdoe@machine.example>
To: Mary Smith <mary@example.net>
Subject: Saying Hello
Date: Fri, 21 Nov 1997 09:55:06 -0600
Message-ID: <1234@local.machine.example>

This is a message just to say hello.
So, "Hello".
EOF



this will create file in ~jack/Maildir/new with the mail inside


