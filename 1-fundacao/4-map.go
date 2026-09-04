package fundacao

import "fmt"

/*
	Map é uma estrutura de dados que permite armazenar uma chave e um valor


	make = cria um map com tamanho e capacidade definidos
	delete = deleta um elemento do map
	for range = itera sobre o map
	len = retorna o tamanho do map
	cap = retorna a capacidade do map
	
	
	
	
*/
func MapTeste() {

	// mapa := map[string]int{
	// 	"apple":  1,
	// 	"banana": 2,
	// 	"cherry": 3,
	// }

	// fmt.Printf("mapa = %v\n len=%d\n", mapa, len(mapa))

	// delete(mapa, "apple")
	// fmt.Printf("mapa = %v\n len=%d\n", mapa, len(mapa))

	// make = cria um map com tamanho e capacidade definidos
	mapa2 := make(map[string]int)
	mapa2["apple"] = 1
	mapa2["banana"] = 2
	mapa2["cherry"] = 3
	mapa2["orange"] = 4
	fmt.Printf("len=%d %v\n", len(mapa2), mapa2)

	for fruta, quantidade := range mapa2 {
		fmt.Printf("fruta = %s, quantidade = %d\n", fruta, quantidade)
	}
}
