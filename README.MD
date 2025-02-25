I have been comparing database performance: MongoDB versus SQLite, and here are the results:

* OrangePi Zero 2W: single board computer with 4 ARM cores and 4 GB of memory
	* SQLite
		* Inserts: around 3000 rows per second
		* Queries: around 5000 rows per second
		* Simultaneous query and update
			* Updates: around 3000 rows per second
			* Queries: around 8000 rows per second
		* CPU load: around 20%
	* MongoDB
		* Inserts: around 1600 per second
		* Queries: around 1400 per second
	* CPU load: around 80%
* Ryzen 5800X3D desktop computer with 8 cores and 16 threads
	* SQLite
		* Storage usage: 165 MB
		* Inserts: 11k per second
			* CPU load: 5%
			* CPU load per 1k operations: 0.455
		* Queries: 87k per second
		* Simultaneous query and update
			* Queries: 63k per second
			* Updates: 28k per second
			* CPU load: 24%
			* CPU load per 1k operations: 0.264
	* MongoDB
		* Storage usage: 162 MB
			* Uncompressed: 250 MB
		* Inserts: 15k per second
			* CPU load: 30%
				* Client: 15%
				* Server: 15%
			* CPU load per 1k operations: 2
		* Queries: 13k per second
		* Simultaneous query and update
			* Queries: 8k per second
			* Updates: 8k per second
			* CPU load: 45%
				* Client: 25%
				* Server: 20%
			* CPU load per 1k operations: 2.81

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
