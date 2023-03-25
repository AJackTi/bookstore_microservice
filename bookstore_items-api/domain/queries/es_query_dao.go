package queries

import "github.com/olivere/elastic"

func (q EsQuery) Build() elastic.Query {
	query := elastic.NewBoolQuery()
	equalQueries := make([]elastic.Query, 0)
	for _, eq := range q.Equals {
		equalQueries = append(equalQueries, elastic.NewMatchQuery(eq.Field, eq.Value))
	}

	query.Must(equalQueries...)
	return query
}
