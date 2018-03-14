---
title: "Escape Standard Library Reference"
slug: scripting-language-stdlib 
type: "docs"
toc: true
---

<style>
h2 {
  font-size: 0.8em;
  font-family: mono;
  background: #4B9CD3;
  padding: 5px;
}
</style>

Standard library functions for the [Escape Scripting Language](../scripting-language/)


# Functions acting on lists

## list_index(n :: integer)

Index a list at position `n`. Usually accessed implicitly using indexing syntax (eg. `list[0]`)

## list_slice(i :: integer, j :: integer)

Slice a list. Usually accessed implicitly using slice syntax (eg. `list[0:5]`)

## env_lookup(key :: string)

Lookup key in environment. Usually called implicitly when using '$'

## join(sep :: string)

Join concatenates the elements of a to create a single string. The separator string sep is placed between elements in the resulting string. 


# Functions acting on strings

## track_minor_version()

Track minor version

## track_version()

Track version

## base64_encode()

Encode string to base64

## trim()

Returns a slice of the string s, with all leading and trailing white space removed, as defined by Unicode. 

## split(sep :: string)

Split slices s into all substrings separated by sep and returns a slice of the substrings between those separators. If sep is empty, Split splits after each UTF-8 sequence.

## base64_decode()

Decode string from base64

## read_file()

Read the contents of a file

## upper(v :: string)

Returns a copy of the string v with all Unicode characters mapped to their upper case

## title(v :: string)

Returns a copy of the string v with all Unicode characters mapped to their title case

## concat(v1 :: string, v2 :: string, ...)

Concatate stringable arguments

## track_major_version()

Track major version

## track_patch_version()

Track patch version

## lower(v :: string)

Returns a copy of the string v with all Unicode characters mapped to their lower case

## replace(old :: string, new :: string, n :: integer)

Replace returns a copy of the string s with the first n non-overlapping instances of old replaced by new. If old is empty, it matches at the beginning of the string and after each UTF-8 sequence, yielding up to k+1 replacements for a k-rune string. If n < 0, there is no limit on the number of replacements.


# Functions acting on integers

## add(y :: integer)

Add two integers


# Unary functions

## timestamp()

Returns a UNIX timestamp


# Functions acting on everything

## id(parameter :: *)

Returns its argument

