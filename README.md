# EDI
EDI library in Go

# ElectronicDataInterchange

# Transmission Structure

```
[ISA]-- Interchange ----------------------
:  [GS]-- Functional Group ---------------
:  :  [ST]-- Transaction Set -------------
:  :  :   ,-- Detail Segment------------.
:  :  :  |                               |
:  :  :  |    DataElement*DataElement    |
:  :  :  |                               |
:  :  :   `-- Detail Segment------------`
:  :  [SE]-- Transaction Set -------------
:  [GE]-- Functional Group ---------------
[IEA]-- Transmission Envelope ------------
 ```
 
 # TODO
 - [ ] @seanbeagle: Write Data Dictionary API to download formatting rules based on interchange version
 - [ ] @seanbeagle: Write validator.go to validate formatting of interchange
