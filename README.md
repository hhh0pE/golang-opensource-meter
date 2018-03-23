Tool to calculate % of OpenSource code used in selected package.

Tested only on windows.
Must work on linux (if not - see and fix line 43. exec.Command syntax)

Usage:
golang-opensource-meter.exe -package=[src_to_package]

For example:
golang-opensource-meter.exe -package=github.com/hhh0pE/golang-opensource-meter

Will display:
StdLib 100%

On other, more big packages normal % is about 5-10%.