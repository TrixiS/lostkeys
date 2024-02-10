package commands

import (
	"fmt"

	"github.com/TrixiS/lostkeys/pkg/command_context"
	"github.com/boltdb/bolt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"golang.design/x/clipboard"
)

func Set(ctx *command_context.CommandContext) error {
	if ctx.CLIContext.NArg() < 2 {
		return fmt.Errorf("provide a key and a value")
	}

	key := ctx.CLIContext.Args().Get(0)
	value := ctx.CLIContext.Args().Get(1)

	db := ctx.Provider.DBFactory()
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := getValuesBucket(tx)

		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), []byte(value))
	})

	if err != nil {
		return fmt.Errorf("set value %w", err)
	}

	fmt.Printf("set %s -> %v\n", key, value)
	return nil
}

func Get(ctx *command_context.CommandContext) error {
	if ctx.CLIContext.NArg() < 1 {
		return fmt.Errorf("provide an existent key")
	}

	key := ctx.CLIContext.Args().First()

	db := ctx.Provider.DBFactory()
	defer db.Close()

	value, err := getValue(db, key)

	if err != nil {
		return fmt.Errorf("get value %w", err)
	}

	if value == nil {
		return fmt.Errorf("key %s not found", key)
	}

	fmt.Println(string(value))
	return nil
}

func List(ctx *command_context.CommandContext) error {
	db := ctx.Provider.DBFactory()
	defer db.Close()

	headerColor := color.New(color.FgBlue, color.Underline).SprintfFunc()
	idColor := color.New(color.FgBlue, color.Bold).SprintfFunc()

	t := table.New("Key", "Value").
		WithHeaderFormatter(headerColor).
		WithFirstColumnFormatter(idColor)

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := getValuesBucket(tx)

		if err != nil {
			return err
		}

		return bucket.ForEach(func(k, v []byte) error {
			t.AddRow(string(k), string(v))
			return nil
		})
	})

	if err != nil {
		return err
	}

	t.Print()
	return nil
}

func Clipboard(ctx *command_context.CommandContext) error {
	if ctx.CLIContext.NArg() < 1 {
		return fmt.Errorf("provide an existent key")
	}

	key := ctx.CLIContext.Args().First()

	db := ctx.Provider.DBFactory()
	defer db.Close()

	value, err := getValue(db, key)

	if err != nil {
		return fmt.Errorf("get value %w", err)
	}

	if value == nil {
		return fmt.Errorf("key %s not found", key)
	}

	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("clipboard init %w", err)
	}

	clipboard.Write(clipboard.FmtText, value)

	fmt.Printf("value of the key %s has been written to your clipboard\n", key)
	return nil
}

func getValuesBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	var valuesBucketName []byte = []byte("values")
	return tx.CreateBucketIfNotExists(valuesBucketName)
}

func getValue(db *bolt.DB, key string) ([]byte, error) {
	tx, err := db.Begin(true)

	if err != nil {
		return nil, err
	}

	bucket, err := getValuesBucket(tx)

	if err != nil {
		return nil, err
	}

	defer tx.Commit()
	return bucket.Get([]byte(key)), nil
}
