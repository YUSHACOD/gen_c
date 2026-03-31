
# General overview of the format
Every generator primitive like data(tables, file embed, etc..), or commands to generate data 
(enums, structs, types, global, etc...) are represented at the top level by @<primitive_name>
It also contains primitives like gen_h_file, gen_c_file, gen_hpp_file, gen_cpp_file  
eg: 
```
@table(Funcs)
@enum(FuncKind)
@enum_to_string_table(FuncKindToString)
@func_globals(FuncGlobals)
@custom(FuncCustom)


@gen_h_file(func)
@gen_c_file(func)
@gen_hpp_file(func)
@gen_cpp_file(func)
```

Primitives may or may not have sub primitives
eg:
```
@table(Funcs) {
	@fields(name type_name args ret) // sub primitive
	[Add                add_op_ft                  `int x, int y`         int]
	[Sub                sub_op_ft                  `int x, int y`         int]
	[Mul                mul_op_ft                  `int x, int y`         int]
}

@enum(FuncKind)
{
	@requires(Funcs f)
	@identifiers(f.name)
}

```

Sub Primtives may be of two kind, 
- 1. 'special', related to the parent primitive: they fill up primitive data
eg:
```
@table(Funcs) {
	@fields(name type_name args ret) // related to the table primitive
	[Add                add_op_ft                  `int x, int y`         int]
	[Sub                sub_op_ft                  `int x, int y`         int]
	[Mul                mul_op_ft                  `int x, int y`         int]
}
```

- 2. 'global', not top level stuff but describes behaviour to the generator
eg:
```
@enum(FuncKind)
{
	@requires(Funcs f) // not related to the @enum primitive but describes dependancies to generator
	@identifiers(@concat(`FK_`, f.name)) // @concat is one of many utility functions to do general
                                         // operations like concatenating, converting uppercase to
                                         // lower, pascal case to snake, and so on and so forth
}


@enum_to_string_table(FuncKindToString, FunKind) // global primitive
```


# Spec

## List of primitives and their special subprimitives

- 1. table -> a simple table with rows and coloumns, to store related data. All the cells may
     contain the following three type of data:
      - 1. C/C++ identifier string alpahnumeric + '_' charachter, shouldn't start with number.
      - 2. '_' just this string, describes that the cell is empty.
      - 3. anything enclosed in backticks(``) should be valid c syntax when composed.

    ## Sub Primitives
    
    table contains only one sub primitive 'fields' this describe all of the fields for a row for
    identification


- 2. enum -> a simple primitve to generate c style enum with

    ## Sub Primitives
    
    'identifier' to descirbe what each enum values identifier should be



- 3. enum_to_string_table -> takes a enum primitive and composes the enum_to_string_table global for 
     it



- 4. func_types -> creates types for functions. 

    ## Sub Primitives

    1. identifier -> type identifier string
    2. args -> arguments for the func type
    3. ret -> return type string for the func type


- 5. func_global -> uses the func types and creates global pointer for the types

    ## Sub Primitives

    1. identifier -> identifier string for the func global 
    2. type -> func type string



- 6. custom -> custom go's text/template template to evaluate and generate c / cpp code.
    
    1. template -> the template string enclosed in backticks







