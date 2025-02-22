Imagine we are going to store data for a web app. Which database performs better: SQLite or MongoDB? I found a few comparisons online, but decided to do my own benchmark anyway.
I have created a simple data structure for test:

type User struct {
	name         string
	passwordHash string
	email        string
	createdAt    time.Time
	level        int
}

Let us begin with SQLite and see how it performs. I set PRAGMA journal_mode = WAL; because this mode is supposed to offer a nice balance between performance and reliability.
* An empty database with NOTHING in it occupies 4 kilobytes on disk.
* An empty database with one table in it and 0 rows occupies 12 kilobytes on disk.

Now I am going to generate 100k rows and see what happens.