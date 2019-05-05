package main

import (
	"github.com/powerman/structlog"
	"time"
)

const (
	KeyAddr   = "_addr"
	KeyPlace  = "_place"
	KeyEN6408 = "_en6408"
)

func now() string {
	return time.Now().Format("15:04:05")
}

func withKeys(logger *structlog.Logger, args ...interface{}) *structlog.Logger {
	var keys []string
	for i, arg := range args {
		if i%2 == 0 {
			k, ok := arg.(string)
			if !ok {
				panic("key must be string")
			}
			keys = append(keys, k)
		}
	}
	return logger.New(args...).PrependSuffixKeys(keys...)
}

func withProductAtPlace(logger *structlog.Logger, place int) *structlog.Logger {
	product := prodsMdl.ProductAt(place)
	return withKeys(logger, "место", place+1, "заводской_номер", product.ProductID,
		"product_id", product.ProductID)
}
