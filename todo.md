## Add timeout settings

When a large file is being transferred, if the network is interrupted, the program will be stuck transferring the file for a long time.

## Add connection interruption handling

If the connection is interrupted during the file transfer process, the program will skip to the next file to continue the transfer, but at this time, the program should try to reconnect before processing the next file.

## Improve the way tohpc.sh runs

Internally, tohpc.sh uses the method of writing the password entered by the user to a file to pass the password to the tohpc program. A better way is to pass the password through a linux pipe.