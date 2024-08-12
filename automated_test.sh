!bin/bash
cd /home/jony/PQHBEVS_HAC_proto_p2p

rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 1x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 1x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 5x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 5x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 10x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 10x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 20x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 20x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 50x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 50x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 75x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 75x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 100x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 100x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 125x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 125x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 150x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 150x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 200x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 200x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 500x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 500x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 750x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 750x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 1000x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 1000x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 1250x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 1250x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 1500x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 1500x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 2000x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 2000x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 2500x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 2500x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 3000x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 3000x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 3500x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 3500x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 4000x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 4000x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 5000x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 5000x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt


rm -rv /home/jony/sod_dbms/data*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/BLOCK
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/EC/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/RO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/internal/role/PO/*
rm -rv /home/jony/PQHBEVS_HAC_proto_p2p/test/levelDB
go test -benchmem -run=^$ -bench ^Benchmark_NID_SERVER$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 7500x >> nid_data.txt
curl --location --request POST 'http://localhost:5669/generate_genesis?pscode=8010&seed=jony&start=1681490380&duration=10800&config=2&verbose=0'
echo "----------------------------------------------------------------------------------">>nid_data.txt
go test -benchmem -run=^$ -bench ^Benchmark_BlockGenearationRate$ github.com/sohelahmedjoni/pqhbevs_hac/test -benchtime 7500x >> newblock_data.txt -timeout 60m
echo "----------------------------------------------------------------------------------">> newblock_data.txt
