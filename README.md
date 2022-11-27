# K-Means Custering via MapReduce Implementation 

<!--
## How Run
- Execute the bash script **build.sh** for building *.go files.  
```bash
  $ ./build.sh
```
- Run **master server**
```bash
  $ ./master
  Master online on port [9001]
```

- Run an arbitrary number of mappers, specifying the port number where each of them will listen

```bash
$ ./mapper 1234
$ ./mapper 4321 # on different shell
$ ./mapper 1212 # on different shell
```
- Run at least one reducer (by default only one reducer will be used as it is not needed for our implementation)
```bash
$ ./reducer 7777
```

![](assets/example_1.png)

- Run client for grep a specific word in one of the files in the master folder
```bash
$ ./clientGrep {word} {file}
```
![](assets/example_2.png)

-->