package parser

import "go/token"

func InspectToken(tok token.Token) (result TokenResult) {
	switch tok {
	// Special tokens
	case token.ILLEGAL:
	case token.EOF:
	case token.COMMENT:

	// literal_beg:
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	case token.IDENT: // main
	case token.INT: // 12345
	case token.FLOAT: // 123.45
	case token.IMAG: // 123.45i
	case token.CHAR: // 'a'
	case token.STRING: // "abc"
	// literal_end

	// operator_beg
	// Operators and delimiters
	case token.ADD: // +
	case token.SUB: // -
	case token.MUL: // *
	case token.QUO: // /
	case token.REM: // %

	case token.AND: // &
	case token.OR: // |
	case token.XOR: // ^
	case token.SHL: // <<
	case token.SHR: // >>
	case token.AND_NOT: // &^

	case token.ADD_ASSIGN: // +=
	case token.SUB_ASSIGN: // -=
	case token.MUL_ASSIGN: // *=
	case token.QUO_ASSIGN: // /=
	case token.REM_ASSIGN: // %=

	case token.AND_ASSIGN: // &=
	case token.OR_ASSIGN: // |=
	case token.XOR_ASSIGN: // ^=
	case token.SHL_ASSIGN: // <<=
	case token.SHR_ASSIGN: // >>=
	case token.AND_NOT_ASSIGN: // &^=

	case token.LAND: // &&
	case token.LOR: // ||
	case token.ARROW: // <-
	case token.INC: // ++
	case token.DEC: // --

	case token.EQL: // ==
	case token.LSS: // <
	case token.GTR: // >
	case token.ASSIGN: // =
	case token.NOT: // !

	case token.NEQ: // !=
	case token.LEQ: // <=
	case token.GEQ: // >=
	case token.DEFINE: // :=
	case token.ELLIPSIS: // ...

	case token.LPAREN: // (
	case token.LBRACK: // [
	case token.LBRACE: // {
	case token.COMMA: // ,
	case token.PERIOD: // .

	case token.RPAREN: // )
	case token.RBRACK: // ]
	case token.RBRACE: // }
	case token.SEMICOLON: // ;
	case token.COLON: // :
	// operator_end

	// keyword_beg
	// Keywords
	case token.BREAK:
	case token.CASE:
	case token.CHAN:
	case token.CONST:
	case token.CONTINUE:

	case token.DEFAULT:
	case token.DEFER:
	case token.ELSE:
	case token.FALLTHROUGH:
	case token.FOR:

	case token.FUNC:
	case token.GO:
	case token.GOTO:
	case token.IF:
	case token.IMPORT:

	case token.INTERFACE:
	case token.MAP:
	case token.PACKAGE:
	case token.RANGE:
	case token.RETURN:

	case token.SELECT:
	case token.STRUCT:
	case token.SWITCH:
	case token.TYPE:
	case token.VAR:
		// keyword_end
	}

	return
}
