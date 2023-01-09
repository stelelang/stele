%{
package main
%}

%union {
	val int
}

%type <val> expr
%token <val> yyINT

%left '+' '-'
%left '*' '/'

%%

expr:
		yyINT { $$ = $1 }
		| expr '+' expr { $$ = $1 + $3 }
		| expr '-' expr { $$ = $1 - $3 }
		| expr '*' expr { $$ = $1 * $3 }
		| expr '/' expr { $$ = $1 / $3 }
		| '(' expr ')' { $$ = $2 }
		;

%%
