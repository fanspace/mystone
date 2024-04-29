package main

import (
	"backend/db"
	log "backend/logger"
)

func main() {
	db.InitDb()
	defer db.CloseDb()
	log.InitLogger()
}
