Imagine we are going to store data for a web app. Which database performs better: SQLite or MongoDB? I found a few comparisons online, but decided to do my own benchmark anyway.
I have created a simple data structure for test:

type User struct {
	name         string
	passwordHash string
	email        string
	createdAt    dateTime
	level        int
	accessToken  string
}

For SQLite I set PRAGMA journal_mode = WAL; because this mode is supposed to offer a nice balance between performance and reliability.
For MongoDB, I installed the newest version 8.0 and limited WiredTiger cache to 250 MB.

I am going to generate one million rows and commit them one by one in a simple loop.
* SQLite
	* Storage space: 156 MB
	    * Vacuum: has no effect on storage space
	* Time: 33 seconds
	* CPU usage: 5%
* MongoDB
	* Storage space
		* Compressed, default: 112 MB
		* Compressed, zstd: 75 MB
		* Compressed, zlib: 75 MB
		* Uncompressed: 195 MB
	* Time: 2 minutes 14 seconds
	* CPU usage: 8%
		* Client: 3%
		* Server: 5%
* Simple JSON file - for reference, not competing in the benchmark
	* Time: 1 millisecond
	* Storage space: 201 MB

Intermediate conclusion:
* SQLite is 4x faster than MongoDB in single threaded mode
	* As expected, because it has no overhead from client-server communication
* Thanks to compression, MongoDB uses less storage space than SQLite
	* Seems like repeating field names (such as "passwordHash") in every of the one million documents is no big deal
	* zlib compression seems to be especially good: using about 2x less storage space than SQLite
	* Enabling zstd compression seems to have negligible CPU usage, so it is basically free storage saving
* MongoDB CPU usage is 60% higher than SQLite
	* However, if we count only server CPU usage of MongoDB, then SQLite and MongoDB have equal CPU usage

Time for multi-threaded test with 4 threads
* SQLite
	* Time: 37 seconds
	* CPU usage: 6%
* MongoDB
	* Time: 35 seconds
	* CPU usage: 31%
		* Client: 11%
		* Server: 20%

SQLite gained nothing from multi-threaded writing. MongoDB is now as fast as SQLite, but uses 5x more CPU.
