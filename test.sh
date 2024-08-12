!bin/bash

go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 10x >> nid_data.txt

go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 10x >> block_data.txt

