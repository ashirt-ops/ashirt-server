# Setting the import path so that flask development properly runs the code
# (Add "." to the import path searching)
import sys
from os.path import abspath, dirname

sys.path.insert(0, dirname(abspath(__file__)))

from main import create_app

app = create_app()
