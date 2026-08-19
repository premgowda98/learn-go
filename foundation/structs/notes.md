## Structs

1. Simmilar to **Dataclass** and **Pydantic** in python
2. When functions attached to structs they are called as **methods**.
3. **Mutation Methods**, those functions which alter the struct data
    1. In mutation function we should pointer as recivere argument, since we are changing the struct data
    2. If not sent as pointer, then only the copy will be edited and not the original object
4. **Creator/Constructor Functions**
    1. Starts with new keyword
    2. Can be used for validation
5. When structs are added in a package, in order to import you need to name all functions with uppercase ane also the fields inside struct should start with upercase
6. **Struct Embedding**
    1. Works like inheritance
    2. Inherits all methods of other struct