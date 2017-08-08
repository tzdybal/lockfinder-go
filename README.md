# lockfinder-go
lockfinder reads stack trace and finds all goroutines holding/waiting on locks.

This simple application reads a stack trace of gouroutines, and for every element of every goroutine stack it opens corresponding source file, and reads it line by line (in backwards direction) to check if corresponding function has any lock.

Summary of all acquired locks is displayed - for every goroutine that holds lock, stack trace is displayed, with extra information about lines that holds some locks.

Currently app is very simple:
 * it looks for "Lock" to decide if line of code acquire lock.
 * every item on every stack trace is processed separately - this may be not effective
 * application is single threaded
 * there is no way to filter output, which may be useful to exclude locks in external libraies
