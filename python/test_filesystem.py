import os
import tempfile

from filesystem import read_dir, read_file_content


def test_read_dir():
    with tempfile.TemporaryDirectory() as tmp_dir:
        # Create test files and directories
        with open(os.path.join(tmp_dir, "file1.txt"), "w") as f:
            f.write("hello")
        os.mkdir(os.path.join(tmp_dir, "dir1"))
        with open(os.path.join(tmp_dir, ".hiddenfile"), "w") as f:
            f.write("hidden")
        os.mkdir(os.path.join(tmp_dir, ".git"))

        entries = read_dir(tmp_dir)

        assert len(entries) == 2
        names = sorted(e.name for e in entries)
        assert names == ["dir1", "file1.txt"]


def test_read_file_content():
    with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
        f.write("hello world")
        tmp_path = f.name

    try:
        content = read_file_content(tmp_path)
        assert content == "hello world"
    finally:
        os.unlink(tmp_path)
