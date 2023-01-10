package main

import logger "gin_study/test_05"

func main() {
	logger.Log.Info("This is a logger")
	logger.Log.Error("is logger")
}
