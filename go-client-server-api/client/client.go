package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fmoraes87/pos_go_expert/go-client-server-api/models"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load("../.env")

	fromCurrency := "USD"
	toCurrency := "BRL"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	serverURL := os.Getenv("SERVER_HOST")

	url := fmt.Sprintf("%s/cotacao?fromCurrency=%s&toCurrency=%s", serverURL, fromCurrency, toCurrency)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Erro ao realizar a chamada: %v", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stdout, "Erro ao obter resposta do servidor: %v", string(body))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Erro ao converter a resposta: %v", err)
		return
	}

	var bid models.BidResponse
	err = json.Unmarshal(body, &bid)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Erro ao processar resposta do servidor: %v", err)
		return
	}

	HandleSaveQuotation(&bid)

	fmt.Fprintf(os.Stdout, "Cotação de USD para BRL: R$ %s", bid.Bid)

}

func HandleSaveQuotation(q *models.BidResponse) {
	fileDir := os.Getenv("FILE_PATH_STORAGE")
	file, err := os.Create(fileDir + "/cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stdout, "Erro ao criar arquivo cotacaol.txt: %v", err)
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", q.Bid))
	if err != nil {
		fmt.Fprintf(os.Stdout, "Erro ao escrever no arquivo cotacaol.txt: %v", err)
		return
	}

	fmt.Println("Arquivo criado/atualizado com sucesso!")
}
