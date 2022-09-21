package main

type Reducer int

func (r *Reducer) Reduce(in string, reply *string) error {
	*reply = in
	return nil
}
