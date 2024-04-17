package assertions

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"
)

func ToBeEqual(actual interface{}, expected interface{}) error {
	if deepEqual(actual, expected) {
		return nil
	}
	return fmt.Errorf("expected {{SEXP}} %v {{EEXP}} but got {{SGOT}} %v {{EGOT}}", expected, actual)
}

func ToNotBeEqual(actual interface{}, expected interface{}) error {
	if !deepEqual(actual, expected) {
		return nil
	}
	return fmt.Errorf("not expected {{SEXP}} %v {{EEXP}} but got {{SGOT}} %v {{EGOT}}", expected, actual)
}

func ValueCompare(actual interface{}, expected interface{}, opp string) error {

	if opp == "==" || opp == "eq" {
		return ToBeEqual(actual, expected)
	} else if opp == "!=" || opp == "ne" {
		return ToNotBeEqual(actual, expected)
	}
	// if !areSameTypes(actual, expected) {
	// 	//return newAssertionError(needSameType)
	// 	return fmt.Errorf("Not have the same type")
	// }

	var actualF float64
	var err error
	switch x := actual.(type) {
	case json.Number:
		actualF, err = x.Float64()
		if err != nil {
			return err
		}
	default:
		actualF, err = cast.ToFloat64E(actual)
		if err != nil {
			actualS, err := cast.ToStringE(actual)
			if err != nil {
				return err
			}

			expectedS, err := cast.ToStringE(expected)
			if err != nil {
				return err
			}

			switch opp {
			case "<", "lt":
				if actualS < expectedS {
					return nil
				}
				return fmt.Errorf("expected {{SEXP}} %v {{EEXP}} to be less than {{SGOT}} %v {{EGOT}} but it wasn't", actual, expected)

			case "<=", "le":
				if actualS <= expectedS {
					return nil
				}
				return fmt.Errorf("expected {{SEXP}} %v {{EEXP}} to be less than or equal to {{SGOT}} %v {{EGOT}} but it wasn't", actual, expected)

			case ">", "gt":
				if actualS > expectedS {
					return nil
				}
				return fmt.Errorf("expected {{SEXP}} %v {{EEXP}} to be greater than {{SGOT}} %v {{EGOT}} but it wasn't", actual, expected)

			case ">=", "ge":
				if actualS >= expectedS {
					return nil
				}
				return fmt.Errorf("expected {{SEXP}} %v {{EEXP}} to be greater than or equal to {{SGOT}} %v {{EGOT}} but it wasn't", actual, expected)

			default:
				return fmt.Errorf("unknown operator: %v", opp)
			}

		}
	}

	expectedF, err := cast.ToFloat64E(expected)
	if err != nil {
		return err
	}

	switch opp {
	case "<", "lt":
		if actualF < expectedF {
			return nil
		}
		return fmt.Errorf("expected {{SEXP}} %v {{EEXP}} to be less than {{SGOT}} %v {{EGOT}} but it wasn't", actual, expected)

	case "<=", "le":
		if actualF <= expectedF {
			return nil
		}
		return fmt.Errorf("expected {{SEXP}} %v {{EEXP}} to be less than or equal to {{SGOT}} %v {{EGOT}} but it wasn't", actual, expected)

	case ">", "gt":
		if actualF > expectedF {
			return nil
		}
		return fmt.Errorf("expected {{SEXP}} %v {{EEXP}} to be greater than {{SGOT}} %v {{EGOT}} but it wasn't", actual, expected)

	case ">=", "ge":
		if actualF >= expectedF {
			return nil
		}
		return fmt.Errorf("expected {{SEXP}} %v {{EEXP}} to be greater than or equal to {{SGOT}} %v {{EGOT}} but it wasn't", actual, expected)

	default:
		return fmt.Errorf("unknown operator: %v", opp)
	}

}
