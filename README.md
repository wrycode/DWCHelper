# This repository holds research, notes and source code for the project. 

Click on "tree" to view the files and folders, and "log" to see
messages attached to each change in the project. This page contains
the overall plan and status report on the project. At the bottom of
the page are my [notes](#notes).

Currently, the software component of the project consists of two main
parts: data visualization functions and the "helper" program.

Both components will share several dependencies. Right now, they are 

* [Darwin Core](https://github.com/tdwg/dwc), a standard "intended to
  facilitate the sharing of information about biological diversity"
 
* the database, [PostgreSQL](https://www.postgresql.org/) (?)
 
* Microsoft Access

We will likely contribute upstream to Darwin Core.

## Data Visualization

TODO

## Helper program

The helper program imports existing datasets from Microsoft Access
~~and moves them to the server~~. It helps the researcher format their
database using aliases and rules from Darwin Core.

[Here](https://git.sr.ht/~wrycode/olduvai/tree/master/DWCHelper.go) is
the current iteration of the helper program. It is a command-line
program written using Go. [More about this decision](#helper).

The program will:

* Automatically try to detect shared variables with Darwin Core, using
  name matching and a built-in database of aliases
* Allow the user to confirm or change these variables
* Automatically set up aliases (so they can continue to use their
  preferred name for variables, while still conforming to the
  standard). This is accomplished with a separate table in the
  database that stores aliases for the visualizer to use when it
  creates graphs and tables. These aliases will be sent to me, so I
  can review them and add them to the built-in alias database to
  improve variable suggestions
* Display the normative definitions of variables from Darwin Core (so
  the user doesn't need to look them up)
* Allow the user to contribute new definitions to Darwin Core, in the
  case that their variable doesn't match any of the existing ones. This
  will appear as a submission option in the interface. Then it will get
  sent to me personally (or someone with the knowledge and authority) to
  review and send upstream.

## Notes <a id="notes"></a>

### Helper Program notes <a id="helper"></a>

#### Limitations

  Commas are now fine in quoted fields. All special characters work.
  
#### Programming language

Why not lisp? Common Lisp lacks the first-class Windows support I
need, and I felt that I was wasting too much time with its outdated
documentation.

In comparison, [Go](https://golang.org) cross-compiles [almost
effortlessly](https://github.com/golang/go/wiki/WindowsCrossCompiling)
to Windows. It's also much easier to work with.

#### Command line?

Creating a command-line program allows me to focus on the logic of the
application and more quickly move to the data visualization part. The
helper program only needs to display text and take input, and the
command-line is perfect for this. No knowledge of programming or
special syntax will be needed by the researcher. Here's an example of
what the output might look like. I've denoted user input with <>.


     $<helper database.csv>
     Do you want to create a new database on the server? y/n <y>
     Please enter a name for the database: <olduvai-researcher-2>
     We will now run through each variable and allow you to choose
     actions for that variable. Typing "Enter" repeatedly will choose
     the default actions, which are to keep the current variable name
     without linking to a DWC term, and to use the inferred type.
     
     Variable "Catalogue Number"
     Possible DWC equivalent terms: 
     
     ====== catalogNumber(1): An identifier (preferably unique) for the
     record within the data set or collection.
     ====== recordNumber(2): An identifier given to the Occurrence at the
     time it was recorded. Often serves as a link between field notes
     and an Occurrence record, such as a specimen collector's number.
     ====== otherCatalogNumbers(3):A list (concatenated and separated) of
     previous or alternate fully qualified catalog numbers or other
     human-used identifiers for the same Occurrence, whether in the
     current or any other data set or collection.
     ====== fieldNumber(4):An identifier given to the event in the
     field. Often serves as a link between field notes and the Event.
     
     Type "Enter" or any combination of 1, 2, 3, 4: <13>
     Terms catalogNumber and otherCatalogNumber selected.
     Do you want to keep "Catalogue Number" as the variable name? y/n
     <n>
     Please type a new variable name, or "1" or "3" for
     "catalogNumber" or "otherCatalogNumber": <catalogueNumber>
     The variable "Catalogue Number" will be imported as
     "catalogueNumber", with the inferred type "INTEGER" from
     database.csv. It will be linked to the DWC terms "catalogNumber"
     and "otherCatalogNumber". You can change this later before
     uploading your changes. Moving on...
     ...
     
The program will continue through each variable, and then allow the
user to modify records if needed (to change the type, for example,
although type inference should be correct 99% of the time). The
researcher will need a password to the server to upload the database.

### Darwin Core 

[Darwin Core "Simple"](http://rs.tdwg.org/dwc/simple/), a predefined
normative subset of dwc terms that "assumes (and allows) no structure
beyond the concept of rows and columns". DK_E_Fauna is almost entirely
compatable with Darwin Core simple without modifications, because
there are no required fields.

[Guidelines for contributing to dwc terms](https://github.com/tdwg/dwc/blob/master/.github/CONTRIBUTING.md)
[Darwin Core text guide](http://rs.tdwg.org/dwc/text/)

[The csv with all dwc terms and definitions](https://github.com/tdwg/dwc/blob/master/vocabulary/term_versions.csv)

[Simple only terms](https://github.com/tdwg/dwc/blob/master/dist/simple_dwc_vertical.csv)

## PostgreSQL server setup notes <a id="server"></a>

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

