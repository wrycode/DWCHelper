# This repository holds research, notes and source code for the project. 

You can leave drive-by comments and messages at
https://todo.sr.ht/~wrycode/olduvai

Click on "tree" to view the files and folders, and "log" to see
messages attached to each change in the project.

Currently, the software component of the project consists of two main
parts: data visualization functions and the "helper" program.

Both components will share several dependencies. Right now, they are 

* [Darwin Core](https://github.com/tdwg/dwc), a standard "intended to
  facilitate the sharing of information about biological diversity"
 
* the database, [PostgreSQL](https://www.postgresql.org/)
 
* Microsoft Access

We will likely contribute upstream to Darwin Core.

As I am now starting to work with Windows (while writing some of the
cross-platform application logic from my Linux laptop), **my current
top focus is setting up the PostgreSQL database server**. Then, in order:

* move/convert the current databases to PostgreSQL, and in the
  process, build the prototype "helper"

* finish the helper (Darwin Core integration, etc.)

* build the data visualization component

## Helper program

The helper program helps format the Microsoft Access database and move
it to the new server. In the process, it will sync with Darwin Core.

The logic will be written in [Common Lisp](https://lisp-lang.org/),
using the [SBCL](http://www.sbcl.org/)
implementation. [Postmodern](http://marijnhaverbeke.nl/postmodern/) is
the library I'm using to interact with PostgreSQL.

The initial conversion will either work with `.accdb` files (using a
third-party dependency) or `.csv` files generated with the built-in
Export Wizard in Access (preferred).

The scope of this part of the project is still under discussion; at a
minimum it will allow a researcher to move their database to the
server after manually changing a few things to fit our data standard.

Additionally, I'd like to:

* Automatically try to detect shared variables with Darwin Core
* Allow the user to confirm or change these variables
* Allow the user to specify aliases (so they can continue to use their
  preferred name for variables, while still conforming to the standard)
* Display the normative definitions of variables from Darwin Core (so
  the user doesn't need to look them up)

One final feature, time-permitting: Allow the user to contribute new
definitions to Darwin Core. This will strengthen the standard. It
would appear as a submission option in the interface. Then it will get
sent to me personally (or someone with the knowledge and authority) to
review and send upstream.

The actual user interface to the helper program hasn't been planned
out. Ideally it will be a simple graphical program with recognizable
forms and buttons, but initially it will probably be an automatically
generated text configuration file.

## Data Visualization

TODO


