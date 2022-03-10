# https://github.com/rubenlab/tohpc

## Design goals

Automation moves data produced by titan to the HPC.

## Deployment

### Build

execute `go build .`

### Execution

Copy the ***tohpc.sh*** file and the newly built ***tohpc*** file to the same directory.

Copy ***config-test.yml*** to the same directory, and rename it to ***config.yml***, modify the contents as needed.

Then execute `./tohpc.sh`,

Then input the password of your private key.

It will run in background, you can view logs in the  ***tohpc.log*** file.

## Configuration instructions

### Source directory

The source directory has two parameters, ***root-dir*** and ***start-level***.

The concept of ***root-dir*** is well understood, and we will automatically move the contents of this directory to the HPC.

But what is the ***start-level***?

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

## Usage

The ***tohpc*** program is deployed on the pre-processing server. It monitors the ***tohpc*** directory on the ***Titan offload drive***.

The destination folder on the HPC is ***/scratch1/projects/rubsak/from-offload***.

So for users, all you have to do is to create a new directory with your name in the ***tohpc*** directory on the Titan offload driver, and create a new project folder in this directory. Cut the datasets you need to your project folders.

Notice, use cut not copy. Since the program automatically removes files, copy will report an error, and there is no reason to use copy.

After the files are moved to the HPC, they will not be deleted but moved to the ***tohpcDustbin*** folder.

So you can check if the files are transfered correctly, and then delete them safely to free up space.

## Required directory structure:

```
tohpc/
├── User folder
│   ├── Project folder
│   │   ├── Dataset folder
│   │   │   ├── frames
│   │   │   └── ...
```

Under the ***tohpc*** folder, there're four levels of directories.

Level 1, user folder level. Each user names a directory with their own name.

Level 2, project folder. Each user may have multiple projects, and each directory matches a research project.

Level 3, dataset folder. The data generated each time the microscope is used corresponds to a dataset.

Level4, folders and files inside dataset. Except for the original files produced by the EM should be placed under the ***frames*** directory, other structures are arbitrary.

**Please notice that only files in under the dataset directories(and their sub folders) will be transfered.**

For example:

```
tohpc/
├── Tianming
│   ├── demo project
│   │   ├── projectFile.txt
│   │   ├── dataset1
│   │   │   ├── frames
│   │   │   ├── dataFile.txt
│   │   │   └── ...
```

**The whole folder "*dataset1*" will be transfered to the HPC. But the file "*projectFile.txt*" will be ignored.**