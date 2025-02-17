package evaluator

import (
	"context"
	"fmt"

	"github.com/mtnmunuklu/alterix/sigma"
)

// evaluateAggregationExpression evaluates an aggregation expression within a Sigma rule
func (rule RuleEvaluator) evaluateAggregationExpression(ctx context.Context, conditionIndex int, aggregation sigma.AggregationExpr) (string, error) {
	var aggregationResult string

	// Determine the type of aggregation expression
	switch agg := aggregation.(type) {
	case sigma.Near:
		return aggregationResult, fmt.Errorf("near isn't supported yet")

	case sigma.Comparison:
		// Evaluate the aggregation function
		aggregationResult, err := rule.evaluateAggregationFunc(ctx, conditionIndex, agg.Func)
		if err != nil {
			return aggregationResult, err
		}

		// Return the aggregation result with the comparison operator and threshold
		return aggregationResult + " " + string(agg.Op) + " " + fmt.Sprintf("%d", int(agg.Threshold)), nil

	default:
		// Return an error if the aggregation expression is not recognized
		return aggregationResult, fmt.Errorf("unknown aggregation expression")
	}
}

// evaluateAggregationFunc evaluates the given aggregation function and returns the resulting query string.
func (rule RuleEvaluator) evaluateAggregationFunc(ctx context.Context, conditionIndex int, aggregation sigma.AggregationFunc) (string, error) {
	var result string
	switch agg := aggregation.(type) {
	case sigma.Count:
		// If the field is not specified, count all records
		if agg.Field == "" {
			// If there is a group by clause, add it to the select statement
			if agg.GroupedBy != "" {
				result = "select " + agg.GroupedBy
			}
			// Add the count function to the select statement
			result += ", count(*)|group having count(*)"
			return result, nil
		} else {
			// If the field is specified, count the number of records for each value of the field
			if len(rule.fieldmappings[agg.Field]) != 0 {
				agg.Field = rule.fieldmappings[agg.Field][0]
			}
			result = "select " + agg.Field
			// If there is a group by clause, add it to the select statement
			if agg.GroupedBy != "" {
				if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
					agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
				}
				result += ", " + agg.GroupedBy
			}
			// Add the count function to the select statement
			result += ", count(*)|group having count(*)"
			return result, nil
		}

	case sigma.Average:
		// Compute the average of the specified field
		if len(rule.fieldmappings[agg.Field]) != 0 {
			agg.Field = rule.fieldmappings[agg.Field][0]
		}
		result = "select " + agg.Field
		// If there is a group by clause, add it to the select statement
		if agg.GroupedBy != "" {
			if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
				agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
			}
			result += ", " + agg.GroupedBy
		}
		// Add the average function to the select statement
		result += ", avg(" + agg.Field + ")|group having avg(" + agg.Field + ")"
		return result, nil

	case sigma.Sum:
		// Compute the sum of the specified field
		if len(rule.fieldmappings[agg.Field]) != 0 {
			agg.Field = rule.fieldmappings[agg.Field][0]
		}
		result = "select " + agg.Field
		// If there is a group by clause, add it to the select statement
		if agg.GroupedBy != "" {
			if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
				agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
			}
			result += ", " + agg.GroupedBy
		}
		// Add the sum function to the select statement
		result += ", sum(" + agg.Field + ")|group having sum(" + agg.Field + ")"
		return result, nil

	case sigma.Min:
		// If the aggregation function is a Min function, map the field to its equivalent in the data source.
		if len(rule.fieldmappings[agg.Field]) != 0 {
			agg.Field = rule.fieldmappings[agg.Field][0]
		}
		// Begin building the query with the SELECT statement and the field to aggregate.
		result = "select " + agg.Field
		// If a group by clause is specified, add it to the query and map the field to its equivalent in the data source.
		if agg.GroupedBy != "" {
			if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
				agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
			}
			result += ", " + agg.GroupedBy
		}
		// Add the aggregation function to the query and set the having clause to filter by the minimum value of the field.
		result += ", min(" + agg.Field + ")|group having min(" + agg.Field + ")"
		return result, nil

	case sigma.Max:
		// If the aggregation function is a Max function, map the field to its equivalent in the data source.
		if len(rule.fieldmappings[agg.Field]) != 0 {
			agg.Field = rule.fieldmappings[agg.Field][0]
		}
		// Begin building the query with the SELECT statement and the field to aggregate.
		result = "select " + agg.Field
		// If a group by clause is specified, add it to the query and map the field to its equivalent in the data source.
		if agg.GroupedBy != "" {
			if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
				agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
			}
			result += ", " + agg.GroupedBy
		}
		// Add the aggregation function to the query and set the having clause to filter by the maximum value of the field.
		result += ", max(" + agg.Field + ")|group having max(" + agg.Field + ")"
		return result, nil

	// If the aggregation function type is not supported, return an error.
	default:
		return result, fmt.Errorf("unsupported aggregation function")
	}
}
