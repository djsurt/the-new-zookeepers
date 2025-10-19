# building report
```bash
docker build . --tag cs249-report-builder
docker run --rm -v "$($PWD):/src" -w '/src' cs249-report-builder make report.pdf
```
