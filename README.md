<!-----------------------------

- File Name : README.md

- Purpose :

- Creation Date : 05-17-2017

- Last Modified : Wed 17 May 2017 10:57:04 PM UTC

- Created By : Kiyor

------------------------------->

# getserverlist requirement

getserverlist require install unbound.

Linux please install `yum install unbound`, `apt-get install unbound`

OSX able to install via brew `brew install unbound`

# Usage:

```
# get all server which domain is www.google.com

getserverlist -top 20 -a us -a ge -a hk www.google.com | request -u https://www.google.com -last-modified 'xxxx' -timestamp 'xxxx' -ratio '0.9'

```

- ratio 0.9 means 100 request 90 is expected then return code will be 0, otherwise return code will be 1. So if the data is sensitive, you able to set ratio=1 means must all good.
