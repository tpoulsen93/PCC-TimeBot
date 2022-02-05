from datetime import datetime as d

f = open("bob.txt", "a")
f.write(str(d.now()))
f.close()

