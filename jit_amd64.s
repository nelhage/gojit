// textflag.h
// Don't profile the marked routine.  This flag is deprecated.
#define NOPROF	1
// It is ok for the linker to get multiple of these symbols.  It will
// pick one of the duplicates to use.
#define DUPOK	2
// Don't insert stack check preamble.
#define NOSPLIT	4
// Put this data in a read-only section.
#define RODATA	8
// This data contains no pointers.
#define NOPTR	16
// This is a wrapper function and should not count as disabling 'recover'.
#define WRAPPER 32
// end textflag.h

// cgocall(*args) with jitcode in the context blob
//   -> runtime·cgocall(jitcode, frame)
TEXT ·cgocall(SB),NOSPLIT,$16
        LEAQ argframe+0(FP), AX
        MOVQ AX, 8(SP)
        MOVQ 8(DX), AX
        MOVQ AX, 0(SP)
        CALL runtime·cgocall(SB)
        RET

TEXT ·jitcall(SB),NOSPLIT,$0
        LEAQ argframe+0(FP), DI
        MOVQ 8(DX), AX
        JMP AX
