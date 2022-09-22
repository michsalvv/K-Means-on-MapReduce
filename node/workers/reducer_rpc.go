package main

import (
	"kmeans-MR/utils"
	"log"
)

type Reducer int

func (r *Reducer) Reduce(in []utils.Point, reply *utils.ReducerResponse) error {
	log.Print(len(in))
	return nil
}
