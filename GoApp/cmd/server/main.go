package main

import (
	"github.com/7empestx/GoHTMXToDoList/internal/db"
	"github.com/7empestx/GoHTMXToDoList/internal/auth"
	"github.com/7empestx/GoHTMXToDoList/internal/router"
)

func main() {
  db.Init();
  auth.Init();
  router.Init();
}
