
package statistics

func Count(vi interface{}, cnt int64, nvi interface{}) interface{} {
	return cnt + 1
}

func Min(vi interface{}, cnt int64, nvi interface{}) interface{} {
	if vi == nil {
		return nvi
	}
	
	if v, ok := vi.(int64); ok {
		nv := nvi.(int64)
		if nv > v {
			return v
		}
		return nv
	}

	if v, ok := vi.(float64); ok {
		nv := nvi.(float64)
		if nv > v {
			return v
		}
		return nv
	}

	return nil
}

func Max(vi interface{}, cnt int64, nvi interface{}) interface{} {
	if vi == nil {
		return nvi
	}
	
	if v, ok := vi.(int64); ok {
		nv := nvi.(int64)
		if nv > v {
			return nv
		}
		return v
	}

	if v, ok := vi.(float64); ok {
		nv := nvi.(float64)
		if nv > v {
			return nv
		}
		return v
	}

	return nil
}

func Sum(vi interface{}, cnt int64, nvi interface{}) interface{} {
	if vi == nil {
		return nvi
	}

	if v, ok := vi.(int64); ok {
		nv := nvi.(int64)
		return v + nv
	}

	if v, ok := vi.(float64); ok {
		nv := nvi.(float64)
		return v + nv
	}

	return nil
}

func Avg(vi interface{}, cnt int64, nvi interface{}) interface{} {
	if vi == nil {
		return nvi
	}

	if v, ok := vi.(int64); ok {
		nv := nvi.(int64)
		return (v * int64(cnt) + nv ) / int64(cnt + 1)
	}

	if v, ok := vi.(float64); ok {
		nv := nvi.(float64)
		return (v * float64(cnt) + nv ) / float64(cnt + 1)
	}

	return nil
}