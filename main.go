package main

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	dbUri := "neo4j://localhost:7687"
	driver, err := neo4j.NewDriver(dbUri, neo4j.BasicAuth("neo4j", "102080", ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close()

	item, err := insertItem(driver)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", item)

	set, err := setItem(driver)
	if err != nil {
		panic(err)
	}
	fmt.Println(set)
}

func matchItem(driver neo4j.Driver) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(
			`match (n:Item{name: $name})
			 return n
			`, map[string]any{
				"name": "Item 1",
			})
		if err != nil {
			return nil, err
		}
		return
	})
}

func setItem(driver neo4j.Driver) (string, error) {

	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
		match (n:Item{name: "Item 1"})
		set n.surname = "anupom"
		return n`,
			map[string]interface{}{
				"id":   2,
				"name": "Item 2",
			})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return "", err
	}

	return "greeting.(string)", nil

}

func insertItem(driver neo4j.Driver) (*Item, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.WriteTransaction(createItemFn)

	if err != nil {
		return nil, err
	}
	return result.(*Item), nil
}

func createItemFn(tx neo4j.Transaction) (interface{}, error) {
	records, err := tx.Run("CREATE (n:Item{id: $id, name: $name }) RETURN n.id, n.name", map[string]interface{}{
		"id":   2,
		"name": "Item 2",
	})
	if err != nil {
		return nil, err
	}

	record, err := records.Single()
	if err != nil {
		return nil, err
	}

	return &Item{
		Id:   record.Values[0].(int64),
		Name: record.Values[1].(string),
	}, nil
}

type Item struct {
	Id   int64
	Name string
}
