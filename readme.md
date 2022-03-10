# https://github.com/rubenlab/tohpc

## Developer Guide

If you want to learn how to build, run, and test tohpc programs, or you want to know why the program is designed the way it is, go to the [developer guide](developer.md).

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

**Please notice that only files under the dataset directories(and their sub folders) will be transfered.**

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
