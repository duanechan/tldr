package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/duanechan/tldr/internal/database"
	"github.com/google/uuid"
	"golang.org/x/term"
	"modernc.org/sqlite"
	_ "modernc.org/sqlite"
)

type credentials struct {
	username string
	password string
}

func main() {
	fmt.Println("⚙️  Creating admin user for TL;DR app...")

	credentials, err := readCredentials()
	if err != nil {
		fmt.Printf("❌  error: %s\n", err.Error())
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", "./tldr.db")
	if err != nil {
		fmt.Printf("❌  error: %s\n", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	q := database.New(db)
	ctx := context.Background()
	userId := uuid.Must(uuid.NewRandom())

	_, err = q.CreateUser(ctx, database.CreateUserParams{
		ID:       userId,
		Username: credentials.username,
		Password: credentials.password,
	})
	if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == 2067 {
		fmt.Println("❌  error: username has been taken")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("❌  error: failed to create user")
		os.Exit(1)
	}

	adminId := uuid.Must(uuid.NewRandom())

	_, err = q.CreateAdmin(ctx, database.CreateAdminParams{
		ID:     adminId,
		UserID: userId,
	})
	if err != nil {
		fmt.Println("❌  error: failed to create admin")
		os.Exit(1)
	}

	fmt.Println("✅ Successfully created admin user.")
}

func readCredentials() (*credentials, error) {
	username, err := readInput(
		"Username: ",
		func() (string, error) {
			var u string
			_, err := fmt.Scanln(&u)
			return u, err
		},
	)
	if err != nil {
		return nil, errors.New("failed to read username")
	}

	password, err := readInput("Password: ", readPassword)
	if err != nil {
		return nil, errors.New("failed to read password")
	}

	confirmPassword, err := readInput("Re-type Password: ", readPassword)
	if err != nil {
		return nil, errors.New("failed to read re-typed password")
	}

	if password != confirmPassword {
		return nil, errors.New("passwords do not match")
	}

	password, err = argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	return &credentials{username, password}, nil
}

func readPassword() (string, error) {
	fd := int(os.Stdin.Fd())

	s, err := term.GetState(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, s)

	c, err := term.ReadPassword(fd)
	fmt.Println()
	return string(c), err
}

func readInput(
	prompt string,
	scanFunc func() (string, error),
) (string, error) {
	fmt.Print(prompt)
	input, err := scanFunc()
	if err == nil {
		return strings.TrimSpace(input), nil
	}
	return "", err
}
