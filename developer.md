# Developer Guide

## Design goals

Automation moves data produced by Titan to the HPC.

## Deployment

### Build

Under the project folder, execute `go build .`

### Execution

Copy the ***tohpc.sh*** file and the newly built ***tohpc*** file to the same directory.

Copy ***config-test.yml*** to the same directory, and rename it to ***config.yml***, modify the contents as needed.

Then execute `./tohpc.sh`,

Then input the password of your private key.

It will run in background, you can view logs in the  ***tohpc.log*** file.

## Command line parameters

- -d, run as daemon, if passed, tohpc will run in background. example: `tohcp -d`
- -pwdfile, use a file stored the private key password, this file will be deleted by tohpc program automatically. example: `tohpc -pwdfile path-to-pwdfile`

If the tohpc program runs in foreground mode and the pwdfile parameter is not specified, the program will ask the user to enter the password of the private key when it starts.

## Configuration instructions

### Source directory

The source directory has three parameters, ***root-dir***, ***overwrite*** and ***start-level***.

#### root-dir

The concept of ***root-dir*** is well understood, and we will automatically move the contents of this directory to the HPC.

#### overwrite

The bool property ***overwrite*** determines whether the program should overwrite existing files, if set to false, the program will rename files instead of overwriting existing ones.

#### start-level

What is the ***start-level***?

This involves our requirements for the directory structure. In order to better realize the automatic archiving of files and the future file search, we hope to have a unified directory structure specifications.

Tat, Tapu and I negotiated a structure as follows:

```
root folder/
├── User folder
│   ├── Project folder
│   │   ├── Dataset folder
│   │   │   ├── frames
│   │   │   └── ...
```

Under the root folder, there're four levels of directories.

Level 1, user folder level. Each user names a directory with their own name.

Level 2, project folder. Each user may have multiple projects, and each directory matches a research project.

Level 3, dataset folder. The data generated each time the microscope is used corresponds to a dataset.

Level4, folders and files inside dataset. Except for the original files produced by the EM should be placed under the ***frames*** directory, other structures are arbitrary.

Then we go back to the ***start-level***. The parameter ***start-level*** means, we only move directories with directory level equal to or greater than ***start-level***, and files within those directories. And the default ***start-level*** is 3, the dataset folder.

Of course, this is not a directory structure design that has been decided and cannot be changed. I hope that by running the ***tohpc*** program, we can inspire everyone to participate in discussions and propose improvement ideas.

This is an example of the folder structure:
```
f4server/
├── Tianming
│   ├── projectname1
│   │   ├── dataset name 1
│   │   │   ├── data
│   │   │   └── frames
│   │   └── dataset name 2
│   │       ├── data
│   │       └── frames
│   └── projectname2
│       ├── dataset name 1
│       └── dataset name 2
├── Tianhui
│   ├── projectname1
│   │   ├── dataset name 1
│   │   │   ├── data
│   │   │   └── frames
│   │   └── dataset name 2
│   │       ├── data
│   │       └── frames
│   └── projectname2
│       ├── dataset name 1
│       │   ├── data
│       │   └── frames
│       └── dataset name 2
│           ├── data
│           └── frames
```

### Dest Folder

The dest folder defines the address, private key, directory and other information of the remote server.

Notice, use a password-encrypted private key to prevent it from being stolen.

### Dustbin

Dustbin defines a trash directory. Files that have been moved to the HPC will not be deleted immediately, but will be moved to the dustbin directory, and the user will delete them after manually checking and confirming that they are correctly transfered.