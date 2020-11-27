TESTDIR="testData/b400"

# Start nodes
./node/node -conf=$TESTDIR/nodes-0.txt -cpu=4 --loglevel=4 &
./node/node -conf=$TESTDIR/nodes-1.txt -cpu=4 --loglevel=4 &
./node/node -conf=$TESTDIR/nodes-2.txt -cpu=4 --loglevel=4 &

sleep 2

./client/client --conf=$TESTDIR/client.txt -metric=1 -batch=400

killall ./node/node
