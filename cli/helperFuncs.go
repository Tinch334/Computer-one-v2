package cli

func sliceMap [T, V any](tLst []T, fn func(T) V) []V {
    vLst := make([]V, len(tLst))

    for i, e := range tLst {
        vLst[i] = fn(e)
    }

    return vLst
}