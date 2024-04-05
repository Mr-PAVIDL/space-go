package commander

import "context"

type Strategy interface {
	Next(ctx context.Context, state *State) Command
}

func s() {

}
