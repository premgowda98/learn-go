a = 25

print("A Before", hex(a))

def some(b):
    print("B Before", hex(b))

    b = b+35

    print("B After", hex(b))

    global a
    a = a+25

    print("A After", hex(a))


some(a)
print("Final", a, hex(a))