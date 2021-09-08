package main

import (
  "github.com/openrelayxyz/plugeth-utils/core"
)

var (
  logger core.Logger
)

type myservice struct{}

func (*myservice) Hello() string {
  return "Hello world"
}

func InitializeNode(l core.Logger, node core.Node, backend core.Backend) {
  logger = l
  logger.Info("Initialized hello")
}

func GetAPIs(node core.Node, backend core.Backend) []core.API {
  defer logger.Info("APIs Initialized")
  return []core.API{
    {
      Namespace: "mynamespace",
      Version: "1.0",
      Service: &myservice{},
      Public: true,
    },
  }
}
