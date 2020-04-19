echo "Generating templates with statik..."
rm -f templates/*.zip
cd templates
zip -o snoop-js.zip snoop.js
cd ..
cp dist/snoop.zip templates/snoop-go.zip
statik -src templates -include=*.zip -dest internal/templates -f
