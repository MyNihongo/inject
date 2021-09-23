OUTPUT_PATH=${GOPATH}/bin/inject.exe

go build -o ${OUTPUT_PATH} ./inject/
echo "saved: ${OUTPUT_PATH}"
read -p "Press enter to continue..."
