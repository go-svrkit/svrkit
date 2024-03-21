// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

#include "textflag.h"

// func getg() *g
TEXT ·getg(SB),NOSPLIT,$0-8
    MOVW g, ret+0(FP)
	RET
