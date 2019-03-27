
# Table of Contents

1.  [Digital Ocean Debian 9 setup (olduvai)](#org7a1987c)
    1.  [Debian 9.7, PostgreSQL 9.6](#org96a8308)
    2.  [from ssh into root:](#org08f12b7)
        1.  [adduser wrycode](#org95b436d)
        2.  [usermod -aG sudo wrycode](#org4f624ef)
        3.  [ufw firewall](#org91bb08e)
        4.  [cp -r ~/.ssh  /home/wrycode](#org5dea4ea)
        5.  [chown -R wrycode:wrycode *home/wrycode*.ssh/](#org4804288)
        6.  [apt-install emacs git man-db mosh postgresql-9.6 postgresql ptop](#orgbf10128)
        7.  [su - postgres](#org1a3a2a9)
        8.  [from this postgres user you can add 'roles' (essentially postgresql users)](#org1205d1a)
        9.  [createdb olduvai](#org8b4ad54)
        10. [psql](#org6aadcf9)
        11. [edit /etc/postgresql/9.6/main/postgresql.conf](#orge7c6966)
        12. [edit /etc/postgresql/9.6/main/pg<sub>hba.conf</sub>:](#orgce792ac)
        13. [systemctl restart postgresql](#org24e4ba8)
    3.  [test the connection from local machine:](#org394ea50)
    4.  [add alias locally in ~/.ssh/config](#orga38c282)
    5.  [scp .tmux.conf and .bashrc to the host](#org3c9c00d)
    6.  [from wrycode user on server:](#org27ae4be)
        1.  [em .bashrc and modify the host-specific settings](#org5b90f91)


<a id="org7a1987c"></a>

# Digital Ocean Debian 9 setup (olduvai)


<a id="org96a8308"></a>

## Debian 9.7, PostgreSQL 9.6


<a id="org08f12b7"></a>

## from ssh into root:


<a id="org95b436d"></a>

### adduser wrycode


<a id="org4f624ef"></a>

### usermod -aG sudo wrycode


<a id="org91bb08e"></a>

### ufw firewall

1.  apt install ufw

2.  ufw app list

3.  ufw allow OpenSSH

4.  ufw allow 60000:61000/udp (mosh)

5.  ufw allow 873 (rsync)

6.  ufw allow 5432 (postgresql)

7.  ufw enable


<a id="org5dea4ea"></a>

### cp -r ~/.ssh  /home/wrycode


<a id="org4804288"></a>

### chown -R wrycode:wrycode *home/wrycode*.ssh/


<a id="orgbf10128"></a>

### apt-install emacs git man-db mosh postgresql-9.6 postgresql ptop


<a id="org1a3a2a9"></a>

### su - postgres


<a id="org1205d1a"></a>

### from this postgres user you can add 'roles' (essentially postgresql users)

info about roles and users:
<https://www.digitalocean.com/community/tutorials/how-to-use-roles-and-manage-grant-permissions-in-postgresql-on-a-vps--2>


<a id="org8b4ad54"></a>

### createdb olduvai


<a id="org6aadcf9"></a>

### psql

1.  create role olduvai;

2.  alter role wrycode with nosuperuser;

3.  \password wrycode


<a id="orge7c6966"></a>

### edit /etc/postgresql/9.6/main/postgresql.conf

uncomment the line with "localhost" in it, and replace localhost with "\*"


<a id="orgce792ac"></a>

### edit /etc/postgresql/9.6/main/pg<sub>hba.conf</sub>:

host olduvai wrycode 0.0.0.0/0 md5


<a id="org24e4ba8"></a>

### systemctl restart postgresql


<a id="org394ea50"></a>

## test the connection from local machine:

psql -h ip -d olduvai -U wrycode


<a id="orga38c282"></a>

## add alias locally in ~/.ssh/config


<a id="org3c9c00d"></a>

## scp .tmux.conf and .bashrc to the host


<a id="org27ae4be"></a>

## from wrycode user on server:


<a id="org5b90f91"></a>

### em .bashrc and modify the host-specific settings

