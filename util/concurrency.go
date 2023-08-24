package util

func OrDone[T any](done <-chan struct{}, c <-chan T) chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}

			}
		}
	}()
	return valStream
}
