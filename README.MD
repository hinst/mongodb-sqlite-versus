I have been benchmarking SQLite versus MongoDB as my side project. Here are the results:

* SQLite is much faster than MongoDB: around 6x...10x times faster
    * The only favorable test benchmark for MongoDB is concurrent insertions
* SQLite uses less CPU power than MongoDB: around 7x...8x more efficient
* MongoDB uses less storage space than SQLite, thanks to zstd compression: -65% storage usage
	* Enabling zstd compression has a negligible impact on CPU usage
	* With compression disabled, SQLite wins in storage usage
* MongoDB scales better with more CPU power and concurrent connections: +10%...+20% more scaling
	* Insertions are special case in SQLite: they are limited to one thread and do not scale at all

|                      | SQLite            | MongoDB     |
| -------------------- | ----------------- | ----------- |
| Speed                | ✅ 6x...10x faster |             |
| CPU usage            | ✅ 7x...8x less    |             |
| Storage uncompressed | ✅ -50% less       |             |
| Storage usage        |                   | ✅ -65% less |
| Thread scaling       | 1.69x             | ✅ 1.81x     |
| Insertions scaling   | no                | ✅ yes       |


TEST CONDITIONS
* Parallel threads: 4
* Batch size: 100 rows
* SQLite journal mode: WAL
* MongoDB compression: zstd

TEST RESULTS
* Ryzen 5800X3D desktop computer with 8 cores and 16 threads
	* SQLite
		* Storage usage: 165 MB
		* Inserts: 10k per second
			* CPU load: 6%
			* CPU load per 1k operations: 0.6
		* Queries: 159k per second
			* CPU load: 22%
			* CPU load per 1k operations: 0.138
		* Simultaneous query and update
			* Queries: 52k per second
			* Updates: 31k per second
			* CPU load: 16%
			* CPU load per 1k operations: 0.193
	* MongoDB, version 7
		* Storage usage: 100 MB
			* Uncompressed: 250 MB
		* Inserts: 15k per second
			* CPU load: 30%
				* Client: 15%
				* Server: 15%
			* CPU load per 1k operations: 2.00
		* Queries: 13k per second
			* CPU load: 30%
				* Client: 15%
				* Server: 15%
			* CPU load per 1k operations: 2.31
		* Simultaneous query and update
			* Queries: 7k per second
			* Updates: 6k per second
			* CPU load: 30%
				* Client: 15%
				* Server: 15%
			* CPU load per 1k operations: 2.31
* OrangePi Zero 2W: single board computer with 4 ARM cores and 4 GB of memory
	* SQLite
		* Inserts: 1.8k per second
			* CPU load: 20%
			* CPU load per 1k operations: 11.1
		* Queries: 23k per second
			* CPU load: 80%
			* CPU load per 1k operations: 3.5
		* Simultaneous query and update
			* Queries: 11k per second
			* Updates: 4k per second
			* CPU load: 50%
			* CPU load per 1k operations: 3.33
	* MongoDB, version 6
		* Inserts: 1.6k per second
			* CPU load: 80%
			* CPU load per 1k operations: 50
		* Queries: 1.5k per second
			* CPU load: 80%
			* CPU load per 1k operations: 53.3
		* Simultaneous query and update
			* Queries: 0.72k per second
			* Updates: 0.72k per second
			* CPU load: 80%
			* CPU load per 1k operations: 55.6
	* CPU load: around 80%

Data structure used for test:
	type User struct {
		name         string
		passwordHash string
		email        string
		createdAt    dateTime
		level        int
		accessToken  string
	}

THREAD SCALING
* Ryzen 5800X3D desktop
	* SQLite: 3 threads -> 6 threads
			* insertions: +4% throughput
			* queries: +44% throughput
			* combined queries: +87% throughput
			* combined updates: +77% throughput
		* Average +53% throughput
			* +69% excluding insertions
	* MongoDB: 3 threads -> 6 threads
			* insertions: +44% throughput
			* queries: +44% throughput
			* combined queries: +98% throughput
			* combined updates: +101% throughput
		* Average +72% throughput
			* +81% excluding insertions

Full source code: https://github.com/hinst/mongodb-sqlite-versus