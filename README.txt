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
        │   ├── locale
        │   ├── nnn
        │   ├── talk
        │   ├── touch
        │   ├── tree
        │   ├── vim

    * there is no outgoing internet from the machine
    * users are chrooted to their homedir

-b
