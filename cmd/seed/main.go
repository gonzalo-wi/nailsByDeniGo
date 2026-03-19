// Utilitario para crear el primer admin en la base de datos.
// Uso: go run cmd/seed/main.go -name="Shei" -email="admin@sheinails.com" -password="tupassword" -role="superadmin"
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	name := flag.String("name", "", "Nombre del admin (requerido)")
	email := flag.String("email", "", "Email del admin (requerido)")
	password := flag.String("password", "", "Contraseña en texto plano (requerido)")
	role := flag.String("role", "admin", "Rol: admin | superadmin")
	flag.Parse()

	if *name == "" || *email == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "Uso: go run cmd/seed/main.go -name=<name> -email=<email> -password=<password> [-role=superadmin]")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error generando hash: %v", err)
	}

	fmt.Println("─── SQL para insertar el admin ─────────────────────────────")
	fmt.Printf(`INSERT INTO admins (name, email, password_hash, role, active, created_at, updated_at)` + "\n")
	fmt.Printf(`VALUES ('%s', '%s', '%s', '%s', true, NOW(), NOW());`+"\n", *name, *email, string(hash), *role)
	fmt.Println("────────────────────────────────────────────────────────────")
}
