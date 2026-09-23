package bancodedados

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Produto struct {
	ID    uint   `gorm:"primaryKey"`
	Nome  string `gorm:"size:100"`
	Preco float64
}

func GoORMTest() {
	// Teste de conexão com o banco de dados usando GORM
	// Aqui você pode adicionar o código para inicializar a conexão e realizar operações de teste

	db, err := gorm.Open(postgres.Open("host=localhost port=5439 user=root password=root dbname=golang_db sslmode=disable"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Produto{}) // Cria a tabela "produtos" no banco de dados, se não existir

	if err != nil {
		panic(err)
	}

	// produto := Produto{Nome: "Produto Teste", Preco: 19.99}

	// err = CreateProduto(db, &produto)

	// if err != nil {
	// 	panic(err)
	// }

	// produtoRecuperado, err := GetProdutoByID(db, 1)

	// if err != nil {
	// 	panic(err)
	// }

	// println("Produto recuperado:", produtoRecuperado.Nome, "Preço:", produtoRecuperado.Preco)

	println("Todos os produtos:") // Lista todos os produtos no banco de dados
	produtos, err := GetAllProdutos(db)

	if err != nil {
		panic(err)
	}

	for _, p := range produtos {
		println("ID:", p.ID, "Nome:", p.Nome, "Preço:", p.Preco)
	}

	//Delete
	deleteID := uint(1) // ID do produto a ser deletado
	err = DEleteProduto(db, deleteID)

	if err != nil {
		panic(err)
	}

	println("Produto com ID", deleteID, "deletado com sucesso.")
	println("Todos os produtos após a exclusão:")

	println("-----") // Atualizado
	var produtoAtualizado Produto

	db.First(&produtoAtualizado, 2)
	produtoAtualizado.Nome = "Produto Atualizado"
	produtoAtualizado.Preco = 29.99

	db.Save(&produtoAtualizado)

	for _, p := range produtos {
		println("ID:", p.ID, "Nome:", p.Nome, "Preço:", p.Preco)
	}
}

func CreateProduto(db *gorm.DB, produto *Produto) error {
	result := db.Create(produto)
	return result.Error
}

func GetProdutoByID(db *gorm.DB, id uint) (*Produto, error) {
	var produto Produto
	result := db.First(&produto, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &produto, nil
}

func GetAllProdutos(db *gorm.DB) ([]Produto, error) {
	var produtos []Produto
	result := db.Find(&produtos)

	if result.Error != nil {
		return nil, result.Error
	}

	return produtos, nil
}

func DEleteProduto(db *gorm.DB, id uint) error {
	result := db.Delete(&Produto{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
