package notable

import (
	"encoding/json"
	"fmt"
	redis "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
	redisurl "github.com/harvesthq/notable/Godeps/_workspace/src/github.com/soveran/redisurl"
	"os"
)

func Notes() []Note {
	db := redisConnection()
	defer db.Close()

	var notes []Note

	for _, id := range noteIDs(db) {
		notes = append(notes, fetchNote(db, id))
	}

	return notes
}

func AddNote(note Note) {
	db := redisConnection()
	defer db.Close()

	id, err := redis.Int64(db.Do("INCR", "notable:note"))
	check(err)

	_, err = db.Do("RPUSH", "notable:notes", id)
	check(err)

	_, err = db.Do("SET", noteKey(id), serialize(note))
	check(err)
}

func Reset() {
	db := redisConnection()
	defer db.Close()

	for _, id := range noteIDs(db) {
		_, err := db.Do("DEL", noteKey(id))
		check(err)
	}

	_, err := db.Do("DEL", "notable:notes")
	check(err)
}

func noteIDs(db redis.Conn) []int64 {
	var ids []int64
	rawIDs, err := redis.Values(db.Do("LRANGE", "notable:notes", 0, -1))
	check(err)

	redis.ScanSlice(rawIDs, &ids)

	return ids
}

func fetchNote(db redis.Conn, id int64) Note {
	noteAsJSON, err := redis.String(db.Do("GET", noteKey(id)))
	check(err)

	return deserialize(noteAsJSON)
}

func noteKey(id int64) string {
	return fmt.Sprintf("notable:note:%d", id)
}

func deserialize(noteAsJSON string) Note {
	var note Note

	err := json.Unmarshal([]byte(noteAsJSON), &note)
	check(err)

	return note
}

func serialize(note Note) string {
	noteAsJSON, err := json.Marshal(note)
	check(err)

	return string(noteAsJSON)
}

func redisConnection() redis.Conn {
	var connection redis.Conn
	var err error

	if len(os.Getenv("REDIS_URL")) > 0 {
		connection, err = redisurl.Connect()
	} else {
		connection, err = redis.Dial("tcp", ":6379")
	}

	check(err)

	return connection
}
