# A collection of tools for the [Nand2Tetris project.](http://nand2tetris.org/)

## Components

Everything needed for a basic assembler and compiler.  A Lexer (assembly -> tokens), parser (tokens -> machine code) and a compiler (watch this space, will have own lexer/parser).

I'm going with a fully fledged lexer/parser for the initial Assembler project so that it can be reused in the compiler project.  The Lexer is blatantly ripped from [Rob Pike's talk here](https://www.youtube.com/watch?v=HxaD_trXwRE).

It's totally overkill for the assembler, (which is really just an exercise in substitution, and would be easy enough to do with regular expressions (eugh)), but I figure this will be more challenging and a good way to get my head around Go (plus partial reuse in the compiler project).

*Postscript:* It turns out that I'm a masochist.  This WAS a good way to learn a new language, but maybe next time don't also learn how to write something like a lexer (which I'm still not happy about, it doesn't feel 'clean' in the same way that the parser does) at the same time too.

## Assembler

Basic assembler that maps symbols/tokens to machine instructions.  Output is a text file with "binary" values written out as string.  Internally they're all going to be represented by 16 bit constants that are then OR'd together and converted to a string representation at the end.  This is because it sounds more 'program-y' but mainly because I cannot bring myself to write this using string concatenation (plus it's good practice as I'm learning Go at the same time).

Bonus points: Handles both Unix and Windows line endings, and has a warning for redundant A-Instructions (i.e. @123 followed by @456 is redundant, @123 will have no effect).

Still to do - tidy up Lexer, and in fact make it dumber.  Right now it's doing a fair bit or error checking that could probably be done more easily in the parser, making the lexer code cleaner.

## Compiler

TBD
