Overview
===

srt is a simple command line app written in Go, the goal is to make the work more efficient for the people who frequently edit subtitles. You can use it to find the error formats within the srt files, such as 
* lint the subtitle number to ensure that it is incremented
* lint the subtitle timeline to ensure that it is incremented
* lint the file to ensure there is no extra space line

Install
===

Once you have [installed Go][golang-install], run the command to install the `srt` tool:
```
    go install github.com/iCell/srt
```
You can also download [the releases](https://github.com/iCell/srt/releases) directly.

Usage
===

Lint a single srt file:
```
    srt lint ~/files/01_file.srt
```
Lint multiple srt files:
```
    srt lint ~/files/01_file.srt ~/files/02_file.srt
```
Lint the files within a directory:
```
    srt lint ~/files/
```
Lint the files within a directory and a file in another path:
```
    srt lint ~/files/ ~/another_files/01_file.srt
```
