package main

const SLASH = `
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
`

const THANKS_FOR_PAYING = `
Hi,

Thanks for paying.

In few seconds paypal will send IPN notification to
https://berserk.red and the account will be enabled.

for more info use:
    curl https://berserk.red

-b
`

const AFTER_REGISTER = `
Hi,
Thanks for registering at https://berserk.red.

The price is 1€ per month for 1GB of space. 
(trial 0.1€ for the first month)

you can pay to it with paypal by following the redirect at:

    https://berserk.red/sub/%s

After you pay you can access it by:

    ssh %s@berserk.red "echo hi > public_html/index.html"

Please send feedback to: jack@baxx.dev

Thanks again.

-b
`

const LICENSE = `
By purchasing, registering for, or using the “Berserk” services (the
  “Services”) you ("referred in the document as "you", "customer",
  "subscriber", "client") enter into a contract with Prymr B.V. KvK:
  68884842, Overtoom 141, 1054 HG Amsterdam, The Netherlands (also
  referred in the document as "berserk", "we"),and you accept and agree
  to the following terms (the “Contract”). The Contract shall apply to
  the supply of the Services, use of the Services after purchase or
  after registering for limited free use where such offer has been
  made available.

Services:
  We provide the service of storing your data, with specific retention
  rate, depending on the the offer you have registered the scope of
  this service might vary.

LIMITATION OF LIABILITY:
  THE CUSTOMER IS LIABLE FOR THE CONTENT ITSELF TO THE FULLEST EXTENT
  PERMITTED BY LAW, BERSERK SHALL NOT BE LIABLE FOR ANY INDIRECT,
  INCIDENTAL, SPECIAL, CONSEQUENTIAL OR PUNITIVE DAMAGES OR LOST
  PROFITS, OR LOST REVENUE ARISING OUT OF THIS AGREEMENT, INCLUDING
  WITHOUT LIMITATION: (1) THE USE OF OR INABILITY TO USE THE SERVICE,
  (2) LOSS OR ALTERATION OF CONTENT, (3) ANY CLAIM ATTRIBUTABLE TO
  ERRORS, OMISSIONS OR OTHER INACCURACIES IN THE SERVICE, (4)
  UNAUTHORIZED ACCESS TO OR ALTERATION OF CONTENT OR OTHER
  TRANSMISSIONS, OR (5) ANY OTHER MATTER RELATING TO THE SERVICE, EVEN
  IF BERSERK HAS BEEN ADVISED OF THE POSSIBILITY OF SUCH DAMAGES. TO THE
  EXTENT PERMITTED BY LAW, YOU AGREE THAT BERSERK’S TOTAL LIABILITY FOR
  DAMAGES RELATED TO THE SERVICE IS LIMITED TO THE TOTAL AMOUNT YOU
  HAVE PAID FOR THE SERVICE OVER THE 12 MONTH PERIOD LEADING UP TO THE
  CAUSE OF THE CLAIM, OR, IF YOUR CLAIM AROSE DURING A FREE TRIAL
  PERIOD, TO THE THEN-CURRENT ANNUAL AMOUNT CHARGED FOR THE SERVICE.

Acceptable Conduct:
  You are responsible for the actions of all users of your account and
  any data that is created, stored, displayed by, or transmitted by
  your account while using Berserk. You will not engage in any activity
  that interferes with or disrupts Berserk's services or networks
  connected to Berserk.

Contract Duration
  You agree that any malicious activities are considered prohibited
  usage and will result in immediate account suspension or
  cancellation without a refund and the possibility that we will
  impose fees; and/or pursue civil remedies without providing advance
  notice.

  You agree that Berserk shall be permitted to charge your credit card on
  a 30 days in advance of providing services. Payment is due
  every 30 days. Once the subscription is cancelled the service is stopped.

  The contract can be cancelled from Berserk at any time.

Services and Data
  Subscriber is solely responsible for the preservation of
  Subscriber's data which Subscriber saves onto its Berserk account (the
  “Data”). Even with respect to Data as to which Subscriber contracts
  for data services provided by Berserk, We shall have no
  responsibility to preserve Data. We shall have no liability for any
  Data that may be lost or inaccessible.

Use of the service is at your own risk.

THE SERVICE IS PROVIDED "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL I BE LIABLE FOR ANY DIRECT, INDIRECT,
INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR
TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE
USE OF THIS SERVICE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH
DAMAGE.

## GDPR

We are not sharing the data with anyone for no purposes what so ever.
We are keeping logs of IP adress registering, logging in, and the
paypal payment notifications for starting/ending the subscription.
`
