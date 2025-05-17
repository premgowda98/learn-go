# GoRoutines

[Video explaining go-routine](https://youtu.be/S-MaTH8WpOM)
[Video explaining scheduler & go-routine](https://youtu.be/MYtUOOizITs)
[Scheduler](https://youtu.be/dKO827FrzUY)
[Scheduler - Must Watch](https://youtu.be/YHRO5WQGh0k)

Inside a Go runtime there are three things - a number of goroutines, a number of machine threads (which is a logical understanding of how many kernel threads the runtime "owns" at any particular point in time) and procs which is the pipeline between goroutines and machine threads.

When you launch a goroutine it is put into a proc, a queue, if you prefer, where it patiently waits for some CPU time on the machine thread.

The number of procs that exist in your runtime is GOMAXPROCS PLUS any procs holding goroutines that are making what the runtime determines will be blocking calls. That is, when a goroutine makes a (potentially) blocking call, the runtime puts that goroutine and proc aside, so that it makes the blocking call, the machine thread that the proc is associated with is given the blocking call, and, blocks.

Instead of the other goroutines that had been in the same proc as the blocker also blocking and waiting, the runtime requests a new kernel thread, and associates that with the new proc so they can continue working.

The above is better documented in runtime/proc.go

Appropriate setting advice
The heuristic for setting GOMAXPROCS is to set it equal to the number of cores that you have in your system, so that all of the cores can be used by the runtime in parallel.

If you set it to a higher value, then the runtime is going to have to swap procs on and off machine threads so that the procs can get the appropriate CPU TIME (anti-starvation)

If you set it to less than the number of cores you have available, then the runtime will not be requesting time on some of the cores, which matters ONLY if you have more goroutines than procs (because it means that you will have idle goroutines, when you could instead be giving them CPU time


## GOMAXPROCS

1. [Medium Explained](https://medium.com/@aditimishra_541/go-performance-with-ubers-automaxprocs-8f31226a92cd)
2. [Ardan Labs](https://www.ardanlabs.com/blog/2024/02/kubernetes-cpu-limits-go.html)
3. [Ardan Labs - Video](https://youtu.be/Dm7yuoYTx54?list=PLq2Nv-Sh8Eba2gEaId35K2aAUFdpbKx9D)
4. [Imp](https://victoriametrics.com/blog/kubernetes-cpu-go-gomaxprocs/)