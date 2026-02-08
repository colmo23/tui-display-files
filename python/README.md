
Install:
```
python3 -m venv .venv && .venv/bin/pip install -r requirements.txt
```
Browse current dir:
```
.venv/bin/python app.py                 # browse current dir
```
Browse specific dir:
```
.venv/bin/python app.py --dir /tmp      # browse specific dir
```
Tests
```
.venv/bin/pytest test_filesystem.py     # run tests
```

