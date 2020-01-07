# hugo
A Go implementation of [BlueFever's chess engine](https://www.youtube.com/watch?v=bGAfaepBco4&list=PLZ1QII7yudbc-Ky058TEaOstZHVbT-2hg) 
with additional modifications and improvements such as a primitive king safety evaluation. 

# Todo:
* Switch to incremental evaluation
* Connect to polyglot & syzygy
* Add dynamic movetime calculation based on position complexity (i.e. tactical position should be analyzed longer than quiet position)
* Add multithreading
