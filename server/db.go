package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type runDB struct {
	items map[string]string
	mu    sync.RWMutex
}

//--- Create a runtime Database
func createDB() runDB {
	f, err := os.Open("db.json")
	if err != nil {
		return runDB{items: map[string]string{}}
	}
	items := map[string]string{}
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		fmt.Println("--- Couldn't decode", err.Error())
		return runDB{items: map[string]string{}}
	}
	return runDB{items: items}
}

//--- Save data to a db.json file
func (r runDB) save() {
	f, err := os.Create("db.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(r.items); err != nil {
		log.Fatal(err)
	}
}

//--- Set data into the runtime Database
func (r runDB) set(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key] = value
}

//--- Get the data from the database
func (r runDB) get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, found := r.items[key]
	return value, found
}

//--- Delete the data from the database
func (r runDB) delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, key)
}
