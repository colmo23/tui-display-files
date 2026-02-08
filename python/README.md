

python3 -m venv .venv && .venv/bin/pip install -r requirements.txt
.venv/bin/python app.py                 # browse current dir
.venv/bin/python app.py --dir /tmp      # browse specific dir
.venv/bin/pytest test_filesystem.py     # run tests

python3 -m venv 
