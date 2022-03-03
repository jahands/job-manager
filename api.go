package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

// PutJobById godoc
// @Summary Insert/replace job with specified id
// @Description Overwrites job if it exists.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} Job
// @Router /v1/:jobId [get]
func PutJobById(c *gin.Context) {
	jobId := c.Param("jobId")
	job := Job{
		JobKey:  jobId,
		Created: time.Now(),
		InUse:   false,
	}

	err := rdb.Set(jobKey(jobId), job, 0).Err()
	if err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{"error creating job"})
		return
	}
	c.JSON(200, SuccessResponse{"success"})
}
