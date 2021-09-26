OUTPUT_PATH=${GOPATH}/bin/mockgen.exe

go build -o ${OUTPUT_PATH} ./mockgen/
echo "saved: ${OUTPUT_PATH}"
read -p "Press enter to continue..."