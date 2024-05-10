# Ringpool
Ringpool is a golang implementation of ring buffer for fixed length element. It can be used in scenario such as packet buffer in network protocol implementation.
Ringpool provides debug tool to trace call stack and channel. For example, you can call AddCallStack at the beginning of a function, and call PopCallStack at the end of the function 
If you enable debug, the Ringpool can also periodically check if there is any unreturned element for longer than a specified period of time, and print out its call stack or which channel it stays in.
