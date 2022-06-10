package main

import (
	"errors"
	"fmt"
	"net/http"
	"pundix-homework/clients"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/libs/math"
)

func setupRoutes(engine *gin.Engine) {
	queryGroup := engine.Group("/query")
	// distribution
	distributionGroup := queryGroup.Group("/distribution")
	{
		distributionGroup.GET("queryParams", QueryParamsHandler)
		distributionGroup.GET("communityPool", CommunityPoolHandler)
		distributionGroup.GET("validatorCommission", ValidatorCommissionHandler)
		distributionGroup.GET("validatorOutstandingRewards", ValidatorOutstandingRewardsHandler)
		// distributionGroup.GET("validatorSlashes", ValidatorSlashesHanlder)
	}

	// bank
	bankGroup := queryGroup.Group("/bank")
	{
		bankGroup.GET("balance", BalanceHandler)
		bankGroup.GET("total", TotalSupplyHandler)
	}
}

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func QueryParamsHandler(c *gin.Context) {
	res, err := clients.DistrQueryClientInstance.QueryParams()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func ValidatorCommissionHandler(c *gin.Context) {
	validator := c.Query("validator")

	if validator == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validator empty"})
		return
	}

	res, err := clients.DistrQueryClientInstance.ValidatorCommission(validator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func parseParams(c *gin.Context) (string, uint64, uint64, uint64, error) {
	validator := c.Query("validator")

	startHeightStr := c.Query("startHeight")
	endHeightStr := c.Query("endHeight")
	limitStr := c.Query("limit")

	startHeight, err0 := strconv.ParseInt(startHeightStr, 10, 64)
	endHeight, err1 := strconv.ParseInt(endHeightStr, 10, 64)
	limit, err2 := strconv.ParseInt(limitStr, 10, 64)
	if err0 != nil || err1 != nil || err2 != nil {
		fmt.Println(err0, err1, err2)
		return "", 0, 0, 0, errors.New("invalid int type")
	}

	if startHeight < 0 || endHeight < 0 || limit < 0 {
		return "", 0, 0, 0, errors.New("int params must be positive")
	}
	endHeight = math.MaxInt64(1, endHeight)
	limit = math.MaxInt64(1, limit)

	if validator == "" {
		return "", 0, 0, 0, errors.New("validator is empty")
	}
	return validator, uint64(startHeight), uint64(endHeight), uint64(limit), nil
}

func ValidatorSlashesHanlder(c *gin.Context) {
	validator, startHright, endHeight, limit, err := parseParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := clients.DistrQueryClientInstance.ValidatorSlashes(validator, startHright, endHeight, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func ValidatorOutstandingRewardsHandler(c *gin.Context) {
	validator := c.Query("validator")
	if validator == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validator empty"})
		return
	}
	res, err := clients.DistrQueryClientInstance.ValidatorOutstandingRewards(validator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func CommunityPoolHandler(c *gin.Context) {
	res, err := clients.DistrQueryClientInstance.CommunityPool()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func BalanceHandler(c *gin.Context) {
	address := c.Query("address")

	res, err := clients.BankQueryClientInstance.Balance(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func TotalSupplyHandler(c *gin.Context) {
	res, err := clients.BankQueryClientInstance.TotalSupply()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
