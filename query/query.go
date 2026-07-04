package query

import "strings"

const (
	FilterOpUnspecified      int32 = 0
	FilterOpExact            int32 = 1
	FilterOpLike             int32 = 2
	FilterOpGT               int32 = 3
	FilterOpLT               int32 = 4
	FilterOpGTE              int32 = 5
	FilterOpLTE              int32 = 6
	FilterOpNull             int32 = 7
	FilterOpNotNull          int32 = 8
	FilterOpLikeCI           int32 = 9
	FilterOpIn               int32 = 10
	FilterOpBetween          int32 = 11
	FilterOpBetweenExclusive int32 = 12

	OrderDirectionUnspecified int32 = 0
	OrderDirectionAsc         int32 = 1
	OrderDirectionDesc        int32 = 2

	LogicalOpUnspecified int32 = 0
	LogicalOpAnd         int32 = 1
	LogicalOpOr          int32 = 2
)

func FilterOpFromString[O ~int32](op string) O {
	switch strings.ToLower(op) {
	case "like":
		return O(FilterOpLike)
	case "gt":
		return O(FilterOpGT)
	case "lt":
		return O(FilterOpLT)
	case "gte":
		return O(FilterOpGTE)
	case "lte":
		return O(FilterOpLTE)
	case "null":
		return O(FilterOpNull)
	case "not_null":
		return O(FilterOpNotNull)
	case "like_ci":
		return O(FilterOpLikeCI)
	case "in":
		return O(FilterOpIn)
	case "between":
		return O(FilterOpBetween)
	case "between_exclusive":
		return O(FilterOpBetweenExclusive)
	default:
		return O(FilterOpExact)
	}
}

func OrderDirectionFromString[O ~int32](direction string) O {
	if strings.EqualFold(direction, "desc") {
		return O(OrderDirectionDesc)
	}
	return O(OrderDirectionAsc)
}

func LogicalOpFromString[L ~int32](op string) L {
	if strings.EqualFold(op, "or") {
		return L(LogicalOpOr)
	}
	return L(LogicalOpAnd)
}
