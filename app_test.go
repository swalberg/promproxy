package main

import (
	"encoding/json"
	"strings"
	"testing"
)

var rError = []byte(`{"status":"error","data":{}}`)
var matrix1 = []byte(`{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"prometheus_tsdb_blocks_loaded","group":"prometheus","instance":"la3stgprom01","job":"prometheus"},"values":[[1532712425,"24"],[1532712455,"24"],[1532712485,"24"],[1532712515,"24"],[1532712545,"24"],[1532712575,"24"],[1532712605,"24"],[1532712635,"24"],[1532712665,"24"],[1532712695,"24"],[1532712725,"24"],[1532712755,"24"],[1532712785,"24"],[1532712815,"24"],[1532712845,"24"],[1532712875,"24"],[1532712905,"24"],[1532712935,"24"],[1532712965,"24"],[1532712995,"24"],[1532713025,"24"],[1532713055,"24"],[1532713085,"24"],[1532713115,"24"],[1532713145,"24"],[1532713175,"24"],[1532713205,"24"],[1532713235,"24"],[1532713265,"24"],[1532713295,"24"],[1532713325,"24"],[1532713355,"24"],[1532713385,"24"],[1532713415,"24"],[1532713445,"24"],[1532713475,"24"],[1532713505,"24"],[1532713535,"24"],[1532713565,"24"],[1532713595,"24"],[1532713625,"24"],[1532713655,"24"],[1532713685,"24"],[1532713715,"24"],[1532713745,"24"],[1532713775,"24"]]}]}}`)
var matrix2 = []byte(`{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"prometheus_tsdb_blocks_loaded","group":"prometheus","instance":"lv1stgprom01","job":"prometheus"},"values":[[1532712425,"24"],[1532712455,"24"],[1532712485,"24"],[1532712515,"24"],[1532712545,"24"],[1532712575,"24"],[1532712605,"24"],[1532712635,"24"],[1532712665,"24"],[1532712695,"24"],[1532712725,"24"],[1532712755,"24"],[1532712785,"24"],[1532712815,"24"],[1532712845,"24"],[1532712875,"24"],[1532712905,"24"],[1532712935,"24"],[1532712965,"24"],[1532712995,"24"],[1532713025,"24"],[1532713055,"24"],[1532713085,"24"],[1532713115,"24"],[1532713145,"24"],[1532713175,"24"],[1532713205,"24"],[1532713235,"24"],[1532713265,"24"],[1532713295,"24"],[1532713325,"24"],[1532713355,"24"],[1532713385,"24"],[1532713415,"24"],[1532713445,"24"],[1532713475,"24"],[1532713505,"24"],[1532713535,"24"],[1532713565,"24"],[1532713595,"24"],[1532713625,"24"],[1532713655,"24"],[1532713685,"24"],[1532713715,"24"],[1532713745,"24"],[1532713775,"24"]]}]}}`)
var merged = `{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"prometheus_tsdb_blocks_loaded","group":"prometheus","instance":"la3stgprom01","job":"prometheus"},"values":[[1532712425,"24"],[1532712455,"24"],[1532712485,"24"],[1532712515,"24"],[1532712545,"24"],[1532712575,"24"],[1532712605,"24"],[1532712635,"24"],[1532712665,"24"],[1532712695,"24"],[1532712725,"24"],[1532712755,"24"],[1532712785,"24"],[1532712815,"24"],[1532712845,"24"],[1532712875,"24"],[1532712905,"24"],[1532712935,"24"],[1532712965,"24"],[1532712995,"24"],[1532713025,"24"],[1532713055,"24"],[1532713085,"24"],[1532713115,"24"],[1532713145,"24"],[1532713175,"24"],[1532713205,"24"],[1532713235,"24"],[1532713265,"24"],[1532713295,"24"],[1532713325,"24"],[1532713355,"24"],[1532713385,"24"],[1532713415,"24"],[1532713445,"24"],[1532713475,"24"],[1532713505,"24"],[1532713535,"24"],[1532713565,"24"],[1532713595,"24"],[1532713625,"24"],[1532713655,"24"],[1532713685,"24"],[1532713715,"24"],[1532713745,"24"],[1532713775,"24"]]},{"metric":{"__name__":"prometheus_tsdb_blocks_loaded","group":"prometheus","instance":"lv1stgprom01","job":"prometheus"},"values":[[1532712425,"24"],[1532712455,"24"],[1532712485,"24"],[1532712515,"24"],[1532712545,"24"],[1532712575,"24"],[1532712605,"24"],[1532712635,"24"],[1532712665,"24"],[1532712695,"24"],[1532712725,"24"],[1532712755,"24"],[1532712785,"24"],[1532712815,"24"],[1532712845,"24"],[1532712875,"24"],[1532712905,"24"],[1532712935,"24"],[1532712965,"24"],[1532712995,"24"],[1532713025,"24"],[1532713055,"24"],[1532713085,"24"],[1532713115,"24"],[1532713145,"24"],[1532713175,"24"],[1532713205,"24"],[1532713235,"24"],[1532713265,"24"],[1532713295,"24"],[1532713325,"24"],[1532713355,"24"],[1532713385,"24"],[1532713415,"24"],[1532713445,"24"],[1532713475,"24"],[1532713505,"24"],[1532713535,"24"],[1532713565,"24"],[1532713595,"24"],[1532713625,"24"],[1532713655,"24"],[1532713685,"24"],[1532713715,"24"],[1532713745,"24"],[1532713775,"24"]]}]}}`
var labelSet1 = []byte(`{"status":"success","data":["ALERTS","go_gc_duration_seconds","go_gc_duration_seconds_count"]}`)
var labelSet2 = []byte(`{"status":"success","data":["ALERTS","go_gc_duration_seconds","testing"]}`)
var labelMerged = `{"status":"success","data":["ALERTS","go_gc_duration_seconds","go_gc_duration_seconds_count","testing"]}`

var series1 = []byte(`{"status":"success","data":[{"__name__":"mysql_up","dc":"la3","environment":"stg","instance":"la3stgbicubedb01","job":"bicube_mysql-db"},{"__name__":"mysql_up","dc":"la3","environment":"stg","instance":"la3stgbicubedb02","job":"bicube_mysql-db"}]}`)
var series2 = []byte(`{"status":"success","data":[{"__name__":"mysql_up","dc":"lv1","environment":"stg","instance":"lv1stgbicubedb01","job":"bicube_mysql-db"},{"__name__":"mysql_up","dc":"lv1","environment":"stg","instance":"lv1stgbicubedb02","job":"bicube_mysql-db"}]}`)

func TestExtractMatrix(t *testing.T) {
	response := ParseResponse(matrix1)

	matrix := ExtractMatrix(response)
	if len(matrix) != 1 {
		t.Errorf("Expected to see 1 time series, got %d", len(matrix))
	}
}

func TestParseResponseError(t *testing.T) {
	response := ParseResponse(rError)
	if response.Successful() {
		t.Error("expected success to be false but it wasn`t")
	}
}

func TestParseResponse(t *testing.T) {
	response := ParseResponse(matrix1)
	if !response.Successful() {
		t.Error("expected success to be true but it wasn`t")
	}
}

func TestMergeSeries(t *testing.T) {
	p1 := ParseResponse(series1)
	p2 := ParseResponse(series2)
	var two = []ApiResponse{p1, p2}

	m1 := Merge(`series`, two)
	got, err := json.Marshal(m1)
	if err != nil {
		t.Errorf("Unexpected error marshalling response: %s", err)
	}

	i := strings.Count(string(got), `instance`)
	if i != 4 {
		t.Errorf("Expected to see 4 instances but got %d from %s", i, got)
	}
}

func TestMergeArray(t *testing.T) {
	p1 := ParseResponse(labelSet1)
	p2 := ParseResponse(labelSet2)
	var two = []ApiResponse{p1, p2}
	m1 := Merge(`label`, two)
	got, err := json.Marshal(m1)
	if err != nil {
		t.Errorf("Unexpected error unmarshalling response: %s", err)
	}
	if labelMerged != string(got) {
		t.Errorf("Expected %s but got %s", labelMerged, got)
	}
}

func TestMergeMatrix(t *testing.T) {
	p1 := ParseResponse(matrix1)
	p2 := ParseResponse(matrix2)
	var two = []ApiResponse{p1, p2}
	m1 := Merge(`query_range`, two)

	matrix := ExtractMatrix(m1)
	if m1.Status != "success" {
		t.Errorf("Expected success status, got %s", m1.Status)
	}

	var qr QueryResult
	err := json.Unmarshal(m1.Data, &qr)
	if err != nil {
		t.Errorf("Expected to unmarshall the QueryResult, got %s", err)
	}

	if len(matrix) != 2 {
		t.Errorf("Expected to see 2 time series, got %d", len(matrix))
	}
}

func BenchmarkMergeMatrix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var args = []ApiResponse{
			ParseResponse(matrix1),
			ParseResponse(matrix2),
		}
		Merge(`query_range`, args)
	}
}
