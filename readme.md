# https://github.com/rubenlab/tohpc

## Developer Guide

If you want to learn how to build, run, and test tohpc programs, or you want to know why the program is designed the way it is, go to the [developer guide](developer.md).

## Usage

The ***tohpc*** program is deployed on the pre-processing server. It monitors the ***tohpc*** directory on the ***Titan offload drive***.

The destination folder on the HPC is ***/scratch1/projects/rubsak/from-offload***.

So for users, all you have to do is to drop your files in the ***tohpc*** directory on the Titan offload driver. The program will transfer these files to the ***/scratch1/projects/rubsak/from-offload*** directory on the HPC, and move the transfered files to the ***tohpcDustbin*** directory on the offload driver.

The program will not delete files, so you can check if the files are transfered correctly, and then delete them safely to free up space.

Notice, use cut not copy. Since the program automatically removes files, copy will report an error, and there is no reason to use copy.
