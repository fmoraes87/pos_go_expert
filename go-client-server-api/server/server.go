package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fmoraes87/pos_go_expert/go-client-server-api/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var db *gorm.DB

func main() {

	_ = godotenv.Load("../.env")

	initDB()

	http.HandleFunc("/cotacao", QuotationHandler)
	http.ListenAndServe(":8080", nil)

	log.Println("Servidor pronto para receber requisições")
}

func initDB() {
	log.Println("Inicializando conexão com banco de dados")
	var err error

	databaseURL := os.Getenv("DATABASE_PATH")

	db, err = gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar no banco de dados: %v", err)
	}

	err = db.AutoMigrate(&models.Exchange{})
	if err != nil {
		log.Fatalf("Erro ao migrar a tabela: %v", err)
	}

	log.Println("Conexão estabelecida")

}

func QuotationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer log.Println("Request finalizada")
	select {
	case <-ctx.Done():
		log.Println("Request cancelada pelo cliente")
	default:
		processRequest(&ctx, w, r)
	}

}

func processRequest(c *context.Context, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		respondWithError(w, http.StatusNotFound, "Caminho inválido")
		return
	}

	fromCurrency := r.URL.Query().Get("fromCurrency")
	toCurrency := r.URL.Query().Get("toCurrency")

	if fromCurrency == "" || toCurrency == "" {
		respondWithError(w, http.StatusBadRequest, `fromCurrency == "" || toCurrency == "" `)
		return
	}

	quotationRequest := models.QuotationRequest{
		FromCurrency: models.Currency{ISO4217Code: fromCurrency},
		ToCurrency:   models.Currency{ISO4217Code: toCurrency},
		RequestData:  time.Now(),
	}

	quotation, err := getCurrentExchangeRate(c, quotationRequest)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			errMsg := "A requisição excedeu o tempo limite"
			respondWithError(w, http.StatusRequestTimeout, errMsg)
			log.Println(errMsg)
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Erro interno no servidor: %v", err))
		}
		return
	}

	bidResponse := models.BidResponse{
		Bid: quotation.Bid,
	}

	saveQuotationToDB(c, quotationRequest, &bidResponse)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bidResponse)
}

func saveQuotationToDB(c *context.Context, quotationRequest models.QuotationRequest, quotation *models.BidResponse) {
	ctx, cancel := context.WithTimeout(*c, 10*time.Millisecond)
	defer cancel()

	bidValue, err := strconv.ParseFloat(string(quotation.Bid), 64)

	if err != nil {
		log.Fatalf("Erro ao converter Bid: %v", err)
	}

	exchange := models.Exchange{
		FromCurrency:  quotationRequest.FromCurrency.ISO4217Code,
		ToCurrency:    quotationRequest.ToCurrency.ISO4217Code,
		Rate:          bidValue,
		RetrievalDate: quotationRequest.RequestData,
	}

	err = db.WithContext(ctx).Create(&exchange).Error
	if err != nil {
		log.Printf("Erro ao criar a cotação: %v", err)
		return
	}

}

func getCurrentExchangeRate(c *context.Context, quotationRequest models.QuotationRequest) (*models.CurrencyQuotation, error) {
	currencyPair := quotationRequest.GetCurrencyPair()
	url := "https://economia.awesomeapi.com.br/json/last/" + currencyPair

	ctx, cancel := context.WithTimeout(*c, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var q map[string]models.CurrencyQuotation
	err = json.Unmarshal(body, &q)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Erro interno servidor: %v", err)
		return nil, err
	}

	quotationKey := quotationRequest.GetQuotationKey()
	quotation, exists := q[quotationKey]
	if !exists {
		return nil, fmt.Errorf("cotação não encontrada para o par de moedas: %v", quotationKey)
	}

	return &quotation, nil

}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponse := models.ErrorResponse{Message: message}
	json.NewEncoder(w).Encode(errorResponse)
}
