// Code generated by "stringer -type=Type"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-128]
	_ = x[EOF-129]
	_ = x[WORD-130]
	_ = x[DIGITS-131]
	_ = x[PLUS-132]
	_ = x[MINUS-133]
	_ = x[BANG-134]
	_ = x[ASTERISK-135]
	_ = x[SLASH-136]
	_ = x[IMPL-137]
	_ = x[MOD-138]
	_ = x[LT-139]
	_ = x[GT-140]
	_ = x[EXISTS-141]
	_ = x[NEXISTS-142]
	_ = x[UNION-143]
	_ = x[EQ-144]
	_ = x[NEQ-145]
	_ = x[DOT-146]
	_ = x[COLON-147]
	_ = x[COMMA-148]
	_ = x[PIPE-149]
	_ = x[SUM-150]
	_ = x[LPAREN-151]
	_ = x[RPAREN-152]
	_ = x[LBRACKET-153]
	_ = x[RBRACKET-154]
	_ = x[AS-155]
	_ = x[IF-156]
	_ = x[IN-157]
	_ = x[MAX-158]
	_ = x[MIN-159]
	_ = x[TO-160]
	_ = x[BOOL-161]
}

const _Type_name = "ILLEGALEOFWORDDIGITSPLUSMINUSBANGASTERISKSLASHIMPLMODLTGTEXISTSNEXISTSUNIONEQNEQDOTCOLONCOMMAPIPESUMLPARENRPARENLBRACKETRBRACKETASIFINMAXMINTOBOOL"

var _Type_index = [...]uint8{0, 7, 10, 14, 20, 24, 29, 33, 41, 46, 50, 53, 55, 57, 63, 70, 75, 77, 80, 83, 88, 93, 97, 100, 106, 112, 120, 128, 130, 132, 134, 137, 140, 142, 146}

func (i Type) String() string {
	i -= 128
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i+128), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
