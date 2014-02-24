TEXT ·Log2(SB),7,$0     
        BSRQ 8(SP), AX  // 8(SP) is the first argument
        MOVQ AX, 16(SP) // store result
        RET
        