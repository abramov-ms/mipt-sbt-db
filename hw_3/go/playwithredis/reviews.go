package playwithredis

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
)

type ID int64

func (id ID) MarshalBinary() ([]byte, error) {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(id))
	return bytes, nil
}

type Review struct {
	Text  string `json:"review"`
	Score int    `json:"score"`
	Type  string `json:"type"`
}

func (r *Review) Dump() {
	fmt.Printf("Type: %s\n", r.Type)
	fmt.Printf("Score: %d\n", r.Score)
	fmt.Println(r.Text)
}

type Movie struct {
	ID      ID       `json:"id"`
	Reviews []Review `json:"reviews"`
}

func (m *Movie) Dump() {
	fmt.Printf("ID = %d\n", m.ID)
	fmt.Println("Reviews:")
	for _, r := range m.Reviews {
		r.Dump()
	}
}

func ParseJSON(filename string) ([]Movie, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	movies := make([]Movie, 0)
	if err := json.Unmarshal(bytes, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}
