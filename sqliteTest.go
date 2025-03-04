package main

import (
	"database/sql"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

const DB_FILE_PATH = "./test-sqlite.db"
const LIBSQL_URL = "http://localhost:8088"
const DB_TIMEOUT = 60 * 60 * 1000
const SQLITE_TEST_MODE_FILE = 0
const SQLITE_TEST_MODE_HTTP = 1

type SqliteTest struct {
	users       []*User
	batchSize   int
	threadCount int
	mode        int
}

func (me *SqliteTest) prepare() {
	if checkFileExists(DB_FILE_PATH) {
		assertError(os.Remove(DB_FILE_PATH))
	}
	if checkFileExists(DB_FILE_PATH + "-shm") {
		assertError(os.Remove(DB_FILE_PATH + "-shm"))
	}
	if checkFileExists(DB_FILE_PATH + "-wal") {
		assertError(os.Remove(DB_FILE_PATH + "-wal"))
	}
	var db = me.open()
	defer me.close(db)
	var setupText = readStringFromFile(executablePath + "/setup.sql")
	assertResultError(db.Exec(setupText))
}

func (me *SqliteTest) open() *sql.DB {
	switch me.mode {
	case SQLITE_TEST_MODE_FILE:
		return assertResultError(sql.Open("sqlite3", "file:"+DB_FILE_PATH+
			"?_journal_mode=WAL&_busy_timeout="+strconv.Itoa(DB_TIMEOUT)))
	case SQLITE_TEST_MODE_HTTP:
		var db = assertResultError(sql.Open("libsql", LIBSQL_URL))
		db.Exec("PRAGMA journal_mode=WAL;")
		db.Exec("PRAGMA busy_timeout=" + strconv.Itoa(DB_TIMEOUT) + ";")
		return db
	default:
		panic("Unknown mode: " + strconv.Itoa(me.mode))
	}
}

func (me *SqliteTest) close(db *sql.DB) *sql.DB {
	if db != nil {
		assertError(db.Close())
	}
	return nil
}

func (me *SqliteTest) run() {
	me.prepare()

	var insertDuration = me.runInserts()
	var insertionsPerSecond = float64(len(me.users)) / insertDuration.Seconds()
	fmt.Printf(TAB+"insertion duration: %v, rows per second: %v\n",
		insertDuration, humanize.CommafWithDigits(insertionsPerSecond, 0))

	var readDuration = me.runQueriesLoop(me.threadCount)
	var readsPerSecond = float64(len(me.users)) / readDuration.Seconds()
	fmt.Printf(TAB+"reading duration: %v, rows per second: %v\n",
		readDuration, humanize.CommafWithDigits(readsPerSecond, 0))

	var combinedReadDuration, combinedUpdateDuration = me.runCombined()
	var combinedReadsPerSecond = float64(len(me.users)) / combinedReadDuration.Seconds()
	var combinedUpdatesPerSecond = float64(len(me.users)) / combinedUpdateDuration.Seconds()
	fmt.Printf(TAB+"combined read & update benchmark: %v reads per second, %v updates per second\n",
		humanize.CommafWithDigits(combinedReadsPerSecond, 0),
		humanize.CommafWithDigits(combinedUpdatesPerSecond, 0))
	fmt.Printf(TAB+TAB+"read duration %v, update duration %v\n",
		combinedReadDuration, combinedUpdateDuration)

	var beginning = time.Now()
	var sizeBefore, sizeAfter = me.compress()
	var compressionDuration = time.Since(beginning)

	fmt.Printf("SQLite file size: %v -> %v, compression duration: %v\n",
		formatFileSize(sizeBefore), formatFileSize(sizeAfter), compressionDuration)
}

func (me *SqliteTest) runInserts() time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for range me.threadCount {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			me.writeUsers(usersChannel)
		}()
	}
	for _, user := range me.users {
		usersChannel <- user
	}
	close(usersChannel)
	waitGroup.Wait()
	var elapsed = time.Since(beginning)

	return elapsed
}

// Reading in SQLite is too fast when nobody is writing
func (me *SqliteTest) runQueriesLoop(threadCount int) time.Duration {
	const count = 10
	var totalDuration time.Duration
	for range count {
		totalDuration += me.runQueries(threadCount)
	}
	return totalDuration / count
}

func (me *SqliteTest) runQueries(threadCount int) time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for range threadCount {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			me.readUsers(usersChannel)
		}()
	}
	for _, user := range me.users {
		usersChannel <- user
	}
	close(usersChannel)
	waitGroup.Wait()
	var elapsed = time.Since(beginning)

	return elapsed
}

func (me *SqliteTest) runUpdates(threadCount int) time.Duration {
	var usersChannel = make(chan *User)
	var waitGroup sync.WaitGroup

	var beginning = time.Now()
	for range threadCount {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			me.updateUsers(usersChannel)
		}()
	}
	for _, user := range me.users {
		usersChannel <- user
	}
	close(usersChannel)
	waitGroup.Wait()
	var elapsed = time.Since(beginning)

	return elapsed
}

func (me *SqliteTest) runCombined() (readDuration time.Duration, updateDuration time.Duration) {
	var threadCount = max(1, me.threadCount/2)
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		updateDuration = me.runUpdates(threadCount)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		readDuration = me.runQueries(threadCount)
	}()

	waitGroup.Wait()
	return
}

func (me *SqliteTest) readUsers(users chan *User) {
	var db *sql.DB
	defer func() { me.close(db) }()
	var counter = 0
	for user := range users {
		if nil == db {
			db = me.open()
		}
		var row = db.QueryRow("SELECT name, passwordHash, accessToken, email, createdAt, level FROM users WHERE id=?", user.SqliteId)
		var userB User = User{SqliteId: user.SqliteId}
		var createdAt int64
		assertError(row.Scan(
			&userB.Name, &userB.PasswordHash, &userB.AccessToken, &userB.Email, &createdAt, &userB.Level))
		assertError(row.Err())
		userB.CreatedAt = time.Unix(createdAt, 0)
		assertCondition(user.compare(&userB), "Users must be equal")
		counter += 1
		if (counter%me.batchSize) == 0 && db != nil {
			db = me.close(db)
		}
	}
}

func (me *SqliteTest) updateUsers(users chan *User) {
	var db *sql.DB
	defer func() { me.close(db) }()
	var counter = 0
	for user := range users {
		if nil == db {
			db = me.open()
		}
		assertResultError(db.Exec("UPDATE users SET level=? WHERE id=?", rand.IntN(100), user.SqliteId))
		counter += 1
		if (counter%me.batchSize) == 0 && db != nil {
			db = me.close(db)
		}
	}
}

func (me *SqliteTest) writeUsers(users chan *User) {
	var db *sql.DB
	defer func() { me.close(db) }()
	var counter = 0
	for user := range users {
		if nil == db {
			db = me.open()
		}
		var row = db.QueryRow("INSERT INTO users (name, passwordHash, accessToken, email, createdAt, level) VALUES (?, ?, ?, ?, ?, ?) RETURNING id",
			user.Name, user.PasswordHash, user.AccessToken, user.Email, user.CreatedAt.Unix(), user.Level)
		assertError(row.Scan(&user.SqliteId))
		counter += 1
		if (counter%me.batchSize) == 0 && db != nil {
			db = me.close(db)
		}
	}
}

func (me *SqliteTest) compress() (int64, int64) {
	var db = me.open()
	defer me.close(db)

	switch me.mode {
	case SQLITE_TEST_MODE_FILE:
		assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))
		var sizeBeforeVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
		assertResultError(db.Exec("VACUUM;"))
		assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))

		var sizeAfterVacuum = assertResultError(os.Stat(DB_FILE_PATH)).Size()
		return sizeBeforeVacuum, sizeAfterVacuum
	case SQLITE_TEST_MODE_HTTP:
		var row = db.QueryRow("SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size();")
		var size int64
		assertError(row.Scan(&size))
		return size, size
	default:
		panic("Unknown mode: " + strconv.Itoa(me.mode))
	}
}
