# Tee

A basic implementation of Unix / Linux tee command

## Build
```
$ make build
```

## Alias
In order to make manual testing easy create an alias
```shell
$ alias tee='./bin/tee'
```

## Usage
```shell
$ ls -al | tee --append a.txt b.txt c.txt

$ tee --append a.txt
Only in a.txt!
Only in a.txt!
```

## Verify
```shell
$ cat a.txt
total 52
drwxrwxr-x 5 adrian adrian 4096 Nov 29 01:37 .
drwxrwxr-x 3 adrian adrian 4096 Nov 28 15:26 ..
-rw-r--r-- 1 adrian adrian  579 Nov 29 01:37 a.txt
drwxrwxr-x 2 adrian adrian 4096 Nov 29 01:26 bin
(Not showing for clarity)
Only in a.txt!

$ cat b.txt
total 52
drwxrwxr-x 5 adrian adrian 4096 Nov 29 01:37 .
drwxrwxr-x 3 adrian adrian 4096 Nov 28 15:26 ..
-rw-r--r-- 1 adrian adrian  579 Nov 29 01:37 a.txt
drwxrwxr-x 2 adrian adrian 4096 Nov 29 01:26 bin
(Not showing for clarity)

$ cat c.txt
total 52
drwxrwxr-x 5 adrian adrian 4096 Nov 29 01:37 .
drwxrwxr-x 3 adrian adrian 4096 Nov 28 15:26 ..
-rw-r--r-- 1 adrian adrian  579 Nov 29 01:37 a.txt
drwxrwxr-x 2 adrian adrian 4096 Nov 29 01:26 bin
(Not showing for clarity)
```

## Unalias
Restore system original tee command
```shell
$ unalias tee
```

## Test
```shell
$ make test
```
