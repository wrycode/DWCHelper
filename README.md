# Description

DWCHelper is a command-line utility to help format and clean up CSV
files (for instance, exported from Microsoft Access). It:

- formats the file according to RFC 4180 (cleans up extra quotes,
  etc.)
- detects and suggests aliases to [Darwin Core](https://dwc.tdwg.org/)
  terms
- detects and suggests terms that may not be used, and can be removed
- allows the user to rename or remove terms
- saves the conversion settings for future runs (to accommodate
  changes to the dataset)
  
# Installation
Windows: An installer for the latest release can be found on the
[Releases](https://github.com/wrycode/DWCHelper/releases) page. 

Linux: You can use the binary provided on the releases page, or easily
build from source with the following steps:

- set your GOPATH
- install the only dependancy: `go get -u github.com/fatih/camelcase`
- clone the repo and run run `go build`

# Usage
Navigate to the location of your CSV dataset in the console and run: 
`DWCHelper <input-filename.csv> <output-filename.csv>`

For Windows users, this means you need to navigate to the folder
containing your CSV file in Windows Explorer, then click in the
navigation bar and type "cmd" (and press Enter). The black command
prompt window that opens up is where you type `DWCHelper
<input-filename.csv> <output-filename.csv>`.

On the first run for each dataset, DWCHelper will prompt you for
various corrections to the data. It will save your choices in the
`.settings` file (in Windows Explorer, it appears as `<filename>.txt`
with the type SETTINGS, but is still a normal text file that you can
open with Notepad) for subsequent runs; if you want to redo the
prompts, simply delete this file.

### Editing `.settings`
The `.settings` file can be edited with a text editor to avoid redoing
the prompts for small changes. DWCHelper is fairly tolerant of errors
in this file and will simply ignore typos and terms that aren't
in your dataset. 

The first line is a CSV list of terms to remove completely from the
dataset during the conversion. 

Any lines after that are term aliases. The first value on each line is
the term to be renamed and the second value is the new name.

# About

DWCHelper is one component of my 2019 Undergraduate Research and
Creativity Award project, which is a collaborative effort with the
Anthropology department at UNCG.

The eventual goal of the project is to provide a tool for researchers
at different sites in Olduvai Gorge, Tanzania to easily share,
compare, and combine datasets and create useful, publishable data
visualizations.

In June of 2019, I will be traveling to Tanzania to excavate and
analyze animal bones, and I hope to gain a broader understanding of
the context surrounding these 1.4 to 2 million-year-old specimens. My
objective is to understand what types of questions researchers may
need answered in their quest to understand this period of human
evolution. 

# Todo

- comment/clean helper functions, tidy up everything
- add better testing/examples
- Continuous Integration and publish releases on Sourcehut instead
- remove stringIsVariation or tighten it up to cut down on false
positives (waiting for more sample aliases before doing this)
- refactor functions with a simplified "database" struct (no need for
the separate "column" struct; I overbuilt that part of the program in
anticipation of a more complicated standard)
