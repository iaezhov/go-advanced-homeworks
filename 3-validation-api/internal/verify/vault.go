package verify

import (
	"encoding/json"
	"fmt"
)

type VaultEntry struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

type Vault struct {
	List []VaultEntry `json:"list"`
}

type Db interface {
	Read() ([]byte, error)
	Write([]byte) error
}

type VaultWithDb struct {
	Vault
	db Db
}

func NewVault(db Db) *VaultWithDb {
	file, err := db.Read()
	if err != nil {
		return &VaultWithDb{
			Vault: Vault{
				List: []VaultEntry{},
			},
			db: db,
		}
	}
	var vault Vault
	err = json.Unmarshal(file, &vault)
	if err != nil {
		fmt.Println(err)
		return &VaultWithDb{
			Vault: Vault{
				List: []VaultEntry{},
			},
			db: db,
		}
	}
	return &VaultWithDb{
		Vault: vault,
		db:    db,
	}
}

func (vault *Vault) ToBytes() ([]byte, error) {
	file, err := json.Marshal(vault)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (vault *VaultWithDb) Save() {
	data, err := vault.ToBytes()
	if err != nil {
		fmt.Println("Не удалось преобразовать в JSON")
		return
	}
	vault.db.Write(data)
}

func (vault *VaultWithDb) Add(input VaultEntry) {
	vault.List = append(vault.List, input)
	vault.Save()
}

func (vault *VaultWithDb) Delete(hash string) bool {
	var list []VaultEntry
	isDeleted := false
	for _, value := range vault.List {
		isMatched := value.Hash == hash
		if !isMatched {
			list = append(list, value)
			continue
		}
		isDeleted = true
	}
	vault.List = list
	vault.Save()
	return isDeleted
}
