package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Block structure
type Block struct {
	Index        int64     `json:"index"`
	Timestamp    string    `json:"timestamp"`
	Data         string    `json:"data"`
	Hash         string    `json:"hash"`
	PreviousHash string    `json:"previous_hash"`
}

// Blockchain structure
type Blockchain struct {
	Chain []Block
}

// Global variable for the blockchain
var blockchain Blockchain

// Function to create a new block
func createBlock(data string, previousHash string) Block {
	block := Block{
		Index:        int64(len(blockchain.Chain) + 1),
		Timestamp:    time.Now().String(),
		Data:         data,
		Hash:         calculateHash(data, previousHash),
		PreviousHash: previousHash,
	}
	return block
}

// Function to calculate the hash of a block
func calculateHash(data string, previousHash string) string {
	record := fmt.Sprintf("%s%s", data, previousHash)
	h := sha256.New()
	h.Write([]byte(record))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Handler for mining a new block
func mine(c *gin.Context) {
	var newBlock Block
	if err := c.ShouldBindJSON(&newBlock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	previousBlock := blockchain.Chain[len(blockchain.Chain)-1]
	newBlock = createBlock(newBlock.Data, previousBlock.Hash)
	blockchain.Chain = append(blockchain.Chain, newBlock)

	c.JSON(http.StatusCreated, newBlock)
}

// Handler to get the blockchain
func getChain(c *gin.Context) {
	c.JSON(http.StatusOK, blockchain)
}

// Main function
func main() {
	// Initialize the blockchain with the genesis block
	blockchain.Chain = append(blockchain.Chain, createBlock("Genesis Block", "0"))

	// Set up the Gin router
	r := gin.Default()
	r.POST("/mine", mine)
	r.GET("/chain", getChain)

	// Start the server
	r.Run(":5000")
}
